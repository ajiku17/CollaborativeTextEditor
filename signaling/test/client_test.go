package test

import (
	"bytes"
	"context"
	"encoding/gob"
	Client "github.com/ajiku17/CollaborativeTextEditor/signaling"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestClientSubscribe(t *testing.T) {
	var r *bytes.Buffer
	var d *gob.Decoder
	var res []string

	url, closeFn := setupTest(t)
	defer closeFn()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	cl := Client.NewClient(ctx, url, "peer1")
	AssertTrue(t, cl != nil)

	err := cl.Dial()
	AssertTrue(t, err == nil)
	defer cl.Close()

	err = cl.Subscribe("doc")
	AssertTrue(t, err == nil)

	msg, err := cl.NextMessage()
	AssertTrue(t, err == nil)

	r = bytes.NewBuffer(msg)
	d = gob.NewDecoder(r)

	err = d.Decode(&res)
	AssertTrue(t, err == nil)
	AssertTrue(t, len(res) == 0)

	cl2 := Client.NewClient(ctx, url, "peer2")
	AssertTrue(t, cl2 != nil)

	err = cl2.Dial()
	AssertTrue(t, err == nil)
	defer cl2.Close()

	err = cl2.Subscribe("doc")
	AssertTrue(t, err == nil)

	msg, err = cl2.NextMessage()
	AssertTrue(t, err == nil)

	r = bytes.NewBuffer(msg)
	d = gob.NewDecoder(r)

	err = d.Decode(&res)
	AssertTrue(t, err == nil)
	AssertTrue(t, reflect.DeepEqual(res, []string{"peer1"}))

	err = cl2.SendMessage(cl.GetPeerId(), "hello " + cl.GetPeerId() + " from " + cl2.GetPeerId())
	AssertTrue(t, err == nil)

	msg, err = cl.NextMessage()
	AssertTrue(t, err == nil)
	AssertTrue(t, string(msg) == "hello " + cl.GetPeerId() + " from " + cl2.GetPeerId())
}

func TestClientConcurrentMessageExchange(t *testing.T) {
	t.Parallel()

	url, closeFn := setupTest(t)
	defer closeFn()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	nClients := 32
	nMessages := 128
	messageLength := 128

	clients := map[*Client.Client]struct{}{}
	clMessages := map[*Client.Client]map[string]map[string]struct{}{}

	for i := 1; i <= nClients; i++ {
		cl := Client.NewClient(ctx, url, "peer" + strconv.Itoa(i))
		AssertTrue(t, cl != nil)

		err := cl.Dial()
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

		go func(cl *Client.Client) {
			for peer, peerMessages := range messages {
				for m, _ := range peerMessages {
					err := cl.SendMessage(peer, m)
					AssertTrue(t, err == nil)
				}
			}
		}(cl)

		wg.Add(1)
		go func(cl *Client.Client) {
			msgs := allMessages[cl.GetPeerId()]

			for len(msgs) > 0 {
				m, err := cl.NextMessage()
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
