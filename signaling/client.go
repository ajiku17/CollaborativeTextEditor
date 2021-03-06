package signaling

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"nhooyr.io/websocket"
	"time"
)

type Client struct {
	url    string
	c      *websocket.Conn
	docId  string
	peerId string
	ctx    context.Context
}

// NewClient URL template: http://[ip]:[port]
func NewClient(ctx context.Context, url string, peerId string) *Client {
	cl := &Client {
		url: url,
		c: nil,
		peerId: peerId,
		ctx: ctx,
	}

	return cl
}

func (c *Client) GetPeerId() string {
	return c.peerId
}

func (c *Client) Dial() error {
	dialURL := fmt.Sprintf("%s/connect?peerId=%s", c.url, c.peerId)
	conn, _, err := websocket.Dial(c.ctx, dialURL, nil)

	if err != nil {
		return err
	}

	c.c = conn

	return nil
}

func (c *Client) Close() error {
	if c.c != nil {
		return c.c.Close(websocket.StatusNormalClosure, "")
	}

	return nil
}

func (c *Client) NextMessageChannel() (chan []byte, chan error) {
	errc := make(chan error, 1)
	msgCh := make(chan []byte, 1)

	go func (c *Client, msgCh chan []byte, errc chan error) {
		msg, err := c.NextMessage()
		if err != nil {
			errc <- err
			return
		}

		msgCh <- msg
	} (c, msgCh, errc)

	return msgCh, errc
}

func (c *Client) NextMessage() ([]byte, error) {
	if c.c == nil {
		return nil, fmt.Errorf("client connection is nil")
	}

	ctx := context.Background()

	typ, data, err := c.c.Read(ctx)

	if err != nil {
		return nil, err
	}

	if typ != websocket.MessageBinary {
		c.c.Close(websocket.StatusUnsupportedData, "expected text data")
		return nil, fmt.Errorf("expected text message but got %v", typ)
	}

	return data, nil
}

func (c *Client) NextMessageTimeout(timeout time.Duration) ([]byte, error) {
	if c.c == nil {
		return nil, fmt.Errorf("client connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	typ, data, err := c.c.Read(ctx)

	if err != nil {
		return nil, err
	}

	if typ != websocket.MessageBinary {
		c.c.Close(websocket.StatusUnsupportedData, "expected text data")
		return nil, fmt.Errorf("expected text message but got %v", typ)
	}

	return data, nil
}

func (c *Client) SendMessage(peer string, msg string) error {
	if c.c == nil {
		return fmt.Errorf("client connection is nil")
	}

	return c.sendData(peer, MESSAGE_FORWARD, []byte(msg))
}

func (c *Client) SendPayload(peer string, msg []byte) error {
	if c.c == nil {
		return fmt.Errorf("client connection is nil")
	}

	return c.sendData(peer, MESSAGE_FORWARD, msg)
}

func (c *Client) Subscribe(docId string) error {
	if c.c == nil {
		return fmt.Errorf("client connection is nil")
	}

	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	subs := Subscription{
		PeerId: c.peerId,
		DocId: docId,
	}

	err := e.Encode(subs)
	if err != nil {
		return err
	}

	c.docId = docId

	return c.sendData("", MESSAGE_SUBSCRIBE, w.Bytes())
}

func (c *Client) sendData(peer string, msgType string, payload []byte) error {
	if c.c == nil {
		return fmt.Errorf("client connection is nil")
	}

	//fmt.Println("sending payload", payload, "to peer", peer)
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	msgStruct := SignalMessage{
		Receiver: peer,
		MsgType: msgType,
		Msg: payload,
	}

	err := e.Encode(msgStruct)
	if err != nil {
		fmt.Println("send data error:", err)
		return err
	}

	toSend := w.Bytes()

	//fmt.Println(c.peerId, "sending message", toSend, "to peer", peer)

	err = c.c.Write(context.Background(), websocket.MessageBinary, toSend)
	return err
}

