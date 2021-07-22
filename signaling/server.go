package signaling

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type SignalingServer struct {

	serveMux http.ServeMux

	mu sync.Mutex
	docSubscribers  map[string] map[string]struct{} // docId -> [peer1, peer2, peer3]
	clients         map[string]client    // peerId -> [peer]

	logf func(f string, v ...interface{})
}

type client struct {
	c             *websocket.Conn
	peerId        string
	inboundMsgs   chan []byte
	outboundMsgs  chan []byte
	terminateSlow func()
}

func registerTypes() {
	gob.Register(SignalMessage{})
	gob.Register(Subscription{})
}

func NewServer() *SignalingServer {
	registerTypes()

	s := &SignalingServer{
		logf:           log.Printf,
		docSubscribers: make(map[string]map[string]struct{}),
		clients:        make(map[string]client),
	}

	s.serveMux.HandleFunc("/connect", s.connectHandler)

	return s
}

func (s *SignalingServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *SignalingServer) connectHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Printf("socket accept error %v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	arguments := r.URL.Query()
	var peerId string

	if peerId = arguments.Get("peerId"); len(peerId) == 0 {
		fmt.Printf("must provide peerId of the client")
		return
	}

	err = s.connect(c, r.Context(), peerId)
	if err != nil {
		if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
			fmt.Printf("subscribeHandler error: %v\n", err)
		}
		return
	}
}

func (s *SignalingServer) connect (c *websocket.Conn, ctx context.Context, peerId string) error {
	sub := client{
		c            : c,
		peerId       : peerId,
		outboundMsgs : make(chan []byte, 16),
		inboundMsgs  : make(chan []byte, 16),
		terminateSlow: func () {
			c.Close(websocket.StatusPolicyViolation, "connection too slow")
		},
	}

	s.putClient(sub)
	defer s.deleteClient(sub)

	errc := make(chan error, 1)

	go sub.receiver(ctx, errc)
	go sub.sender(ctx, errc)

	go s.requestProcessor(ctx, sub, errc)

	return <- errc
}

func (s *SignalingServer) requestProcessor(ctx context.Context, cl client, errc chan error) {
	for {
		msg := <- cl.inboundMsgs

		response, err := s.processRequest(ctx, cl, msg)
		if err != nil {
			errc <- err
			return
		}

		if response != nil {
			cl.outboundMsgs <- response
		}
	}
}

func (s *SignalingServer) processRequest(ctx context.Context, cl client, rqs []byte) ([]byte, error) {
	var rsp []byte

	r := bytes.NewBuffer(rqs)
	dec := gob.NewDecoder(r)

	msg := SignalMessage{}

	err := dec.Decode(&msg)
	if err != nil {
		return nil, err
	}

	if msg.MsgType == MESSAGE_SUBSCRIBE {
		//fmt.Println("decoding", msg.Msg)

		r := bytes.NewBuffer(msg.Msg)
		dec := gob.NewDecoder(r)

		subscription := Subscription{}

		err := dec.Decode(&subscription)
		if err != nil {
			return nil, fmt.Errorf("invalid subscription request: %s", err)
		}

		s.putSubscriber(subscription.DocId, cl)

		rsp, err = s.getPeerData(subscription.DocId, cl)

		if err != nil {
			return nil, fmt.Errorf("failed to send peer data {%v}", err)
		}
	} else if msg.MsgType == MESSAGE_FORWARD {
		//fmt.Println("request processor: sending ", msg.Msg, "to", msg.Receiver)
		s.sendMessageToPeer(msg.Receiver, msg.Msg)
	} else {
		return nil, fmt.Errorf("unsuported message type")
	}

	return rsp, nil
}

func (s *SignalingServer) sendMessageToPeer(peerId string, msg []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sub, ok := s.clients[peerId]; ok {
		sub.outboundMsgs <- msg
	}
}

func (s *client) receiver(ctx context.Context, errc chan error) {
	for {
		//fmt.Println("receiver trying to read")
		msg, err := read(ctx, s.c)

		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) ||
			websocket.CloseStatus(err) == websocket.StatusUnsupportedData {
			//fmt.Println("receiver continue; error", err)
			continue
		}

		if err != nil {
			if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
				fmt.Println("receiver error", err)
			}
			errc <- err
			return
		}

		if msg != nil {
			//fmt.Printf("received message: %v\n", msg)

			s.inboundMsgs <- msg
		}
	}
}

func (s *client) sender(ctx context.Context, errc chan error) {
	for {
		select {
		case msg := <- s.outboundMsgs:
			err := writeTimeout(ctx, time.Second * 5, s.c, msg)
			if err != nil {
				fmt.Println("Sender error:", err)
				errc <- err
				return
			}
		case <- ctx.Done():
			errc <- ctx.Err()
			return
		}
	}
}

func writeTimeout (ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}

func readTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	typ, data, err := c.Read(ctx)

	if err != nil {
		return nil, err
	}

	if typ != websocket.MessageText {
		c.Close(websocket.StatusUnsupportedData, "expected text data")
		return nil, fmt.Errorf("expected text message but got %v", typ)
	}

	return data, nil
}

func read(ctx context.Context, c *websocket.Conn) ([]byte, error) {
	typ, data, err := c.Read(ctx)

	if err != nil {
		return nil, err
	}

	if typ != websocket.MessageText {
		c.Close(websocket.StatusUnsupportedData, "expected text data")
		return nil, fmt.Errorf("expected text message but got %v", typ)
	}

	return data, nil
}

func (s *SignalingServer) getPeerData(docId string, sub client) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	peerIds := []string{}
	for peerId, _ := range s.docSubscribers[docId] {
		if peerId != sub.peerId {
			peerIds = append(peerIds, peerId)
		}
	}

	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	err := e.Encode(peerIds)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (s *SignalingServer) putClient (cl client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[cl.peerId] = cl
}

func (s *SignalingServer) deleteClient (sub client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, sub.peerId)
}

func (s *SignalingServer) putSubscriber (docId string, sub client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.docSubscribers[docId]; ok {
		s.docSubscribers[docId][sub.peerId] = struct{}{}
	}  else {
		s.docSubscribers[docId] = map[string] struct{} {
			sub.peerId: {},
		}
	}

	s.clients[sub.peerId] = sub
}

func (s *SignalingServer) deleteSubscriber (docId string, sub client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.docSubscribers[docId], sub.peerId)
}