package test

import (
	"encoding/json"
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/tracker"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func AssertTrue(t *testing.T, condition bool) {
	if !condition {
		t.Helper()
		t.Errorf("assertion failed")
	}
}

func TestRegister(t *testing.T) {

}

func setup() {

}

func teardown() {

}

func parseGet(r io.Reader) []string {
	peerList := make([]string, 0)

	rspData, err := ioutil.ReadAll(r)
	if err != nil {
		return peerList
	}

	if len(rspData) == 0 {
		fmt.Println("error while retrieving data", err)
		return peerList
	}

	err = json.Unmarshal(rspData, &peerList)
	if err != nil && err != io.EOF {
		fmt.Println("error while parsing", err)
	}

	return peerList
}

func TestGet(t *testing.T) {
	tr := tracker.NewHttpTracker()

	tr.Table.Register("doc1", "peer1")

	s := httptest.NewServer(tr)

	get, err := http.Get(s.URL + "/get?doc=doc1")
	if err != nil {
		t.Error(err)
	}

	peerList := parseGet(get.Body)

	AssertTrue(t, len(peerList) == 1)
	AssertTrue(t, peerList[0] == "peer1")

	tr.Table.Register("doc1", "peer2")
	get, err = http.Get(s.URL + "/get?doc=doc1")
	if err != nil {
		t.Error(err)
	}

	peerList = parseGet(get.Body)

	AssertTrue(t, len(peerList) == 2)
	AssertTrue(t, peerList[0] == "peer1")
	AssertTrue(t, peerList[1] == "peer2")

	s.Close()
}

func TestTracker(t *testing.T) {
	tr := tracker.NewHttpTracker()

	s := httptest.NewServer(tr)

	time.Sleep(time.Millisecond * 10)

	c := tracker.NewClient(s.URL)
	c.Register("doc2", "peer1")

	peers := c.Get("doc2")
	AssertTrue(t, len(peers) == 1)
	AssertTrue(t, peers[0] == "peer1")

	peers = c.Get("doc1")
	AssertTrue(t, len(peers) == 0)

	c.Register("doc3", "peer3")

	peers = c.Get("doc3")
	AssertTrue(t, len(peers) == 1)
	AssertTrue(t, peers[0] == "peer3")

	c.Register("doc2", "peer3")

	peers = c.Get("doc2")
	AssertTrue(t, len(peers) == 2)
	AssertTrue(t, peers[0] == "peer1")
	AssertTrue(t, peers[1] == "peer3")

	s.Close()
}