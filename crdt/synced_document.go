package crdt

import (
	"fmt"
)

type SyncedDocument struct {
	site int
	updateManager *DocumentUpdateManager
	Document     Document
}

func NewSynchedDoc(site int) *SyncedDocument {
	serverUrl := "localhost:8081"
	doc := NewBasicDocument(NewBasicPositionManager())
	synchedDoc := SyncedDocument{site, NewDocumentUpdateManager(serverUrl), doc}
	return &synchedDoc
}

func (doc *SyncedDocument) GetSite() int {
	return doc.site
}

func (doc *SyncedDocument)InsertAtIndex(val string, index int, site int) Position {
	position := doc.Document.InsertAtIndex(val, index, doc.site)
	doc.updateManager.Insert(position, val, site)
	return position
}


func (doc *SyncedDocument) InsertAtPosition(pos Position, val string) {
	doc.Document.InsertAtPosition(pos, val)
	// TODO: send server an acknowledgement request
}

func (doc *SyncedDocument) DeleteAtIndex(index int) {
	position := doc.Document.DeleteAtIndex(index)
	doc.updateManager.Delete(position, doc.site)
}


func (doc *SyncedDocument) DeleteAtPosition(pos Position) {
	doc.Document.DeleteAtPosition(pos)
	// TODO: send server an acknowledgement request
}


func (doc *SyncedDocument)PrintDocument() {
	fmt.Printf("Document for client site N %d is : \n", doc.site)
	fmt.Println(doc.Document.ToString())
}