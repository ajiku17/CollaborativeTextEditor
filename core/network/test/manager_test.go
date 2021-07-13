package test

import (
	"testing"
)

// tests if the manager calls PeerConnect listener
// Connection takes place between two peers
func TestSimplePeerConnect(t *testing.T) {
}

// tests if the manager calls PeerDisconnect listener
// Connection takes place between two peers
func TestSimplePeerDisconnect(t *testing.T) {

}

// Connection takes place between multiple peers
func TestPeerConnectMultiple(t *testing.T) {

}

// Connection takes place between multiple peers
func TestDisconnectMultiple(t *testing.T) {

}

// a stress test where peers come and go rapidly
// manager should be able to handle large amount of connections
func TestPeerChurn(t *testing.T) {

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