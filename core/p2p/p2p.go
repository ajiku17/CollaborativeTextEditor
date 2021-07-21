package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/signaling"
	"github.com/pion/webrtc/v3"
	"sync"
	"time"
)

type PeerConnectionCallback func(string, *PeerConn, interface{})
type PeerConnectionRequestCallback func (*PeerConn, ConnOffer, interface{})

type ApplyMsg struct {
	CommandValid   bool
	Command        interface{}
	CommandIndex   int
	Snapshot       bool
	SnapshotData []byte
}

type P2P struct {
	mu      sync.Mutex

	stopped bool
	conn    *signaling.Client
	peerId  string
	config  webrtc.Configuration

	peerConnectionCallback        PeerConnectionCallback
	peerConnectionRequestCallback PeerConnectionRequestCallback

	inbound  chan []byte
	outbound chan []byte

	applyCh chan ApplyMsg
	msgQueues map[string] chan interface{}
}

func New(signalingURL string, peerId string, stunURL string, applyCh chan ApplyMsg) *P2P {
	c := new(P2P)

	cn := signaling.NewClient(context.Background(), signalingURL, peerId)

	c.stopped = false
	c.conn = cn
	c.peerConnectionCallback = nil
	c.peerId = peerId
	c.applyCh = applyCh
	c.inbound = make(chan []byte, 16)
	c.outbound = make(chan []byte, 16)
	c.msgQueues = make(map[string] chan interface{})
	c.config = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{stunURL},
			},
		},
	}

	return c
}

func (p *P2P) OnPeerConnection(callback PeerConnectionCallback) {
	p.peerConnectionCallback = callback
}

func (p *P2P) OnPeerConnectionRequest(callback PeerConnectionRequestCallback) {
	p.peerConnectionRequestCallback = callback
}

func (p *P2P) Start() error {
	err := p.conn.Dial()
	if err != nil {
		fmt.Println("p2p signaling dial error:", err)
		return err
	}

	go p.receiver()
	go p.sender()

	errc := make(chan error, 10)
	go p.msgProcessor(errc)

	return nil
}

func (p *P2P) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stopped = true
}

func (p *P2P) receiver() {
	for {
		p.mu.Lock()
		stopped := p.stopped
		p.mu.Unlock()

		if stopped {
			fmt.Println(p.peerId, "receiver stopped")
			return
		}

		fmt.Println("waiting for message")
		message, err := p.conn.NextMessage()
		fmt.Println(p.peerId, "received payload")
		if err != nil {
			fmt.Println("p2p receiver error", err)
			continue
		}

		p.inbound <- message
	}
}

func (p *P2P) sender() {
	for {
		p.mu.Lock()
		stopped := p.stopped
		p.mu.Unlock()

		if stopped {
			return
		}

		timer := time.After(2 * time.Second)

		select {
		case msg := <- p.outbound:
			fmt.Println("p2p sending data", msg)

			err := p.conn.SendPayload("peer", msg)
			if err != nil {
				fmt.Println("p2p sender error:", err)
				continue
			}
		case <- timer:
			continue
		}
	}
}

func (p *P2P) msgProcessor(errc chan error) {
	for {
		msg := <- p.inbound

		err := p.processMsg(msg)
		if err != nil {
			errc <- err
			return
		}
	}
}

func (p *P2P) processMsg(msg []byte) error {
	m := P2PMessage{}

	err := json.Unmarshal(msg, &m)
	if err != nil {
		fmt.Println("invalid P2P message", msg)
		return err
	}

	fmt.Println(p.peerId, "received message", m)
	switch m.MsgType {
	case CONN_OFFER:
		offer := ConnOffer{}
		err := json.Unmarshal(m.Msg, &offer)
		if err != nil {
			fmt.Println("invalid ConnOffer", m.Msg)
			return err
		}
		fmt.Println(p.peerId, "received offer", offer)
		go p.connectionRequested(offer)
	case CONN_ANSWER:
		answer := ConnAnswer{}
		err := json.Unmarshal(m.Msg, &answer)
		if err != nil {
			fmt.Println("invalid ConnOffer", m.Msg)
			return err
		}
		err = p.dispatchMsg(answer.Sender, answer)
	case ICE_CANDIDATE:
		fmt.Println(p.peerId, "received ice candidate from signaling server")
		icecandidate := ICECandidateMsg{}
		err := json.Unmarshal(m.Msg, &icecandidate)
		if err != nil {
			fmt.Println("invalid ConnOffer", m.Msg)
			return err
		}
		err = p.dispatchMsg(icecandidate.Sender, icecandidate)
	}

	return nil
}

