package synceddoc

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/core/network"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type SyncedDocument struct {
	id utils.UUID

	cursorPosition int

	localDocument crdt.Document
	syncManager   network.Manager
}

func (syncedDoc *SyncedDocument) Connect() {

}

func (syncedDoc *SyncedDocument) Disconnect() {

}

// New creates a new, empty document
func New() Document {
	syncedDoc := new (SyncedDocument)

	return syncedDoc
}

// Open downloads a document having the specified ID
func Open(docId string) (Document, error) {
	syncedDoc := new (SyncedDocument)

	return syncedDoc, nil
}

// Load deserializes serializedData and creates a document
func Load(serializedData []byte) (Document, error) {
	syncedDoc := new (SyncedDocument)

	return syncedDoc, nil
}

func (syncedDoc *SyncedDocument) GetID() utils.UUID {
	return ""
}

func (syncedDoc *SyncedDocument) SetChangeListener(listener ChangeListener) {

}

func (syncedDoc *SyncedDocument) SetPeerConnectedListener(listener PeerConnectedListener) {

}

func (syncedDoc *SyncedDocument) SetPeerDisconnectedListener(listener PeerDisconnectedListener) {

}

func (syncedDoc *SyncedDocument) Serialize() []byte {
	return []byte{}
}

func (syncedDoc *SyncedDocument) InsertAtIndex(index int, val string) {

}

func (syncedDoc *SyncedDocument) DeleteAtIndex(index int) {

}

func (syncedDoc *SyncedDocument) SetCursor(index int) {

}

func (syncedDoc *SyncedDocument) Close() {

}

func (syncedDoc *SyncedDocument) ToString() string {
	return ""
}