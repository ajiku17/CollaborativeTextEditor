package Client

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/signaling"
	"nhooyr.io/websocket"
)

type Client struct {
	url    string
	c      *websocket.Conn
	docId  string
	peerId string
}

// NewClient URL template: http://[ip]:[port]
func NewClient(ctx context.Context, url string, peerId string) (*Client, error) {
	dialURL := fmt.Sprintf("%s/connect?peerId=%s", url, peerId)
	conn, _, err := websocket.Dial(ctx,
		dialURL,
		nil)

	if err != nil {
		return nil, err
	}

	cl := &Client {
		url: url,
		c: conn,
		peerId: peerId,
	}

	return cl, nil
}

func (c *Client) Close() error {
	return c.c.Close(websocket.StatusNormalClosure, "")
}

func (c *Client) NextMessage() ([]byte, error) {
	ctx := context.Background()

	typ, data, err := c.c.Read(ctx)

	if err != nil {
		return nil, err
	}

	if typ != websocket.MessageText {
		c.c.Close(websocket.StatusUnsupportedData, "expected text data")
		return nil, fmt.Errorf("expected text message but got %v", typ)
	}

	return data, nil
}

func (c *Client) SendMessage(peer string, msg string) error {
	return c.SendData(peer, signaling.MESSAGE_FORWARD, []byte(msg))
}

func (c *Client) Subscribe(docId string) error {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	subs := signaling.Subscription {
		PeerId: c.peerId,
		DocId: docId,
	}

	err := e.Encode(subs)
	if err != nil {
		return err
	}

	c.docId = docId

	return c.SendData("", signaling.MESSAGE_SUBSCRIBE, w.Bytes())
}

func (c *Client) SendData(peer string, msgType string, payload []byte) error {
	//fmt.Println("sending message", payload, "to peer", peer)
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	msgStruct := signaling.SignalMessage {
		Receiver: peer,
		MsgType: msgType,
		Msg: payload,
	}

	err := e.Encode(msgStruct)
	if err != nil {
		return err
	}

	toSend := w.Bytes()

	//fmt.Println("sending message", toSend, "to peer", peer)

	err = c.c.Write(context.Background(), websocket.MessageText, toSend)
	return err
}

