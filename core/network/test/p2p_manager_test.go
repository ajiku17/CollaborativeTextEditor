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
	wg.Add(2)

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

		wg.Done()

	}()

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

		wg.Done()

	}()

	wg.Wait()

	time.Sleep(500 * time.Millisecond) // wait for change monitor

	fmt.Println(doc1.ToString(), doc2.ToString())
	AssertTrue(t, doc1.ToString() == doc2.ToString())

	wg.Add(1)
	go func () {

		doc3.LocalInsert(0, "c")
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

		wg.Done()
	}()

	wg.Wait()

	time.Sleep(1000 * time.Millisecond) // wait for change monitor

	fmt.Println(doc1.ToString(), doc2.ToString(), doc3.ToString())
	AssertTrue(t, len(doc1.ToString()) == 33)
	AssertTrue(t, len(doc1.ToString()) == len(doc2.ToString()))
	AssertTrue(t, len(doc2.ToString()) == len(doc3.ToString()))
	AssertTrue(t, doc1.ToString() == doc2.ToString())
	AssertTrue(t, doc2.ToString() == doc3.ToString())

	trackerC4 := tracker.NewClient(tUrl)
	siteId4 := "4"
	doc4, err := synceddoc.Open(siteId4, string(doc1.GetID()))
	AssertTrue(t, err == nil)
	m4 := network.NewP2PManager(utils.UUID(siteId4), doc4, sUrl, trackerC4)

	wg.Add(1)
	go func () {
		m4.Start()

		doc4.LocalInsert(0, "c")
		doc4.LocalInsert(1, "e")
		doc4.LocalInsert(2, "l")
		doc4.LocalInsert(3, "o")
		doc4.LocalInsert(4, "w")
		doc4.LocalInsert(5, "o")
		doc4.LocalInsert(6, "r")
		doc4.LocalInsert(7, "l")
		doc4.LocalInsert(8, "d")
		doc4.LocalInsert(3, "l")
		doc4.LocalInsert(5, " ")

		wg.Done()
	}()

	wg.Wait()

	time.Sleep(1500 * time.Millisecond)

	fmt.Println("doc1", doc1.ToString(), "doc2", doc2.ToString(), "doc3", doc3.ToString(), "doc4", doc4.ToString())
	AssertTrue(t, len(doc1.ToString()) == 44)
	AssertTrue(t, len(doc1.ToString()) == len(doc2.ToString()))
	AssertTrue(t, len(doc2.ToString()) == len(doc3.ToString()))
	AssertTrue(t, doc1.ToString() == doc2.ToString())
	AssertTrue(t, doc2.ToString() == doc3.ToString())
	AssertTrue(t, doc2.ToString() == doc4.ToString())
}

func insertString(d synceddoc.Document, text string) {
	for i, s := range text {
		d.LocalInsert(i, string(s))
	}
}

func Test2Peers(t *testing.T) {
	sUrl, tUrl, closeFn := setupTest(t)
	defer closeFn()

	trackerC1 := tracker.NewClient(tUrl)
	siteId1 := "1"
	doc1 := synceddoc.New(siteId1)

	trackerC2 := tracker.NewClient(tUrl)
	siteId2 := "2"
	doc2, err := synceddoc.Open(siteId2, string(doc1.GetID()))
	AssertTrue(t, err == nil)

	m1 := network.NewP2PManager(utils.UUID(siteId1), doc1, sUrl, trackerC1)
	m2 := network.NewP2PManager(utils.UUID(siteId2), doc2, sUrl, trackerC2)

	var wg sync.WaitGroup
	wg.Add(2)

	text1 := "Hello world"
	text2 := "Hello indeed"

	go func () {
		m1.Start()

		insertString(doc1, text1)

		wg.Done()
	} ()

	go func () {
		m2.Start()

		insertString(doc2, text2)

		wg.Done()
	} ()

	wg.Wait()

	time.Sleep(500 * time.Millisecond) // wait for synchronization

	fmt.Println("doc1:", doc1.ToString(), "doc2:", doc2.ToString())
	AssertTrue(t, len(doc1.ToString()) == len(text1) + len(text2))
	AssertTrue(t, len(doc1.ToString()) == len(doc2.ToString()))
	AssertTrue(t, doc1.ToString() == doc2.ToString())
}