func (p *P2P) dispatchMsg(peerId string, msg interface{}) error {
	p.mu.Lock()
	c, ok := p.msgQueues[peerId]
	p.mu.Unlock()

	if !ok {
		return fmt.Errorf("message queue not found for %v", peerId)
	}

	c <- msg

	return nil
}

func (p *P2P) DistributeMsg() {
	fmt.Println("requesting Msg distribution")
}

func (p *P2P) RequestSnapshot() {
	fmt.Println("requesting snapshot")
}

func (p *P2P) signalAnswer(peerId string, answer webrtc.SessionDescription) error {
	sdp, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	ans := ConnAnswer {
		Sender: p.peerId,
		SDP:    sdp,
	}

	ansBytes, err := json.Marshal(ans)
	if err != nil {
		return err
	}

	msg := P2PMessage {
		MsgType: CONN_ANSWER,
		Msg:     ansBytes,
	}

	err = p.signalMessage(peerId, msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *P2P) signalICECandidate(peerId string, c *webrtc.ICECandidate) error {

	ice := ICECandidateMsg {
		Sender:       p.peerId,
		IceCandidate: c.ToJSON().Candidate,
	}

	iceBytes, err := json.Marshal(ice)
	if err != nil {
		return err
	}

	msg := P2PMessage {
		MsgType: ICE_CANDIDATE,
		Msg:     iceBytes,
	}

	err = p.signalMessage(peerId, msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *P2P) signalOffer(peerId string, offer webrtc.SessionDescription) error {
	sdp, err := json.Marshal(offer)
	if err != nil {
		return err
	}

	off := ConnOffer {
		Sender: p.peerId,
		SDP:    sdp,
	}

	offBytes, err := json.Marshal(off)
	if err != nil {
		return err
	}

	msg := P2PMessage {
		MsgType: CONN_OFFER,
		Msg:     offBytes,
	}

	err = p.signalMessage(peerId, msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *P2P) signalMessage(peerId string, msg P2PMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	fmt.Println("sending Msg to", peerId, msg)
	err = p.conn.SendPayload(peerId, payload)

	return err
}

func (p *P2P) SetupConn(peer *PeerConn, peerId string) error {
	inboundSignals := make(chan interface{}, 100)

	p.mu.Lock()
	p.msgQueues[peerId] = inboundSignals
	p.mu.Unlock()

	return p.setupConn(peer, peerId, inboundSignals)
}

func (p *P2P) setupConn(peer *PeerConn, peerId string, inboundSignals chan interface{}) error {
	peer.endpointId = peerId

	errc := make(chan error, 10)

	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	// Create a new RTCPeerConnection
	var err error
	peer.Conn, err = webrtc.NewPeerConnection(p.config)
	if err != nil {
		panic(err)
	}

	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
	peer.Conn.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peer.Conn.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else if onICECandidateErr := p.signalICECandidate(peer.endpointId, c); onICECandidateErr != nil {
			panic(onICECandidateErr)
		}
	})

	peer.OnICECandidateReceived(func (msg ICECandidateMsg) {
		fmt.Println("adding ice candidate in", p.peerId)
		if candidateErr := peer.Conn.AddICECandidate(webrtc.ICECandidateInit{Candidate: msg.IceCandidate});
			candidateErr != nil {
			panic(candidateErr)
		}
	})

	peer.OnAnswer(func (answer ConnAnswer) {
		fmt.Println(p.peerId, "received answer", answer)
		sdp := webrtc.SessionDescription{}
		if sdpErr := json.Unmarshal(answer.SDP, &sdp); sdpErr != nil {
			fmt.Println("answer: invalid sdp")
			panic(sdpErr)
		}

		if sdpErr := peer.Conn.SetRemoteDescription(sdp); sdpErr != nil {
			panic(sdpErr)
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		for _, c := range pendingCandidates {
			if onICECandidateErr := p.signalICECandidate(peer.endpointId, c); onICECandidateErr != nil {
				panic(onICECandidateErr)
			}
		}
	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peer.Conn.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("%s Peer Connection State has changed: %s\n", p.peerId, s.String())

		if s != webrtc.PeerConnectionStateConnected {
			errc <- fmt.Errorf(s.String())
		}
	})

	peer.Channel, err = peer.Conn.CreateDataChannel("data", nil)
	if err != nil {
		return err
	}

	// Register channel opening handling
	peer.Channel.OnOpen(func() {
		fmt.Println(p.peerId, "channel opened wooo")
		//fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", peer.Channel.Label(), peer.Channel.ID())
		errc <- nil
	})

	// Register text message handling
	peer.Channel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("%s Message from DataChannel '%s': '%s'\n", peer.endpointId, peer.Channel.Label(), string(msg.Data))
		peer.OnMessageCallback(msg.Data)
	})

	go peer.handleSignals(inboundSignals)

	// Create an offer to send to the other process
	offer, err := peer.Conn.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	// Note: this will start the gathering of ICE candidates
	if err = peer.Conn.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	err = p.signalOffer(peer.endpointId, offer)
	if err != nil {
		panic(err)
	}

	// Block while an error hasn't occurred
	err = <- errc
	fmt.Println("peer ", peer, "error", err)
	if err != nil {
		return fmt.Errorf("could not establish connection: %s", err)
	}

	return nil
}

