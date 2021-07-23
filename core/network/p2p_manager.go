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

	trackerC tracker.Client

	conns   map[*p2p.PeerConn]struct{}
	connsMu sync.Mutex

	inbound  chan P2PMessage
	outbound chan P2PMessage

	killed bool
	mu sync.Mutex
}

func New (siteId utils.UUID, doc synceddoc.Document, signalingURL string, track tracker.Client) Manager {
	m := new(P2PManager)

	m.id = siteId
	m.doc = doc
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
	m.p2p = p2p.New(m.signalingURL, string(m.id), STUN_URL)

	err := m.trackerC.Register(string(m.doc.GetID()), string(m.id))
	if err != nil {
		fmt.Println("P2P manager: error registering with tracker", err)
	}

	peers, err := m.trackerC.Get(string(m.doc.GetID()))
	if err != nil {
		fmt.Println("P2P manager: error fetching peer list from tracker", err)
	}

	// Setup p2p
	m.p2p.OnPeerConnectionRequest(func(conn *p2p.PeerConn, offer p2p.ConnOffer, aux interface{}) {
		conn.OnMessage(func (msg []byte) {
			fmt.Printf("%s received a message from %s %s\n", m.p2p.GetPeerId(), conn.GetEndpoint(), string(msg))
			p2pMsg, err := DecodeP2PMessage(msg)
			if err == nil {
				m.inbound <- p2pMsg
			}
		})
	})

	m.p2p.OnPeerConnection(func(endpointPeerId string, conn *p2p.PeerConn, aux interface{}) {
		fmt.Println(m.p2p.GetPeerId(), "received a connection from ", conn.GetEndpoint())

		m.connsMu.Lock()
		defer m.connsMu.Unlock()

		m.conns[conn] = struct{}{}
	})

	err = m.p2p.Start()
	if err != nil {
		fmt.Println("P2P manager: error starting p2p", err)
	}

	// Setup connections
	wg := sync.WaitGroup{}
	for _, p := range peers {
		wg.Add(1)
		go func (peerId string) {
			conn := p2p.NewConn(peerId)

			conn.OnMessage(func(msg []byte) {
				fmt.Printf("%s received a message from %s %s\n", m.p2p.GetPeerId(), conn.GetEndpoint(), string(msg))
				p2pMsg, err := DecodeP2PMessage(msg)
				if err == nil {
					m.inbound <- p2pMsg
				}
			})

			err := m.p2p.SetupConn(conn, conn.GetEndpoint())
			wg.Done()

			if err != nil {
				fmt.Println("P2P manager: error while setting up connection with", conn.GetEndpoint())
				return
			}

			m.connsMu.Lock()
			m.conns[conn] = struct{}{} // save newly setup connection
			m.connsMu.Unlock()
		} (p)
	}

	// wait for connections
	wg.Wait()

	go m.synchronizer()
}

func (m *P2PManager) synchronizer() {
	go m.sender()
	go m.requestProcessor()
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
			m.connsMu.Lock()
			for conn := range m.conns {
				go func(conn* p2p.PeerConn, p2pMsg P2PMessage) {
					byteMsg, err := EncodeP2PMessage(p2pMsg)
					if err != nil {
						fmt.Println(m.id, "sender error: failed to encode msg", p2pMsg)
						return
					}

					err = conn.SendMessage(byteMsg)
					if err != nil {
						fmt.Println(m.id, "sender error: failed to send msg to", conn.GetEndpoint())
						return
					}

				}(conn, msg)
			}
			fmt.Println(msg)
			m.connsMu.Unlock()
		case <- timer:
			fmt.Println(m.id, "sender timer fired off")
		}
	}
}

func (m *P2PManager) requestProcessor () {
	for {
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
			fmt.Println(msg)
			sendRsp, response, err := m.processRequest(msg)
			if err != nil {
				fmt.Println(m.id, "error processing request from")
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

	return false, P2PMessage{}, nil
}

func (m *P2PManager) Stop() {
	m.p2p.Stop()
}

func (m *P2PManager) Kill() {
	m.connsMu.Lock()
	defer m.connsMu.Unlock()

	m.killed = true
}