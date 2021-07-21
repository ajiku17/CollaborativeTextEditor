package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/core/network"
	"github.com/ajiku17/CollaborativeTextEditor/core/server"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

// tests if the manager calls PeerConnect listener
// Connection takes place between two peers
func TestSimplePeerConnect(t *testing.T) {
	server := server.NewServer()
	_, manager := connect()

	time.Sleep(5 * time.Second)
	AssertTrue(t, server.IsConnected(manager.GetId()))
}

// tests if the manager calls PeerDisconnect listener
// Connection takes place between two peers
func TestSimplePeerDisconnect(t *testing.T) {
	server := server.NewServer()
	_, manager := connect()
	time.Sleep(5 * time.Second)
	manager.Stop()
	AssertTrue(t, !manager.(*network.NetworkClient).IsAlive())

	manager.Start()
	manager.Kill()
	time.Sleep(5 * time.Second)
	AssertTrue(t, !server.IsConnected(manager.GetId()) && !manager.(*network.NetworkClient).IsAlive())
}

// Connection takes place between multiple peers
func TestPeerConnectMultiple(t *testing.T) {
	server := server.NewServer()
	_, manager1 := connect()
	_, manager2 := connect()
	_, manager3 := connect()
	time.Sleep(5 * time.Second)

	AssertTrue(t, server.IsConnected(manager1.GetId()) &&  server.IsConnected(manager2.GetId()) &&  server.IsConnected(manager3.GetId()))
}

// Connection takes place between multiple peers
func TestDisconnectMultiple(t *testing.T) {
	server := server.NewServer()
	_, manager1 := connect()
	_, manager2 := connect()
	_, manager3 := connect()
	time.Sleep(5 * time.Second)
	manager1.Kill()
	manager2.Kill()
	manager3.Kill()
	time.Sleep(5 * time.Second)

	AssertTrue(t, !server.IsConnected(manager1.GetId()) && !manager1.(*network.NetworkClient).IsAlive())
	AssertTrue(t, !server.IsConnected(manager2.GetId()) && !manager2.(*network.NetworkClient).IsAlive())
	AssertTrue(t, !server.IsConnected(manager3.GetId()) && !manager3.(*network.NetworkClient).IsAlive())
}

// a stress test where peers come and go rapidly
// manager should be able to handle large amount of connections
func TestPeerChurn(t *testing.T) {
	server := server.NewServer()
	managers := make([]*network.Manager, 0)

	for i := 0; i < 80; i++ {
		_, manager := connect()
		managers = append(managers, &manager)
	}

	time.Sleep(5 * time.Second)
	allConnected := true

	for _, manager := range managers {
		if(!server.IsConnected((*manager).GetId())) {
			allConnected = false
		}
	}

	AssertTrue(t, allConnected)
	for _, manager := range managers {
		(*manager).Kill()
	}
}

// TestSimpleBroadcastMessage tests if the manager calls MessageReceive listener
// connection takes place between two peers
func TestSimpleBroadcastMessage(t *testing.T) {
	server.NewServer()
	document1, manager1 := connect()
	document2, manager2 := connect()
	time.Sleep(5 * time.Second)

	document1.LocalInsert(0, "H")
	time.Sleep(time.Second)
	document2.LocalInsert(document2.GetDocument().Length(), "i")
	time.Sleep(time.Second)
	document1.LocalInsert(document1.GetDocument().Length(), "!")
	time.Sleep(5 * time.Second)

	fmt.Println(document1.GetDocument().ToString())
	fmt.Println(document2.GetDocument().ToString())
	AssertTrue(t, document1.GetDocument().ToString() == document2.GetDocument().ToString())
	manager1.Kill()
	manager2.Kill()
}

// a stress test with two peers and lots of messages
// network manager should not lose any of the broadcast messages
func TestStressBroadcastMessage(t *testing.T) {
	server.NewServer()
	document1, manager1 := connect()
	document2, manager2 := connect()
	time.Sleep(2 * time.Second)

	for i := 0; i < 50; i++ {
		if i % 2 == 0 {
			document1.LocalInsert(document1.GetDocument().Length(), "a" + fmt.Sprint(i))
			time.Sleep(time.Second)
		} else {
			document2.LocalInsert(document2.GetDocument().Length(), "b" + fmt.Sprint(i))
			time.Sleep(time.Second)
		}
	}
	time.Sleep(5 * time.Second)

	fmt.Println(document1.GetDocument().ToString())
	fmt.Println(document2.GetDocument().ToString())
	AssertTrue(t, document1.GetDocument().ToString() == document2.GetDocument().ToString())
	manager1.Kill()
	manager2.Kill()
}

