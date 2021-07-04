package test

import (
	"testing"
)

// calls InsertAt and checks if peer documents are identical.
func TestInsertAt(t *testing.T) {

}

// calls InsertAt concurrently on multiple peers and checks if peer documents are identical.
func TestInsertAtConcurrent(t *testing.T) {

}

// calls DeleteAt and checks if peer documents are identical.
func TestDeleteAt(t *testing.T) {

}

// calls DeleteAt concurrently on multiple peers and checks if peer documents are identical.
func TestDeleteAtConcurrent(t *testing.T) {

}

// calls SetCursor and checks if peer documents are identical.
func TestSetCursor(t *testing.T) {

}

// calls DeleteAt concurrently on multiple peers and checks if peer documents are identical.
func TestSetCursorConcurrent(t *testing.T) {

}

// calls InsertAt, DeleteAt and SetCursor lots of times from different peers.
// peers should eventually have identical documents.
func TestStressDocumentOperations(t *testing.T) {

}

// calls serialize on the document.
// returned value should later be deserialized into a valid document.
func TestSerialize(t *testing.T) {

}

// make changes on the document offline, and later call connect.
// peers should receive those changes after connect is called.
func TestConnect(t *testing.T) {

}

// after calling disconnect peer should not receive any more updates from other peers.
func TestDisconnect(t *testing.T) {

}

// calls connect and disconnect multiple times on a single peer.
// that peer should receive every update eventually from other peers.
func TestConnectDisconnect(t *testing.T) {

}

// calls connect and disconnect multiple times on a multiple peers.
// peers should eventually have identical documents.
func TestConnectDisconnectMultiple(t *testing.T) {

}

// calls connect and disconnect multiple times on a multiple peers with
// lots of document modifications in between. peers should eventually have identical documents.
func TestStressConnectDisconnectMultiple(t *testing.T) {

}