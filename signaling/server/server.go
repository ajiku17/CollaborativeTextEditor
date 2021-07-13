package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type SignalingServer struct {

	serveMux http.ServeMux

	subscribersMu sync.Mutex
	subscribers   map[string] []subscriber

	logf func(f string, v ...interface{})
}

type subscriber struct {
	c             *websocket.Conn
	peerId        string
	sdpDesc       string
	msgChan       chan []byte
	terminateSlow func()
}

func NewServer() *SignalingServer {
	s := &SignalingServer{
		logf:                    log.Printf,
		subscribers:             make(map[string][]subscriber),
	}

	s.serveMux.HandleFunc("/subscribe", s.subscribeHandler)

	return s
}

func (s *SignalingServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *SignalingServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.logf("socket accept error %v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	arguments := r.URL.Query()
	var docId, peerId string

	if docId = arguments.Get("doc"); len(docId) == 0 {
		s.logf("must provide docId to subscribe to")
		return
	}

	if peerId = arguments.Get("peerId"); len(peerId) == 0 {
		s.logf("must provide peerId of the subscriber")
		return
	}

	err = s.subscribe(c, r.Context(), docId, peerId)

	if errors.Is(err, context.Canceled) {
		return
	}

	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}

	if err != nil {
		s.logf("error: %v", err)
		return
	}
}

func (s *SignalingServer) subscribe (c *websocket.Conn, ctx context.Context, docId, peerId string) error {
	ctx = c.CloseRead(ctx)

	sub := subscriber {
		peerId : peerId,
		msgChan: make(chan []byte, 16),
		terminateSlow: func () {
			c.Close(websocket.StatusPolicyViolation, "connection too slow")
		},
	}

	s.putSubscriber(docId, sub)
	defer s.deleteSubscriber(docId)

	err := s.sendPeerData(docId, sub)
	if err != nil {
		return fmt.Errorf("failed to send peer data {%v}", err)
	}

	go s.listen(sub)

	for {
		select {
		case msg := <- sub.msgChan:
			err := writeTimeout(ctx, time.Second * 5, c, msg)
			if err != nil {
				return err
			}
		case <- ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *SignalingServer) listen(sub subscriber) {

}

func writeTimeout (ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}

func (s *SignalingServer) sendPeerData(docId string, sub subscriber) error {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	peerIds := []string{}
	for _, sb := range s.subscribers[docId] {
		if sb.peerId != sub.peerId {
			peerIds = append(peerIds, sb.peerId)
		}
	}

	payload, err := json.Marshal(peerIds)
	if err != nil {
		return err
	}

	select {
	case sub.msgChan <- payload:
	default:
		go sub.terminateSlow()
	}

	return nil
}

func (s *SignalingServer) putSubscriber (docId string, sub subscriber) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	if val, ok := s.subscribers[docId]; ok {
		s.subscribers[docId] = append(val, sub)
	}  else {
		s.subscribers[docId] = []subscriber{sub}
	}
}

func (s *SignalingServer) deleteSubscriber (docId string) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	delete(s.subscribers, docId)
}


func echo(ctx context.Context, c *websocket.Conn) error {
	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}