// a stress test with multiple peers and lots of messages
// network manager should not lose any of the broadcast messages
func TestStressBroadcastMessageMultiplePeers(t *testing.T) {
	server.NewServer()
	document1, manager1 := connect()
	document2, manager2 := connect()
	document3, manager3 := connect()
	time.Sleep(2 * time.Second)

	for i := 0; i < 50; i++ {
		if i % 3 == 0 {
			document1.LocalInsert(document1.GetDocument().Length(), "a" + fmt.Sprint(i))
			time.Sleep(time.Second)
		} else if i % 3 == 1 {
			document2.LocalInsert(document2.GetDocument().Length(), "b" + fmt.Sprint(i))
			time.Sleep(time.Second)
		} else {
			document3.LocalInsert(document3.GetDocument().Length(), "b" + fmt.Sprint(i))
			time.Sleep(time.Second)
		}
	}
	time.Sleep(5 * time.Second)

	fmt.Println(document1.GetDocument().ToString())
	fmt.Println(document2.GetDocument().ToString())
	AssertTrue(t, document1.GetDocument().ToString() == document2.GetDocument().ToString() && document1.GetDocument().ToString() == document3.GetDocument().ToString())
	manager1.Kill()
	manager2.Kill()
	manager3.Kill()
}

// network manager is told to broadcast messages while it is disconnected.
// after calling connect, manager should synchronize all of those messages.
// two peers
func TestOfflineThenConnect(t *testing.T) {
	server.NewServer()
	document1, manager1 := connect()
	time.Sleep(5 * time.Second)

	document1.LocalInsert(0, "H")
	time.Sleep(time.Second)
	document1.LocalInsert(document1.GetDocument().Length(), "i")
	time.Sleep(time.Second)

	document2, manager2 := connect()
	time.Sleep(time.Second)
	document2.LocalInsert(document2.GetDocument().Length(), "!")
	time.Sleep(time.Second)

	document1.LocalInsert(document1.GetDocument().Length(), "!")
	time.Sleep(5 * time.Second)

	AssertTrue(t, document1.GetDocument().ToString() == document2.GetDocument().ToString())
	manager1.Kill()
	manager2.Kill()
}

// connect is called with multiple peers present in the network.
// all of those peers should be up to date after a few seconds.
func TestOfflineThenConnectMultiple(t *testing.T) {
	server.NewServer()
	document1, manager1 := connect()
	time.Sleep(5 * time.Second)

	document1.LocalInsert(0, "H")
	time.Sleep(time.Second)
	document1.LocalInsert(document1.GetDocument().Length(), "i")
	time.Sleep(time.Second)

	document2, manager2 := connect()
	time.Sleep(time.Second)
	document2.LocalInsert(document2.GetDocument().Length(), "!")
	time.Sleep(time.Second)

	document1.LocalInsert(document1.GetDocument().Length(), " ")
	time.Sleep(5 * time.Second)
	document1.LocalInsert(document1.GetDocument().Length(), "A")
	time.Sleep(5 * time.Second)
	document3, manager3 := connect()
	time.Sleep(time.Second)

	document3.LocalInsert(document3.GetDocument().Length(), "l")
	time.Sleep(time.Second)
	document3.LocalInsert(document3.GetDocument().Length(), "l")
	time.Sleep(time.Second)

	AssertTrue(t, document1.GetDocument().ToString() == document2.GetDocument().ToString())
	manager1.Kill()
	manager2.Kill()
	manager3.Kill()
}

// network manager should not receive any more updates after calling disconnect
func TestSimpleDisconnect(t *testing.T) {
	server.NewServer()
	document1, manager1 := connect()
	document2, manager2 := connect()
	time.Sleep(2 * time.Second)

	document1.LocalInsert(0, "H")
	time.Sleep(time.Second)

	document2.LocalInsert(document2.GetDocument().Length(), "i")
	time.Sleep(time.Second)

	manager1.Stop()
	time.Sleep(time.Second)

	document2.LocalInsert(document2.GetDocument().Length(), "!")
	time.Sleep(5 * time.Second)

	AssertTrue(t, document1.GetDocument().ToString() != document2.GetDocument().ToString())
	manager1.Kill()
	manager2.Kill()
}

// stress test with lots of connects, disconnects and message broadcasts between two peers
func TestStressConnectDisconnect(t *testing.T) {
}

// stress test with lots of connects, disconnects and message broadcasts between multiple peers
func TestStressConnectDisconnectMultiple(t *testing.T) {

}




func AssertTrue(t *testing.T, condition bool) {
	if !condition {
		t.Helper()
		t.Errorf("assertion failed")
	}
}

func connect() (synceddoc.Document, network.Manager) {
	site := utils.GenerateNewUUID()
	doc := synceddoc.New(string(site))
	manager := network.NewDocumentManager(site, &doc)
	manager.Start()
	return doc, manager
}