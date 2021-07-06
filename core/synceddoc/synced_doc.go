package synceddoc

import (
	"log"
	"strconv"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/core/network"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type SyncedDocument struct {
	id utils.UUID

	cursorPosition int

	LocalDocument crdt.Document
	syncManager network.Manager
}

func (syncedDoc *SyncedDocument) Connect() {
	// TODO: synchronizes changes
	syncedDoc.syncManager.Connect()
}

func (syncedDoc *SyncedDocument) Disconnect() {
	panic("implement me")
}

// New creates a new, empty document
func New() Document {
	id := utils.GenerateNewID()
	syncedDoc := SyncedDocument{id, 0, crdt.NewBasicDocument(crdt.NewBasicPositionManager()), &network.DocumentManager{Id:id}}
	return &syncedDoc
}

// Open downloads a document having the specified ID
func Open(docId string) Document {
	syncedDoc := new (SyncedDocument)

	return syncedDoc
}

// Load deserializes serializedData and creates a document
func Load(serializedData []byte) Document {
	syncedDoc := new (SyncedDocument)
	syncedDoc.LocalDocument.Deserialize(serializedData)
	return syncedDoc
}

func (syncedDoc *SyncedDocument) GetID() utils.UUID {
	return syncedDoc.id
}

func onMessageReceive(message interface{}) {
	notify := message.(network.ToNotify)
	notify.ToNotifyDocuments <- notify.AddCurrent
}

func (syncedDoc *SyncedDocument) SetChangeListener(listener ChangeListener) {
	go syncedDoc.syncManager.SetOnMessageReceiveListener(onMessageReceive)

	notified := make(map[utils.PackedDocument]struct{})
	for { 	
		packedDocument := <- syncedDoc.syncManager.(*network.DocumentManager).ToNotify.ToNotifyDocuments
		if packedDocument != nil {
			if _, ok := notified[*packedDocument]; ok {
				continue
			}
			notified[*packedDocument] = struct{}{}
			syncedDoc.action(*packedDocument)
			
			time.Sleep(time.Second)
		}
	}
}

func (syncedDoc *SyncedDocument) SetPeerConnectedListener(listener PeerConnectedListener) {

}

func (syncedDoc *SyncedDocument) SetPeerDisconnectedListener(listener PeerDisconnectedListener) {

}

func (syncedDoc *SyncedDocument)action(packedDocument utils.PackedDocument) {
	index, _ := strconv.Atoi(packedDocument.Index)
	if packedDocument.Action == CHANGE_INSERT {
		syncedDoc.LocalDocument.InsertAtPosition(crdt.ToBasicPosition(packedDocument.Position), packedDocument.Value)
	} else if packedDocument.Action == CHANGE_DELETE {
		syncedDoc.LocalDocument.DeleteAtIndex(index)
	}
}

func (syncedDoc *SyncedDocument) Serialize() []byte {
	res, err := syncedDoc.LocalDocument.Serialize()
	if err != nil {
		log.Fatalf("Error while serializing data: %s", err)
	}
	return res
}

func (syncedDoc *SyncedDocument) InsertAtIndex(index int, val string) {
	id, _ := strconv.Atoi(string(syncedDoc.id))
	position := syncedDoc.LocalDocument.InsertAtIndex(val, index, id)
	syncedDoc.syncManager.BroadcastMessage(utils.PackedDocument{string(syncedDoc.id), strconv.Itoa(index), crdt.BasicPositionToString(position.(crdt.BasicPosition)), val, CHANGE_INSERT})
}

func (syncedDoc *SyncedDocument) DeleteAtIndex(index int) {
	syncedDoc.LocalDocument.DeleteAtIndex(index)
	syncedDoc.syncManager.BroadcastMessage(utils.PackedDocument{string(syncedDoc.id), strconv.Itoa(index), "", "", CHANGE_DELETE})
}

func (syncedDoc *SyncedDocument) SetCursor(index int) {
	syncedDoc.cursorPosition = index
	syncedDoc.syncManager.BroadcastMessage(utils.PackedDocument{string(syncedDoc.id), strconv.Itoa(index), "", "", CHANGE_PEER_CURSOR})
}

func (syncedDoc *SyncedDocument) Close() {

}
