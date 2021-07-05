package synceddoc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/core/network"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"sync"
)

type SyncedDocument struct {
	id     utils.UUID
	siteId utils.UUID

	cursorPosition      int
	peerCursorPositions map[utils.UUID] int

	localDocument crdt.Document
	syncManager   network.Manager

	killed bool
	mu     sync.Mutex
}

func (doc *SyncedDocument) Connect() {
	doc.syncManager.Connect()
}

func (doc *SyncedDocument) Disconnect() {
	doc.syncManager.Disconnect()
}

func initDocState(doc *SyncedDocument) {
	doc.siteId = utils.GenerateNewID()
	doc.localDocument = crdt.NewBasicDocument(crdt.NewBasicPositionManager(doc.siteId))
	doc.cursorPosition = 0
	doc.peerCursorPositions = make(map[utils.UUID]int)
	doc.killed = false
}

func setListeners(doc *SyncedDocument, changeListener ChangeListener,
	peerConnectedListener PeerConnectedListener,
	peerDisconnectedListener PeerDisconnectedListener) {

	setChangeListener(doc, changeListener)
	setPeerConnectedListener(doc, peerConnectedListener)
	setPeerDisconnectedListener(doc, peerDisconnectedListener)
}

func setPeerDisconnectedListener(doc *SyncedDocument, listener PeerDisconnectedListener) {
	doc.syncManager.SetPeerDisconnectedListener(func (peerId utils.UUID, aux interface{}) {
		listener(peerId)
	})
}

func setPeerConnectedListener(doc *SyncedDocument, listener PeerConnectedListener) {
	doc.syncManager.SetPeerConnectedListener(func (peerId utils.UUID, aux interface{}) {
		if cursorIndex, ok := aux.(int); ok {
			listener(peerId, cursorIndex)
		} else {
			fmt.Printf("Peer connected listener: received invalid argument of type %T", aux)
		}
	})
}

func setChangeListener(doc *SyncedDocument, listener ChangeListener) {
	doc.syncManager.SetOnMessageReceiveListener(func (message interface{}) {
		switch message.(type) {
		case ChangeCRDTInsert:
			change := message.(ChangeCRDTInsert)
			insertIndex := doc.localDocument.InsertAtPosition(change.Position, change.Value)

			listener(CHANGE_INSERT, ChangeInsert {Index: insertIndex, Value: change.Value})
		case ChangeCRDTDelete:
			change := message.(ChangeCRDTDelete)
			deleteIndex := doc.localDocument.DeleteAtPosition(change.Position)

			listener(CHANGE_DELETE, ChangeDelete {Index: deleteIndex})
		case ChangePeerCursor:
			listener(CHANGE_PEER_CURSOR, message)
		}
	})
}

// New creates a new, empty document
func New(syncManager network.Manager, changeListener ChangeListener,
	peerConnectedListener PeerConnectedListener,
	peerDisconnectedListener PeerDisconnectedListener) Document {

	syncedDoc := new (SyncedDocument)

	syncedDoc.syncManager = syncManager
	initDocState(syncedDoc)
	setListeners(syncedDoc, changeListener, peerConnectedListener, peerDisconnectedListener)

	return syncedDoc
}

// Open downloads a document having the specified ID
func Open(docId string, syncManager network.Manager, changeListener ChangeListener,
	peerConnectedListener PeerConnectedListener,
	peerDisconnectedListener PeerDisconnectedListener) (Document, error) {

	syncedDoc := new (SyncedDocument)

	syncedDoc.syncManager = syncManager
	initDocState(syncedDoc)
	setListeners(syncedDoc, changeListener, peerConnectedListener, peerDisconnectedListener)

	syncedDoc.Connect()

	return syncedDoc, nil
}

// Load deserializes serializedData and creates a document
func Load(serializedData []byte, syncManager network.Manager, changeListener ChangeListener,
	peerConnectedListener PeerConnectedListener,
	peerDisconnectedListener PeerDisconnectedListener) (Document, error) {

	syncedDoc := new (SyncedDocument)

	r := bytes.NewBuffer(serializedData)
	d := gob.NewDecoder(r)

	err := d.Decode(&syncedDoc.id)
	if err != nil {
		return nil, err
	}

	var documentContent []byte
	err = d.Decode(&documentContent)
	if err != nil {
		return nil, err
	}

	err = syncedDoc.localDocument.Deserialize(documentContent)
	if err != nil {
		return nil, err
	}

	syncedDoc.syncManager = syncManager
	initDocState(syncedDoc)
	setListeners(syncedDoc, changeListener, peerConnectedListener, peerDisconnectedListener)

	syncedDoc.Connect()

	return syncedDoc, nil
}

func (doc *SyncedDocument) GetID() utils.UUID {
	return doc.id
}

func (doc *SyncedDocument) SetChangeListener(listener ChangeListener) {
	setChangeListener(doc, listener)
}

func (doc *SyncedDocument) SetPeerConnectedListener(listener PeerConnectedListener) {
	setPeerConnectedListener(doc, listener)
}

func (doc *SyncedDocument) SetPeerDisconnectedListener(listener PeerDisconnectedListener) {
	setPeerDisconnectedListener(doc, listener)
}

func (doc *SyncedDocument) Serialize() ([]byte, error) {
	var result []byte

	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	err := e.Encode(doc.id)
	if err != nil {
		return nil, err
	}

	documentContent, err := doc.localDocument.Serialize()
	if err != nil {
		return nil, err
	}

	result = append(result, w.Bytes()...)
	result = append(result, documentContent...)

	return result, nil
}

func (doc *SyncedDocument) InsertAtIndex(index int, val string) {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	pos := doc.localDocument.InsertAtIndex(val, index)

	doc.syncManager.BroadcastMessage(ChangeCRDTInsert {Position: pos, Value: val})
}

func (doc *SyncedDocument) DeleteAtIndex(index int) {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	pos := doc.localDocument.DeleteAtIndex(index)

	doc.syncManager.BroadcastMessage(ChangeCRDTDelete{Position: pos})
}

func (doc *SyncedDocument) SetCursor(index int) {
	doc.cursorPosition = index
}

func (doc *SyncedDocument) Close() {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	doc.killed = true

	// free resources
	doc.localDocument = nil
}

func (doc *SyncedDocument) ToString() string {
	return "[Document " + string(doc.id) + "]" + doc.localDocument.ToString()
}