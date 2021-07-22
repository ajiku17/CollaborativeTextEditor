package test

import (
	"crypto/rand"
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/p2p"
	"github.com/ajiku17/CollaborativeTextEditor/signaling"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const STUN_URL = "stun:stun.l.google.com:19302"

func AssertTrue(t *testing.T, condition bool) {
	if !condition {
		t.Helper()
		t.Errorf("assertion failed")
	}
}

func TestConnect(t *testing.T) {
	url, closeFn := setupTest(t)
	defer closeFn()

	p1 := p2p.New(url, "sandro", STUN_URL)
	defer p1.Stop()

	p2 := p2p.New(url, "tamo", STUN_URL)
	defer p2.Stop()

	connected := false

	p2.OnPeerConnection(func(s string, conn *p2p.PeerConn, aux interface{}) {
		fmt.Println(s, "connected!")
		connected = true
	})
	
	p2.OnPeerConnectionRequest(func(conn *p2p.PeerConn, offer p2p.ConnOffer, aux interface{}) {
		conn.OnMessage(func (msg []byte) {
			fmt.Println("tamo received a message", msg)
		})
	})
	
	err := p1.Start()
	AssertTrue(t, err == nil)

	err = p2.Start()
	AssertTrue(t, err == nil)

	peer1 := p2p.NewConn("tamo")

	peer1.OnMessage(func (msg []byte) {
		fmt.Println("Halooo from message", msg)
	})

	err = p1.SetupConn(peer1, "tamo")
	AssertTrue(t, err == nil)

	time.Sleep(500 * time.Millisecond)

	err = peer1.SendMessage([]byte("hello from test"))
	AssertTrue(t, err == nil)

	time.Sleep(500 * time.Millisecond)

	AssertTrue(t, connected)
}

func TestMultiplePeerMessageExchange(t *testing.T) {
	url, closeFn := setupTest(t)
	defer closeFn()

	nClients := 25
	nMessages := 128
	messageLength := 16

	conns := make(map[*p2p.P2P] map[string] *p2p.PeerConn)

	peers := []*p2p.P2P{}
	p2pMessages := map[*p2p.P2P]map[string]map[string]struct{}{}

	receivedMessages := map[string][]string{}

	mu := sync.Mutex{}

	connectedCounter := 0
	for i := 1; i <= nClients; i++ {
		p := p2p.New(url, "peer" + strconv.Itoa(i), STUN_URL)

		p.OnPeerConnectionRequest(func(conn *p2p.PeerConn, offer p2p.ConnOffer, aux interface{}) {
			//fmt.Printf("%s peer connection request. endpoint: %s address: %p\n", p.GetPeerId(), conn.GetEndpoint(), conn)
			mu.Lock()
			receivedMessages[p.GetPeerId()] = []string{}
			mu.Unlock()

			conn.OnMessage(func (msg []byte) {
				//fmt.Printf("%s %p received a message from %s %p %s\n", p.GetPeerId(), p, conn.GetEndpoint(), conn, string(msg))
				mu.Lock()
				receivedMessages[p.GetPeerId()] = append(receivedMessages[p.GetPeerId()], string(msg))
				mu.Unlock()
			})

			connectedCounter++

			if _, ok := conns[p]; ok {
				conns[p][conn.GetEndpoint()] = conn
			} else {
				m := make(map[string] *p2p.PeerConn)
				m[conn.GetEndpoint()] = conn
				conns[p] = m
			}
		})

		p.OnPeerConnection(func(s string, conn *p2p.PeerConn, i interface{}) {
			fmt.Println(p.GetPeerId(), "received a connection from ", conn.GetEndpoint())
		})

		err := p.Start()
		AssertTrue(t, err == nil)

		peers = append(peers, p)

		p2pMessages[p] = randMessages(nClients, nMessages, messageLength)
		delete(p2pMessages[p], p.GetPeerId()) // remove messages for self
	}

	allMessages := map[string][]string{}

	for _, messages := range p2pMessages {
		for peer, peerMessages := range messages {
			if _, ok := allMessages[peer]; !ok {
				allMessages[peer] = []string{}
			}
			for m, _ := range peerMessages {
				allMessages[peer] = append(allMessages[peer], m)
			}
		}
	}

	for i, p := range peers {
		func (peer *p2p.P2P) {
			for j := i + 1; j < len(peers); j++ {
				fmt.Println("setting up connection from", peer.GetPeerId(), "to", peers[j].GetPeerId())

				conn := p2p.NewConn(peers[j].GetPeerId())

				mu.Lock()
				receivedMessages[peer.GetPeerId()] = []string{}
				mu.Unlock()

				conn.OnMessage(func(msg []byte) {
					//fmt.Printf("%s %p received a message2 from %s %p %s\n", peer.GetPeerId(), p, conn.GetEndpoint(), conn, string(msg))
					mu.Lock()
					receivedMessages[peer.GetPeerId()] = append(receivedMessages[peer.GetPeerId()], string(msg))
					mu.Unlock()
				})

				err := peer.SetupConn(conn, conn.GetEndpoint())
				AssertTrue(t, err == nil)

				if _, ok := conns[p]; ok {
					conns[p][conn.GetEndpoint()] = conn
				} else {
					m := make(map[string]*p2p.PeerConn)
					m[conn.GetEndpoint()] = conn
					conns[p] = m
				}
			}
		}(p)
	}

	time.Sleep(100 * time.Second) // long idle period. connections should still be valid
	//fmt.Println(p2pMessages)
	fmt.Println("connectedCounter", connectedCounter)
	AssertTrue(t, connectedCounter == nClients * (nClients - 1) / 2)

	wg := sync.WaitGroup{}
	sentMessages := uint32(0)
	for conn, msgs := range p2pMessages {
		wg.Add(1)
		go func (p2p *p2p.P2P, msgs map[string]map[string]struct{}) {
			for receiver, ms := range msgs {
				connections, ok := conns[p2p]
				AssertTrue(t, ok)
				//fmt.Println(p2p.GetPeerId(), "connections", connections)
				peerConn, ok := connections[receiver]
				//fmt.Println(p2p.GetPeerId(), "sending messages to", peerConn.GetEndpoint())
				AssertTrue(t, ok)
				for m := range ms {
					atomic.AddUint32(&sentMessages, 1)
					err := peerConn.SendMessage([]byte(m))
					AssertTrue(t, err == nil)
				}
			}
			wg.Done()
		}(conn, msgs)
	}

	wg.Wait()

	time.Sleep(1 * time.Second) // wait for messages to distribute
	//fmt.Println("sent messages", sentMessages)
	//fmt.Println("received messages", receivedMessages)
	//fmt.Println("all messages", allMessages)

	//verify distribution was successful
	for peerId, msgs := range allMessages {
		received, ok := receivedMessages[peerId]
		AssertTrue(t, ok)

		for _, m := range msgs {

			index := -1
			//fmt.Println("looking for ", m, "in", msgs)
			for i, r := range received {
				if r == m {
					index = i
					break
				}
			}

			AssertTrue(t, index > -1)


			var copyElems []string

			copyElems = append(copyElems, received[:index]...)
			copyElems = append(copyElems, received[index+1:]...)

			received = copyElems[:]
		}

		fmt.Println("Performing assertion")
		AssertTrue(t, len(received) == 0)
	}
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
	return s
}

func setupTest(t *testing.T) (url string, closeFn func()) {
	s := signaling.NewServer()

	srv := httptest.NewServer(s)
	return srv.URL, func() {
		srv.Close()
	}
}