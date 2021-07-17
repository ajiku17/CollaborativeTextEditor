package test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/signaling"
	"github.com/ajiku17/CollaborativeTextEditor/signaling/server"
	"net/http/httptest"
	"nhooyr.io/websocket"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func AssertTrue(t *testing.T, condition bool) {
	if !condition {
		t.Helper()
		t.Errorf("assertion failed")
	}
}

func TestSubscribe(t *testing.T) {
	var r *bytes.Buffer
	var d *gob.Decoder
	var res []string

	url, closeFn := setupTest(t)
	defer closeFn()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	cl, err := newClient(ctx, url, "peer1")
	AssertTrue(t, err == nil)
	defer cl.close()

	err = cl.subscribe("doc")
	AssertTrue(t, err == nil)

	msg, err := cl.nextMessage()
	AssertTrue(t, err == nil)

	r = bytes.NewBuffer(msg)
	d = gob.NewDecoder(r)

	err = d.Decode(&res)
	AssertTrue(t, err == nil)
	AssertTrue(t, len(res) == 0)

	cl2, err := newClient(ctx, url, "peer2")
	AssertTrue(t, err == nil)
	defer cl2.close()

	err = cl2.subscribe("doc")
	AssertTrue(t, err == nil)

	msg, err = cl2.nextMessage()
	AssertTrue(t, err == nil)

	r = bytes.NewBuffer(msg)
	d = gob.NewDecoder(r)

	err = d.Decode(&res)
	AssertTrue(t, err == nil)
	AssertTrue(t, reflect.DeepEqual(res, []string{"peer1"}))

	err = cl2.sendMessage(cl.peerId, "hello " + cl.peerId + " from " + cl2.peerId)
	AssertTrue(t, err == nil)

	msg, err = cl.nextMessage()
	AssertTrue(t, err == nil)
	AssertTrue(t, string(msg) == "hello " + cl.peerId + " from " + cl2.peerId)
}

func TestConcurrentMessageExchange(t *testing.T) {
	t.Parallel()

	url, closeFn := setupTest(t)
	defer closeFn()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	nClients := 32
	nMessages := 128
	messageLength := 128

	clients := map[*client]struct{}{}
	clMessages := map[*client]map[string]map[string]struct{}{}

	for i := 1; i <= nClients; i++ {
		cl, err := newClient(ctx, url, "peer" + strconv.Itoa(i))
		AssertTrue(t, err == nil)
		clients[cl] = struct{}{}

		clMessages[cl] = randMessages(nClients, nMessages, messageLength)
	}

	allMessages := map[string][]string{}

	for _, messages := range clMessages {
		for peer, peerMessages := range messages {
			if _, ok := allMessages[peer]; !ok {
				allMessages[peer] = []string{}
			}
			for m, _ := range peerMessages {
				allMessages[peer] = append(allMessages[peer], m)
			}
		}
	}

	var wg sync.WaitGroup
	for cl, messages := range clMessages {
		cl := cl
		messages := messages

		go func(cl *client) {
			for peer, peerMessages := range messages {
				for m, _ := range peerMessages {
					err := cl.sendMessage(peer, m)
					AssertTrue(t, err == nil)
				}
			}
		}(cl)

		wg.Add(1)
		go func(cl *client) {
			msgs := allMessages[cl.peerId]

			for len(msgs) > 0 {
				m, err := cl.nextMessage()
				AssertTrue(t, err == nil)
				//fmt.Println(err)

				index := -1
				//fmt.Println("looking for ", m, "in", msgs)
				for i, msg := range msgs {
					if msg == string(m) {
						index = i
						break
					}
				}

				AssertTrue(t, index > -1)

				var copyElems []string

				copyElems = append(copyElems, msgs[:index]...)
				copyElems = append(copyElems, msgs[index+1:]...)

				msgs = copyElems[:]
			}

			wg.Done()
		}(cl)
	}

	wg.Wait()
}

func randMessages(nClients, nMessages, length int) map[string]map[string]struct{} {
	res := map[string]map[string]struct{}{}

	for i := 1; i <= nClients; i++ {
		peer := "peer" + strconv.Itoa(i)
		messages := map[string]struct{}{}
		for j := 1; j <= nMessages; j++ {
			m := randString(length)
			if _, ok := messages[m]; ok {
				j--
				continue
			}
			messages[m] = struct{}{}
		}
		res[peer] = messages
	}

	return res
}

// randString generates a random string with length n.
func randString(n int) string {
	b := make([]byte, n)
	_, err := rand.Reader.Read(b)
	if err != nil {
		panic(fmt.Sprintf("failed to generate rand bytes: %v", err))
	}

	s := strings.ToValidUTF8(string(b), "_")
	s = strings.ReplaceAll(s, "\x00", "_")
	if len(s) > n {
		return s[:n]
	}
	if len(s) < n {
		// Pad with =
		extra := n - len(s)
		return s + strings.Repeat("=", extra)
	}
	return "hello"
	return s
}

func setupTest(t *testing.T) (url string, closeFn func()) {
	s := server.NewServer()

	srv := httptest.NewServer(s)
	return srv.URL, func() {
		srv.Close()
	}
}

type client struct {
	url    string
	c      *websocket.Conn
	docId  string
	peerId string
}

func newClient(ctx context.Context, url string, peerId string) (*client, error) {
	dialURL := fmt.Sprintf("%s/connect?peerId=%s", url, peerId)
	conn, _, err := websocket.Dial(ctx,
		dialURL,
		nil)

	if err != nil {
		return nil, err
	}

	cl := &client {
		url: url,
		c: conn,
		peerId: peerId,
	}

	return cl, nil
}

func (c *client) close() error {
	return c.c.Close(websocket.StatusNormalClosure, "")
}

func (c *client) nextMessage() ([]byte, error) {
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

func (c *client) sendMessage(peer string, msg string) error {
	return c.sendData(peer, signaling.MESSAGE_FORWARD, []byte(msg))
}

func (c *client) subscribe(docId string) error {
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

	return c.sendData("", signaling.MESSAGE_SUBSCRIBE, w.Bytes())
}

func (c *client) sendData(peer string, msgType string, payload []byte) error {
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