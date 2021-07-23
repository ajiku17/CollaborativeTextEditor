package p2p

import (
	"fmt"
	"github.com/pion/webrtc/v3"
	"sync"
	"time"
)

type ConnectionStateChangeCallback func ()
type ChannelCreateCallback func ()
type MessageCallback func (msg []byte)

type ICECandidateCallback func(candidate ICECandidateMsg)
type ConnOfferCallback func(offer ConnOffer)
type ConnAnswerCallback func(answer ConnAnswer)

type PeerConn struct {
	endpointId string

	OnConnectionStateChange ConnectionStateChangeCallback

	OnChannelCreateCallback ChannelCreateCallback
	OnMessageCallback       MessageCallback

	Conn          *webrtc.PeerConnection
	Channel       *webrtc.DataChannel

	mu         sync.Mutex
	terminated bool

	onICECandidateCallback ICECandidateCallback
	onConnAnswerCallback   ConnAnswerCallback
}

func NewConn(endpoint string) *PeerConn {
	p := new(PeerConn)

	p.endpointId = endpoint
	p.terminated = false

	return p
}

func (p *PeerConn) GetEndpoint() string {
	return p.endpointId
}

func (p *PeerConn) OnICECandidateReceived(callback ICECandidateCallback) {
	p.onICECandidateCallback = callback
}

func (p *PeerConn) OnAnswer(callback ConnAnswerCallback) {
	//fmt.Println("setting endpoint", p.GetEndpoint(), "answer handler")
	p.onConnAnswerCallback = callback
}

func (p *PeerConn) OnMessage(callback MessageCallback) {
	p.OnMessageCallback = callback
}

func (p *PeerConn) Close() error {
	if err := p.Conn.Close(); err != nil {
		//fmt.Printf("cannot close peer.Conn: %v\n", err)
		return err
	}

	return nil
}

func (p *PeerConn) SendMessage(data []byte) error {
	//fmt.Println("sending message len", len(data), "to", p.GetEndpoint())
	err := p.Channel.Send(data)

	return err
}

func (m *P2P) handleSignals(p *PeerConn, msgQueue chan interface{}, errc chan error) {
	//fmt.Println(m.peerId, "handling signals from", p.GetEndpoint(), "answer callback: ", p.onConnAnswerCallback)
	for {
		p.mu.Lock()
		terminated := p.terminated
		p.mu.Unlock()

		if terminated {
			return
		}

		timer := time.After(10 * time.Second)

		select {
		case msg := <- msgQueue:
			//fmt.Println(m.peerId, "received data from", p.endpointId)

			switch msg.(type) {
			case ICECandidateMsg:
				//fmt.Println("received ice candidate message", msg)
				p.onICECandidateCallback(msg.(ICECandidateMsg))
			case ConnAnswer:
				//fmt.Println("received answer")
				p.onConnAnswerCallback(msg.(ConnAnswer))
			case ConnRefuse:
				//fmt.Println("connection refused")
				errc <- fmt.Errorf("connection refused")
			}
		case <- timer:
			//fmt.Println("peer signal timer fired off")
			continue
		}
	}
}