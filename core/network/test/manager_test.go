package test

import (
	"encoding/gob"
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
	manager := getConnectedManager()

	time.Sleep(5 * time.Second)
	AssertTrue(t, server.IsConnected(manager.GetId()))
}

// tests if the manager calls PeerDisconnect listener
// Connection takes place between two peers
func TestSimplePeerDisconnect(t *testing.T) {

}

// Connection takes place between multiple peers
func TestPeerConnectMultiple(t *testing.T) {
	server := server.NewServer()
	manager1 := getConnectedManager()
	manager2 := getConnectedManager()
	manager3 := getConnectedManager()

	time.Sleep(5 * time.Second)

	AssertTrue(t, server.IsConnected(manager1.GetId()) &&  server.IsConnected(manager2.GetId()) &&  server.IsConnected(manager3.GetId()))
}

// Connection takes place between multiple peers
func TestDisconnectMultiple(t *testing.T) {

}

// a stress test where peers come and go rapidly
// manager should be able to handle large amount of connections
func TestPeerChurn(t *testing.T) {
	server := server.NewServer()

	managers := make([]*network.DummyManager, 0)

	for i := 0; i < 80; i++ {
		manager := getConnectedManager()
		managers = append(managers, manager)
	}

	time.Sleep(5 * time.Second)
	allConnected := true
	
	for _, manager := range managers {
		if(!server.IsConnected(manager.GetId())) {
			allConnected = false
		}
	}

	AssertTrue(t, allConnected)
}

// TestSimpleBroadcastMessage tests if the manager calls MessageReceive listener
// connection takes place between two peers
func TestSimpleBroadcastMessage(t *testing.T) {
	// server := server.NewServer()
	// manager1 := getConnectedManager()
	// manager2 := getConnectedManager()	
}

// a stress test with two peers and lots of messages
// network manager should not lose any of the broadcast messages
func TestStressBroadcastMessage(t *testing.T) {

}

// a stress test with multiple peers and lots of messages
// network manager should not lose any of the broadcast messages
func TestStressBroadcastMessageMultiplePeers(t *testing.T) {

}

// network manager is told to broadcast messages while it is disconnected.
// after calling connect, manager should synchronize all of those messages.
// two peers
func TestOfflineThenConnect(t *testing.T) {

}

// connect is called with multiple peers present in the network.
// all of those peers should be up to date after a few seconds.
func TestOfflineThenConnectMultiple(t *testing.T) {

}

// network manager should not receive any more updates after calling disconnect
func TestSimpleDisconnect(t *testing.T) {

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

func getConnectedManager() *network.DummyManager{
	manager := network.NewDummyManager(utils.GenerateNewID()).(*network.DummyManager)
	connect(manager)
	return manager
}

func connect(manager *network.DummyManager) {
	gob.Register(synceddoc.ConnectRequest{})
	manager.Connect()
	manager.BroadcastMessage(synceddoc.ConnectRequest{manager.GetId()})
}