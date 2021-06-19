package crdt

import (
	"fmt"
	"net/http"
	"time"
)

type SyncedDocument struct {
	site int
	updateManager *DocumentUpdateManager
	Document     Document
}

func NewSynchedDoc(site int) *SyncedDocument {
	serverUrl := "localhost:8081"
	doc := NewBasicDocument(NewBasicPositionManager())
	//Get remote document

	synchedDoc := SyncedDocument{site, &DocumentUpdateManager{serverUrl, &http.Client{Timeout: 5 * time.Minute}}, doc}
	return &synchedDoc
}

func (doc *SyncedDocument) GetSite() int {
	return doc.site
}

func (doc *SyncedDocument)InsertAtIndex(val string, index int, site int) Position {
	position := doc.Document.InsertAtIndex(val, index, doc.site)
	doc.updateManager.Insert(position, val, site)
	// doc.RemoteDocument.InsertAtPosition(position, val)
	return position
}

// func (client *Client) Insert(val string, index int) {
// 	position := client.document.InsertAtIndex(val, index, client.site)
// 	client.clientServer.SendInsertRequest(position, val, client.site)
// }

func (doc *SyncedDocument) InsertAtPosition(pos Position, val string) {
	doc.Document.InsertAtPosition(pos, val)
	// TODO: send server an acknowledgement request
}

func (doc *SyncedDocument) DeleteAtIndex(index int) {
	position := doc.Document.DeleteAtIndex(index)
	doc.updateManager.Delete(position, doc.site)
	// doc.RemoteDocument.DeleteAtPosition(position)
}

// func (client *Client) Delete(index int) {
// 	position := client.document.DeleteAtIndex(index)
// 	client.clientServer.SendDeleteRequest(position, client.site)
// }

func (doc *SyncedDocument) DeleteAtPosition(pos Position) {
	doc.Document.DeleteAtPosition(pos)
	// TODO: send server an acknowledgement request
}


func (doc *SyncedDocument)PrintDocument() {
	fmt.Printf("Document for client site N %d is : \n", doc.site)
	fmt.Println(doc.Document.ToString())
}