func (p *P2P) connectionRequested (offer ConnOffer) {
	peer := NewConn(offer.Sender)

	inboundSignals := make(chan interface{}, 100)
	errc := make(chan error, 10)

	p.mu.Lock()
	p.msgQueues[peer.endpointId] = inboundSignals
	p.mu.Unlock()

	p.peerConnectionRequestCallback(peer, offer, nil)

	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	// Create a new RTCPeerConnection
	var err error
	peer.Conn, err = webrtc.NewPeerConnection(p.config)
	if err != nil {
		panic(err)
	}

	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
	peer.Conn.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peer.Conn.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else if onICECandidateErr := p.signalICECandidate(peer.endpointId, c); onICECandidateErr != nil {
			panic(onICECandidateErr)
		}
	})

	peer.OnICECandidateReceived(func (candidate ICECandidateMsg) {
		fmt.Println("adding ice candidate in answer")
		if candidateErr := peer.Conn.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate.IceCandidate}); candidateErr != nil {
			panic(candidateErr)
		}
	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peer.Conn.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("%s Peer Connection State has changed: %s\n", p.peerId, s.String())

		if s != webrtc.PeerConnectionStateConnected {
			errc <- fmt.Errorf(s.String())
		}
	})

	peer.Conn.OnDataChannel(func(d *webrtc.DataChannel) {
		peer.Channel = d

		d.OnOpen(func() {
			fmt.Println(p.peerId, "channel opened wooo")
			errc <- nil

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
			fmt.Printf("%s Message from DataChannel '%s': '%s'\n", peer.endpointId, peer.Channel.Label(), string(msg.Data))
			if peer.OnMessageCallback != nil {
				peer.OnMessageCallback(msg.Data)
			}
		})
	})

	sdp := webrtc.SessionDescription{}
	if err := json.Unmarshal(offer.SDP, &sdp); err != nil {
		panic(err)
	}

	if err := peer.Conn.SetRemoteDescription(sdp); err != nil {
		panic(err)
	}

	// Create an answer to send to the other process
	answer, err := peer.Conn.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	err = p.signalAnswer(offer.Sender, answer)
	if err != nil {
		panic(err)
	}
	// Sets the LocalDescription, and starts our UDP listeners
	err = peer.Conn.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	candidatesMux.Lock()
	for _, c := range pendingCandidates {
		onICECandidateErr := p.signalICECandidate(peer.endpointId, c)
		if onICECandidateErr != nil {
			panic(onICECandidateErr)
		}
	}
	candidatesMux.Unlock()

	go peer.handleSignals(inboundSignals)

	err = <- errc

	if err == nil {
		p.peerConnectionCallback(peer.endpointId, peer, nil)
	}

	return
}