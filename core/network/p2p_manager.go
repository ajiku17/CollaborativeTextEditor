package network

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/p2p"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/tracker"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"sync"
	"time"
)


const STUN_URL = "stun:stun.l.google.com:19302"

type P2PManager struct {
	id utils.UUID
	doc synceddoc.Document

	p2p *p2p.P2P

	signalingURL string

	trackerC *tracker.Client

	conns   map[*p2p.PeerConn]struct{}
	connsMu sync.Mutex

	inbound  chan P2PMessage
	outbound chan P2PMessage

	stopped bool
	killed bool
	mu sync.Mutex
}

func NewP2PManager(siteId utils.UUID, doc synceddoc.Document, signalingURL string, track *tracker.Client) Manager {
	m := new(P2PManager)

	m.id = siteId
	m.doc = doc
	m.stopped = false
	m.killed = false

	m.signalingURL = signalingURL
	m.trackerC = track

	m.inbound = make(chan P2PMessage, 100)
	m.outbound = make(chan P2PMessage, 100)

	m.conns = make(map[*p2p.PeerConn]struct{})

	return m
}

func (m *P2PManager) GetId() utils.UUID {
	return m.id
}

func (m *P2PManager) Start() {
	m.stopped = false

	m.p2p = p2p.New(m.signalingURL, string(m.id), STUN_URL)

	peers, err := m.trackerC.RegisterAndGet(string(m.doc.GetID()), string(m.id))
	if err != nil {
		fmt.Println("P2P manager: error registering and getting from tracker", err)
	}

	peers, err = m.trackerC.Get(string(m.doc.GetID()))
	if err != nil {
		fmt.Println("P2P manager: error registering and getting from tracker", err)
	}

	// Setup p2p
	m.p2p.OnPeerConnectionRequest(func(conn *p2p.PeerConn, offer p2p.ConnOffer, aux interface{}) {
		conn.OnMessage(func (msg []byte) {
			p2pMsg, err := DecodeP2PMessage(msg)
			//fmt.Printf("%s received a message from %s %s\n", m.p2p.GetPeerId(), conn.GetEndpoint(), string(msg))
			if err == nil {
				m.inbound <- p2pMsg
			}
		})
	})

	m.p2p.OnPeerConnection(func(endpointPeerId string, conn *p2p.PeerConn, aux interface{}) {
		//fmt.Println(m.p2p.GetPeerId(), "received a connection from ", conn.GetEndpoint())

		m.connsMu.Lock()
		defer m.connsMu.Unlock()

		m.conns[conn] = struct{}{}

		docState := m.doc.GetCurrentState()
		//fmt.Println(m.id, "sending patch request to newly established", conn.GetEndpoint())
		m.outbound <- P2PMessage{
			Sender:   string(m.id),
			Receiver: conn.GetEndpoint(),
			Data:     nil,
			IsPatch:  false,
			Patch:    nil,
			IsState:  true,
			State:    docState,
		}


		m.doc.(*synceddoc.SyncedDocument).OnPeerConnect(utils.UUID(endpointPeerId), 0, nil)

	})

	m.p2p.OnPeerDisconnection(func(endpointPeerId string, conn *p2p.PeerConn, aux interface{}) {
		m.doc.(*synceddoc.SyncedDocument).OnPeerDisconnect(utils.UUID(endpointPeerId), nil)
	})

	err = m.p2p.Start()
	if err != nil {
		fmt.Println("P2P manager: error starting p2p", err)
	}

	fmt.Println(m.id, "received peer ids", peers)

	// Setup connections
	//wg := sync.WaitGroup{}
	for _, p := range peers {
		if string(m.id) == p {
			continue
		}

		//wg.Add(1)
		go func (peerId string) {
			conn := p2p.NewConn(peerId)

			conn.OnMessage(func(msg []byte) {
				//fmt.Printf("%s received a message from %s %s\n", m.p2p.GetPeerId(), conn.GetEndpoint(), string(msg))
				p2pMsg, err := DecodeP2PMessage(msg)
				if err == nil {
					m.inbound <- p2pMsg
				}
			})

			//fmt.Println(m.id, "manager setting up connection with", conn.GetEndpoint())
			err := m.p2p.SetupConn(conn, conn.GetEndpoint())
			//wg.Done()

			if err != nil {
				//fmt.Println(m.id, "P2P manager: error while setting up connection with", conn.GetEndpoint(), err)
				return
			}


			m.connsMu.Lock()
			m.conns[conn] = struct{}{} // save newly setup connection
			m.connsMu.Unlock()

			//fmt.Println(m.id, "sending patch request to newly established", conn.GetEndpoint())

		} (p)
	}

	// wait for connections
	//wg.Wait()

	m.startSynchronizer()
}

func (m *P2PManager) startSynchronizer() {
	fmt.Println("starting startSynchronizer")
	go m.changeMonitor()
	go m.sender()
	go m.requestProcessor()
	go m.backgroundSync()
}

func (m *P2PManager) backgroundSync() {
	for {
		var killed bool
		m.mu.Lock()
		killed = m.killed
		m.mu.Unlock()

		if killed {
			fmt.Println(m.id, "background sync has been killed")
			return
		}

		docState := m.doc.GetCurrentState()
		m.outbound <- P2PMessage{
			Sender:   string(m.id),
			Receiver: "",
			Data:     nil,
			IsPatch:  false,
			Patch:    nil,
			IsState:  true,
			State:    docState,
		}

		time.Sleep(5 * time.Second) // sleep for a while
	}
}

