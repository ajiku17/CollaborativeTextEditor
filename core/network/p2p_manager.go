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

	track  tracker.Client

	conns   map[*p2p.PeerConn]struct{}
	connsMu sync.Mutex

	inbound chan []byte

	killed bool
	mu sync.Mutex
}

func New (siteId utils.UUID, doc synceddoc.Document, signalingURL string, track tracker.Client) Manager {
	m := new(P2PManager)

	m.id = siteId
	m.doc = doc
	m.killed = false

	m.signalingURL = signalingURL
	m.track = track

	m.inbound = make(chan []byte, 100)
	m.conns = make(map[*p2p.PeerConn]struct{})

	return m
}

func (m *P2PManager) GetId() utils.UUID {
	return m.id
}

func (m *P2PManager) Start() {
	m.p2p = p2p.New(m.signalingURL, string(m.id), STUN_URL)

	err := m.track.Register(string(m.doc.GetID()), string(m.id))
	if err != nil {
		fmt.Println("P2P manager: error registering with tracker", err)
	}

	peers, err := m.track.Get(string(m.doc.GetID()))
	if err != nil {
		fmt.Println("P2P manager: error fetching peer list from tracker", err)
	}

	// Setup p2p
	m.p2p.OnPeerConnectionRequest(func(conn *p2p.PeerConn, offer p2p.ConnOffer, aux interface{}) {
		conn.OnMessage(func (msg []byte) {
			fmt.Printf("%s received a message from %s %s\n", m.p2p.GetPeerId(), conn.GetEndpoint(), string(msg))
			m.inbound <- msg
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
				m.inbound <- msg
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
	go m.receiver()
}

func (m *P2PManager) sender () {
	for {
		var killed bool
		m.mu.Lock()
		killed = m.killed
		m.mu.Unlock()

		if killed {
			fmt.Println(m.id, "synchronizer has been killed")
			return
		}

		timer := time.After(5 * time.Second)

		select {
		case msg := <- m.inbound:
			//m.doc.ApplyRemoteOp()
			fmt.Println(msg)
		case <- timer:
			fmt.Println(m.id, "synchronizer timer fired off")
		}
	}
}

func (m *P2PManager) receiver () {
	for {
		var killed bool
		m.mu.Lock()
		killed = m.killed
		m.mu.Unlock()

		if killed {
			fmt.Println(m.id, "synchronizer has been killed")
			return
		}

		timer := time.After(5 * time.Second)

		select {
		case msg := <- m.inbound:
			//m.doc.ApplyRemoteOp()
			fmt.Println(msg)
		case <- timer:
			fmt.Println(m.id, "synchronizer timer fired off")
		}
	}
}

func (m *P2PManager) Stop() {
	m.p2p.Stop()
}

func (m *P2PManager) Kill() {
	m.connsMu.Lock()
	defer m.connsMu.Unlock()

	m.killed = true
}