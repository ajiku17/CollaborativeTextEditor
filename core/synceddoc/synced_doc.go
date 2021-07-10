package synceddoc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type LogEntry struct {
	PeerId utils.UUID

	LogState map[utils.UUID] int
}

type SyncedDocument struct {
	id     utils.UUID
	siteId utils.UUID

	cursorPosition      int
	peerCursorPositions map[utils.UUID] int

	localDocument crdt.Document

	log []LogEntry

	onChange         ChangeListener
	onPeerConnect    PeerConnectedListener
	onPeerDisconnect PeerDisconnectedListener

	killed bool
	mu     sync.Mutex
}

func (d *SyncedDocument) ConnectSignals(changeListener ChangeListener,
	peerConnectedListener PeerConnectedListener,
	peerDisconnectedListener PeerDisconnectedListener) {

	d.setListeners(changeListener, peerConnectedListener, peerDisconnectedListener)
}

func initDocState(d *SyncedDocument) {
	d.localDocument = crdt.NewBasicDocument(crdt.NewBasicPositionManager(d.siteId))
	d.cursorPosition = 0
	d.peerCursorPositions = make(map[utils.UUID]int)
	d.killed = false
}

func (d *SyncedDocument) setListeners(changeListener ChangeListener,
	peerConnectedListener PeerConnectedListener,
	peerDisconnectedListener PeerDisconnectedListener) {

	setChangeListener(d, changeListener)
	setPeerConnectedListener(d, peerConnectedListener)
	setPeerDisconnectedListener(d, peerDisconnectedListener)
}

func registerTypes() {
	gob.Register(LogEntry{})
	gob.Register(ConnectRequest{})
	gob.Register(MessageInsert{})
	gob.Register(MessageCRDTInsert{})
	gob.Register(MessageDelete{})
	gob.Register(MessageCRDTDelete{})
	gob.Register(MessagePeerCursor{})
}

func setPeerDisconnectedListener(d *SyncedDocument, listener PeerDisconnectedListener) {
	d.onPeerDisconnect = listener
}

func setPeerConnectedListener(d *SyncedDocument, listener PeerConnectedListener) {
	d.onPeerConnect = listener
}

func setChangeListener(d *SyncedDocument, listener ChangeListener) {
	d.onChange = listener
}

// New creates a new, empty document
func New(siteId string) Document {
	doc := new (SyncedDocument)

	doc.id = utils.GenerateNewUUID()
	doc.siteId = utils.UUID(siteId)
	initDocState(doc)
	registerTypes()

	return doc
}

// Open downloads a document having the specified ID
func Open(siteId string, docId string) (Document, error) {
	doc := new (SyncedDocument)

	doc.id = utils.UUID(docId)
	doc.siteId = utils.UUID(siteId)

	initDocState(doc)

	return doc, nil
}

// Load deserializes serializedData and creates a document
func Load(siteId string, serializedData []byte) (Document, error) {
	doc := new (SyncedDocument)

	doc.id = utils.GenerateNewUUID()
	doc.siteId = utils.UUID(siteId)
	initDocState(doc)

	r := bytes.NewBuffer(serializedData)
	d := gob.NewDecoder(r)

	err := d.Decode(&doc.id)
	if err != nil {
		return nil, err
	}

	err = d.Decode(&doc.log)
	if err != nil {
		return nil, err
	}

	var documentContent []byte
	err = d.Decode(&documentContent)
	if err != nil {
		return nil, err
	}

	err = doc.localDocument.Deserialize(documentContent)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (d *SyncedDocument) GetID() utils.UUID {
	return d.id
}

func (d *SyncedDocument) SetChangeListener(listener ChangeListener) {
	setChangeListener(d, listener)
}

func (d *SyncedDocument) SetPeerConnectedListener(listener PeerConnectedListener) {
	setPeerConnectedListener(d, listener)
}

func (d *SyncedDocument) SetPeerDisconnectedListener(listener PeerDisconnectedListener) {
	setPeerDisconnectedListener(d, listener)
}

func (d *SyncedDocument) Serialize() ([]byte, error) {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	err := e.Encode(d.id)
	if err != nil {
		return nil, err
	}

	err = e.Encode(d.log)
	if err != nil {
		return nil, err
	}

	documentContent, err := d.localDocument.Serialize()
	if err != nil {
		return nil, err
	}

	err = e.Encode(documentContent)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (d *SyncedDocument) LocalInsert(index int, val string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.localDocument.InsertAtIndex(val, index)

	d.incrementPeerLastLogSequence(d.siteId)
}

func (d *SyncedDocument) RemoteInsert(peerId utils.UUID, position crdt.Position, val string, aux interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()

	index := d.localDocument.InsertAtPosition(position, val)

	d.incrementPeerLastLogSequence(peerId)

	if d.onChange != nil {
		d.onChange(MESSAGE_INSERT, MessageInsert{Index: index, Value: val}, aux)
	}
}

func (d *SyncedDocument) LocalDelete(index int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.localDocument.DeleteAtIndex(index)

	d.incrementPeerLastLogSequence(d.siteId)
}

func (d *SyncedDocument) RemoteDelete(peerId utils.UUID, position crdt.Position, aux interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()

	index := d.localDocument.DeleteAtPosition(position)

	d.incrementPeerLastLogSequence(peerId)

	if d.onChange != nil {
		d.onChange(MESSAGE_DELETE, MessageDelete{Index: index}, aux)
	}
}

func (d *SyncedDocument) ApplyRemoteOp(peerId utils.UUID, op Op, aux interface{}) {
	switch op.(type) {
	case crdt.OpInsert:
		crdtOp := op.(crdt.OpInsert)
		d.RemoteInsert(peerId, crdtOp.Pos, crdtOp.Val, aux)
	case crdt.OpDelete:
		crdtOp := op.(crdt.OpDelete)
		d.RemoteDelete(peerId, crdtOp.Pos, aux)
	default:
		fmt.Println("[SyncedDoc] unknown op")
	}
}

func (d *SyncedDocument) incrementPeerLastLogSequence(peerId utils.UUID) {
	var newState map[utils.UUID] int

	if len(d.log) == 0 {
		newState = make(map[utils.UUID] int)
	} else {
		newState = d.log[len(d.log) - 1].LogState
	}

	newState[peerId]++

	d.log = append(d.log, LogEntry{PeerId: peerId, LogState: newState})
}

func (d *SyncedDocument) SetCursor(index int) {
	d.cursorPosition = index
}

func (d *SyncedDocument) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.killed = true

	// free resources
	d.localDocument = nil
}

func (d *SyncedDocument) ToString() string {
	return d.localDocument.ToString()
}

func (d *SyncedDocument) String() string {
	return "[Document " + string(d.id) + "]" + d.localDocument.ToString()
}

func (d *SyncedDocument) GetDocument() crdt.Document {
	return d.localDocument
}