func (m *P2PManager) changeMonitor() {
	lastChangeIndex := -1
	//fmt.Println(m.id, "starting change monitor")
	for {
		var killed bool
		m.mu.Lock()
		killed = m.killed
		m.mu.Unlock()
		if killed {
			fmt.Println(m.id, "change monitor has been killed")
			return
		}

		changesAfter, ind := m.doc.GetLocalOpsFrom(lastChangeIndex)
		lastChangeIndex = ind

		for _, op := range changesAfter {
			opData, err := synceddoc.EncodeOp(op)
			if err != nil {
				fmt.Println("failed to encode Op")
				continue
			}

			//fmt.Println(m.id, "sending new op with index", op.PeerOpIndex, "to others")
			m.outbound <- P2PMessage{
				Sender:   string(m.id),
				Receiver: "",
				Data:     opData,
				IsPatch:  false,
				Patch:    nil,
				IsState:  false,
				State:    nil,
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func (m *P2PManager) sender () {
	for {
		var killed bool
		m.mu.Lock()
		killed = m.killed
		m.mu.Unlock()

		if killed {
			fmt.Println(m.id, "sender has been killed")
			return
		}

		timer := time.After(5 * time.Second)

		select {
		case msg := <- m.outbound:
			if msg.Receiver == "" {
				m.sendToAll(msg)
			} else {
				m.sendToPeer(msg.Receiver, msg)
			}

		case <- timer:
			fmt.Println(m.id, "sender timer fired off")
		}
	}
}

func (m *P2PManager) sendToAll(msg P2PMessage) {
	//fmt.Println(m.id, "trying to send")
	m.connsMu.Lock()
	defer m.connsMu.Unlock()
	//fmt.Println(m.id, "took lock to send", m.conns)
	for conn := range m.conns {
		//fmt.Println(m.id, "sending to", conn.GetEndpoint())
		go m.sendToConn(conn, msg)
	}
}

func (m *P2PManager) sendToPeer(peerId string, msg P2PMessage) {
	//fmt.Println(m.id, "trying to send")
	m.connsMu.Lock()
	defer m.connsMu.Unlock()
	//fmt.Println(m.id, "took lock to send", m.conns)
	for conn := range m.conns {
		if conn.GetEndpoint() == peerId {
			//fmt.Println(m.id, "sending to", conn.GetEndpoint())
			go m.sendToConn(conn, msg)
		}
	}
}
func (m *P2PManager) sendToConn(conn *p2p.PeerConn, msg P2PMessage) {
	byteMsg, err := EncodeP2PMessage(msg)
	if err != nil {
		fmt.Println(m.id, "sender error: failed to encode msg", err)
		return
	}

	err = conn.SendMessage(byteMsg)
	if err != nil {
		fmt.Println(m.id, "sender error: failed to send msg to", conn.GetEndpoint(), "error:", err)
		return
	}
}

func (m *P2PManager) requestProcessor () {
	for {
		fmt.Println("Request processed by p2p manager ", m.stopped)


		m.mu.Lock()
		stopped := m.stopped
		m.mu.Unlock()
		if stopped {
			continue
		}

		var killed bool
		m.mu.Lock()
		killed = m.killed
		m.mu.Unlock()

		if killed {
			fmt.Println(m.id, "receiver has been killed")
			return
		}

		timer := time.After(5 * time.Second)

		select {
		case msg := <- m.inbound:
			sendRsp, response, err := m.processRequest(msg)
			if err != nil {
				fmt.Println(m.id, "error processing request:", err)
				continue
			}

			if sendRsp {
				m.outbound <- response
			}

		case <- timer:
			fmt.Println(m.id, "receiver timer fired off")
		}
	}
}

func (m *P2PManager) processRequest(msg P2PMessage) (bool, P2PMessage, error) {
	if msg.Receiver != "" && msg.Receiver != string(m.id) {
		return false, P2PMessage{}, fmt.Errorf("%v, received somebody else's msg. receiver: %v", m.id, msg.Receiver)
	}

	if msg.IsPatch {
		m.doc.ApplyPatch(msg.Patch)
		return false, P2PMessage{}, nil
	}

	if msg.IsState {
		p := m.doc.CreatePatch(msg.State)
		return true, P2PMessage{
			Sender:   string(m.id),
			Receiver: msg.Sender,
			Data:     nil,
			IsPatch:  true,
			Patch:    p,
			IsState:  false,
			State:    nil,
		}, nil
	}

	op, err := synceddoc.DecodeOp(msg.Data)
	if err != nil {
		return false, P2PMessage{}, err
	}

	//fmt.Println(m.id, "applying op index", op.PeerOpIndex, "from", op.PeerId)
	m.doc.ApplyRemoteOp(op, nil)

	return false, P2PMessage{}, nil
}

func (m *P2PManager) Stop() {
	m.stopped = true
	m.p2p.Stop()
}

func (m *P2PManager) Kill() {
	m.connsMu.Lock()
	defer m.connsMu.Unlock()

	m.killed = true
}