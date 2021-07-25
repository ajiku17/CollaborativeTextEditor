package p2p

import (
	"fmt"
	"github.com/pion/webrtc/v3"
	"sync"
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
	p.mu.Lock()
	p.terminated = true
	p.mu.Unlock()

	fmt.Println("closing connection to", p.endpointId)

	err := p.Channel.Close()
	if err != nil {
		fmt.Printf("cannot close channel. endpoint: %v, error: %v\n", p.endpointId, err)
		return nil
	}

	err = p.Conn.Close()

	if err != nil {
		fmt.Printf("cannot close peer.Conn. endpoint: %v, error: %v\n", p.endpointId, err)
		return err
	}

	return nil
}

func (p *PeerConn) SendMessage(data []byte) error {
	//fmt.Println("sending message len", len(data), "to", p.GetEndpoint())
	err := p.Channel.Send(data)

	return err
}