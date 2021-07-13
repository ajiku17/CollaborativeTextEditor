package test

import (
	"context"
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/signaling/server"
	"net/http/httptest"
	"nhooyr.io/websocket"
	"testing"
	"time"
)

func AssertTrue(t *testing.T, condition bool) {
	if !condition {
		t.Helper()
		t.Errorf("assertion failed")
	}
}

func TestServer(t *testing.T) {
	url, closeFn := setupTest(t)
	defer closeFn()

	ctx, cancel :=context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	cl, err := newClient(ctx, url, "doc", "peer1", "sdpString")
	AssertTrue(t, err == nil)

	defer cl.close()

	msg, err := cl.nextMessage()
	AssertTrue(t, err == nil)

	fmt.Println("received message", msg)

	cl2, err := newClient(ctx, url, "doc", "peer2", "sdpString2")
	AssertTrue(t, err == nil)

	defer cl2.close()

	msg2, err := cl2.nextMessage()
	AssertTrue(t, err == nil)

	fmt.Println("received message", msg2)
}

func setupTest(t *testing.T) (url string, closeFn func()) {
	s := server.NewServer()

	srv := httptest.NewServer(s)
	return srv.URL, func() {
		srv.Close()
	}
}

type client struct {
	url string
	c   *websocket.Conn
	docId string
	peerId string
	sdp string
}

func newClient(ctx context.Context, url string, docId, peerId, sdp string) (*client, error) {
	dialURL := fmt.Sprintf("%s/subscribe?doc=%s&peerId=%s&sdp=%s", url, docId, peerId, sdp)
	conn, _, err := websocket.Dial(ctx,
		dialURL,
		nil)

	if err != nil {
		return nil, err
	}

	cl := &client {
		url: url,
		c: conn,
		docId : docId,
		peerId: peerId,
		sdp: sdp,
	}

	return cl, nil
}

func (c *client) close() error {
	return c.c.Close(websocket.StatusNormalClosure, "")
}


func (c *client) nextMessage() (string, error) {
	typ, data, err := c.c.Read(context.Background())

	if err != nil {
		return "", err
	}

	if typ != websocket.MessageText {
		c.c.Close(websocket.StatusUnsupportedData, "expected text data")
		return "", fmt.Errorf("expected text message but got %v", typ)
	}

	return string(data), nil
}
