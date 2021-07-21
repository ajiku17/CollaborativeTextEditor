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

func (p *PeerConn) OnICECandidateReceived(callback ICECandidateCallback) {
	p.onICECandidateCallback = callback
}

func (p *PeerConn) OnAnswer(callback ConnAnswerCallback) {
	p.onConnAnswerCallback = callback
}

func (p *PeerConn) OnChannelCreate(fn func ()) {
	p.OnChannelCreateCallback = fn

	p.Conn.OnDataChannel(func(d *webrtc.DataChannel) {
		p.Channel = d
		
		d.OnOpen(func() {
			p.OnChannelCreateCallback()

			for range time.NewTicker(5 * time.Second).C {
				message := "Hello world"
				fmt.Printf("Sending '%s'\n", message)

				// Send the message as text
				sendTextErr := d.SendText(message)
				if sendTextErr != nil {
					panic(sendTextErr)
				}
			}
		})
		
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("%s Message from DataChannel '%s': '%s'\n", p.endpointId, p.Channel.Label(), string(msg.Data))
			p.OnMessageCallback(msg.Data)
		})
	})
}

func (p *PeerConn) OnMessage(callback MessageCallback) {
	p.OnMessageCallback = callback
}

func (p *PeerConn) Close() error {
	if err := p.Conn.Close(); err != nil {
		fmt.Printf("cannot close peer.Conn: %v\n", err)
		return err
	}

	return nil
}

func (p *PeerConn) SendMessage(data []byte) error {
	err := p.Channel.Send(data)

	return err
}

func (p *PeerConn) handleSignals(msgQueue chan interface{}) {
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
			fmt.Println("p2p received data", msg, "from", p.endpointId)

			switch msg.(type) {
			case ICECandidateMsg:
				fmt.Println("received ice candidate message", msg)
				p.onICECandidateCallback(msg.(ICECandidateMsg))
			case ConnAnswer:
				fmt.Println("received answer", msg)
				p.onConnAnswerCallback(msg.(ConnAnswer))
			}
		case <- timer:
			fmt.Println("peer conn timer fired off")
			continue
		}
	}
}