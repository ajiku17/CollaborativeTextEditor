package test

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/p2p"
	"github.com/ajiku17/CollaborativeTextEditor/signaling"
	"net/http/httptest"
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

	applyCh1 := make(chan p2p.ApplyMsg, 100)
	applyCh2 := make(chan p2p.ApplyMsg, 100)

	p1 := p2p.New(url, "sandro", STUN_URL, applyCh1)
	defer p1.Stop()

	p2 := p2p.New(url, "tamo", STUN_URL, applyCh2)
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

}

func setupTest(t *testing.T) (url string, closeFn func()) {
	s := signaling.NewServer()

	srv := httptest.NewServer(s)
	return srv.URL, func() {
		srv.Close()
	}
}