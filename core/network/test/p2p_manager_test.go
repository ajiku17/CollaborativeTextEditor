package test

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/network"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/signaling"
	"github.com/ajiku17/CollaborativeTextEditor/tracker"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"net/http/httptest"
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

func TestCreate(t *testing.T) {
	sUrl, tUrl, closeFn := setupTest(t)
	defer closeFn()

	trackerC1 := tracker.NewClient(tUrl)
	siteId1 := "1"
	doc1 := synceddoc.New(siteId1)

	trackerC2 := tracker.NewClient(tUrl)
	siteId2 := "2"
	doc2, err := synceddoc.Open(siteId2, string(doc1.GetID()))
	AssertTrue(t, err == nil)

	trackerC3 := tracker.NewClient(tUrl)
	siteId3 := "3"
	doc3, err := synceddoc.Open(siteId3, string(doc1.GetID()))
	AssertTrue(t, err == nil)

	m1 := network.NewP2PManager(utils.UUID(siteId1), doc1, sUrl, trackerC1)
	m2 := network.NewP2PManager(utils.UUID(siteId2), doc2, sUrl, trackerC2)
	m3 := network.NewP2PManager(utils.UUID(siteId3), doc3, sUrl, trackerC3)

	m3.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func () {
		m1.Start()

		doc1.LocalInsert(0, "h")
		doc1.LocalInsert(1, "e")
		doc1.LocalInsert(2, "l")
		doc1.LocalInsert(3, "o")
		doc1.LocalInsert(4, "w")
		doc1.LocalInsert(5, "o")
		doc1.LocalInsert(6, "r")
		doc1.LocalInsert(7, "l")
		doc1.LocalInsert(8, "d")
		doc1.LocalInsert(3, "l")
		doc1.LocalInsert(5, " ")

		time.Sleep(3 * time.Second)

		wg.Done()

		m1.Stop()
	}()

	wg.Add(1)
	go func () {
		m2.Start()

		doc2.LocalInsert(0, "w")
		doc2.LocalInsert(1, "C")
		doc2.LocalInsert(2, "d")
		doc2.LocalInsert(3, "o")
		doc2.LocalInsert(4, "i")
		doc2.LocalInsert(5, "n")
		doc2.LocalInsert(6, "d")
		doc2.LocalInsert(7, "e")
		doc2.LocalInsert(8, "e")
		doc2.LocalInsert(3, "d")
		doc2.LocalInsert(5, " ")

		time.Sleep(3 * time.Second)

		wg.Done()

		m2.Stop()
	}()

	wg.Wait()

	fmt.Println(doc1.ToString(), doc2.ToString())
	AssertTrue(t, doc1.ToString() == doc2.ToString())

	wg.Add(1)
	go func () {
		//m3.Start()

		doc3.LocalInsert(0, "h")
		doc3.LocalInsert(1, "e")
		doc3.LocalInsert(2, "l")
		doc3.LocalInsert(3, "o")
		doc3.LocalInsert(4, "w")
		doc3.LocalInsert(5, "o")
		doc3.LocalInsert(6, "r")
		doc3.LocalInsert(7, "l")
		doc3.LocalInsert(8, "d")
		doc3.LocalInsert(3, "l")
		doc3.LocalInsert(5, " ")

		time.Sleep(3 * time.Second)

		wg.Done()

		m3.Stop()
	}()

	wg.Wait()

	fmt.Println(doc1.ToString(), doc2.ToString(), doc3.ToString())
	AssertTrue(t, doc1.ToString() == doc2.ToString())
	AssertTrue(t, doc2.ToString() == doc3.ToString())
}

func setupTest(t *testing.T) (signalingURL string, trackerURL string, closeFn func()) {
	s := signaling.NewServer()
	tr := tracker.NewHttpTracker()

	srvSignal := httptest.NewServer(s)
	srvTracker := httptest.NewServer(tr)
	return srvSignal.URL, srvTracker.URL, func() {
		srvSignal.Close()
		srvTracker.Close()
	}
}