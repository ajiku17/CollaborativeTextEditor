package crdt

import (
	"fmt"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type SyncedDocument struct {
	site int
	updateManager DocumentUpdateManager
	Document     Document
}

func NewSynchedDoc(site int) *SyncedDocument {
	serverUrl := "localhost:8081"
	doc := NewBasicDocument(NewBasicPositionManager())
	synchedDoc := SyncedDocument{site, NewDocumentUpdateManager(serverUrl), doc}
	go synchedDoc.sync()
	return &synchedDoc
}

func (doc *SyncedDocument) GetSite() int {
	return doc.site
}

func (doc *SyncedDocument) GetLastIndex() int {
	return (*doc).Document.Length()
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

func (doc *SyncedDocument)sync() {
	notified := make(map[utils.PackedDocument]struct{})
	// i := 50
	// for i != 0 {
	for { 	
		packedDocument := doc.updateManager.Notify()
		// fmt.Printf("Notified %s\n", packedDocument)
		if packedDocument != nil {
			if _, ok := notified[*packedDocument]; ok {
				// i--
				continue
			}
			notified[*packedDocument] = struct{}{}
			if packedDocument.Action == "Insert" {
				doc.InsertAtPosition(ToBasicPosition(packedDocument.Position), packedDocument.Value)
			} else {
				doc.DeleteAtPosition(ToBasicPosition(packedDocument.Position))
			}
		}
		// i--
	}
}