func Test3PeersLateConnect(t *testing.T) {
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

	var wg sync.WaitGroup
	wg.Add(2)

	text1 := "Hello world"
	text2 := "Hello indeed"

	go func () {
		m1.Start()

		insertString(doc1, text1)

		wg.Done()
	} ()

	go func () {
		m2.Start()

		insertString(doc2, text2)

		wg.Done()
	} ()

	wg.Wait()

	time.Sleep(1000 * time.Millisecond) // wait for synchronization

	fmt.Println("doc1:", doc1.ToString(), "doc2:", doc2.ToString())
	AssertTrue(t, len(doc1.ToString()) == len(text1) + len(text2))
	AssertTrue(t, len(doc1.ToString()) == len(doc2.ToString()))
	AssertTrue(t, doc1.ToString() == doc2.ToString())

	m3.Start()

	time.Sleep(1000 * time.Millisecond) // wait for synchronization

	fmt.Println("doc1:", doc1.ToString(), "doc2:", doc2.ToString(), "doc3:", doc3.ToString())
	AssertTrue(t, len(doc1.ToString()) == len(text1) + len(text2))
	AssertTrue(t, len(doc1.ToString()) == len(doc2.ToString()))
	AssertTrue(t, len(doc2.ToString()) == len(doc3.ToString()))
	AssertTrue(t, doc1.ToString() == doc2.ToString())
	AssertTrue(t, doc2.ToString() == doc3.ToString())
}

func TestManyPeersLateConnect(t *testing.T) {
	sUrl, tUrl, closeFn := setupTest(t)
	defer closeFn()

	trackerC1 := tracker.NewClient(tUrl)
	siteId1 := "1"
	doc1 := synceddoc.New(siteId1)

	trackerC2 := tracker.NewClient(tUrl)
	siteId2 := "2"
	doc2, err := synceddoc.Open(siteId2, string(doc1.GetID()))
	AssertTrue(t, err == nil)

	m1 := network.NewP2PManager(utils.UUID(siteId1), doc1, sUrl, trackerC1)
	m2 := network.NewP2PManager(utils.UUID(siteId2), doc2, sUrl, trackerC2)

	var wg sync.WaitGroup
	wg.Add(2)

	text1 := "Hello world"
	text2 := "Hello indeed"

	go func () {
		m1.Start()

		insertString(doc1, text1)

		wg.Done()
	} ()

	go func () {
		m2.Start()

		insertString(doc2, text2)

		wg.Done()
	} ()

	wg.Wait()

	time.Sleep(1000 * time.Millisecond) // wait for synchronization

	fmt.Println("doc1:", doc1.ToString(), "doc2:", doc2.ToString())
	AssertTrue(t, len(doc1.ToString()) == len(text1) + len(text2))
	AssertTrue(t, len(doc1.ToString()) == len(doc2.ToString()))
	AssertTrue(t, doc1.ToString() == doc2.ToString())

	trackerC3 := tracker.NewClient(tUrl)
	siteId3 := "3"
	doc3, err := synceddoc.Open(siteId3, string(doc1.GetID()))
	AssertTrue(t, err == nil)

	trackerC4 := tracker.NewClient(tUrl)
	siteId4 := "4"
	doc4, err := synceddoc.Open(siteId4, string(doc1.GetID()))
	AssertTrue(t, err == nil)

	trackerC5 := tracker.NewClient(tUrl)
	siteId5 := "5"
	doc5, err := synceddoc.Open(siteId5, string(doc1.GetID()))
	AssertTrue(t, err == nil)

	trackerC6 := tracker.NewClient(tUrl)
	siteId6 := "6"
	doc6, err := synceddoc.Open(siteId6, string(doc1.GetID()))
	AssertTrue(t, err == nil)

	m4 := network.NewP2PManager(utils.UUID(siteId4), doc4, sUrl, trackerC4)
	m5 := network.NewP2PManager(utils.UUID(siteId5), doc5, sUrl, trackerC5)
	m6 := network.NewP2PManager(utils.UUID(siteId6), doc6, sUrl, trackerC6)
	m3 := network.NewP2PManager(utils.UUID(siteId3), doc3, sUrl, trackerC3)

	text3 := "text3"
	text4 := "text4"
	text5 := "text5"
	text6 := "text6"

	insertString(doc3, text3)
	insertString(doc4, text4)

	m3.Start()
	m4.Start()

	m5.Start()
	m6.Start()

	insertString(doc5, text5)
	insertString(doc6, text6)

	time.Sleep(4000 * time.Millisecond)

	AssertTrue(t, len(doc1.ToString()) == len(text1) + len(text2) + len(text3) + len(text4) + len(text5) + len(text6))
	AssertTrue(t, doc1.ToString() == doc2.ToString())
	AssertTrue(t, doc2.ToString() == doc3.ToString())
	AssertTrue(t, doc3.ToString() == doc4.ToString())
	AssertTrue(t, doc4.ToString() == doc5.ToString())
	AssertTrue(t, doc5.ToString() == doc6.ToString())
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