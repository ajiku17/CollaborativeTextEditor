package client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/crdt"
)

type SynchedDocument struct {
	Site int
	RemoteDocument *RemoteDocument
	Document     crdt.Document
}

func NewSynchedDoc(site int) *SynchedDocument {
	serverUrl := "http://localhost:8081/"
	manager := new(crdt.BasicPositionManager)
	doc := new(crdt.BasicDocument)
	manager.PositionManagerInit()
	doc.DocumentInit(manager)
	synchedDoc := SynchedDocument{site, &RemoteDocument{serverUrl, &http.Client{Timeout: 5 * time.Minute}}, doc}
	return &synchedDoc
}

func (doc *SynchedDocument) GetSite() int {
	return doc.Site
}

func (doc *SynchedDocument)InsertAtIndex(val string, index int, site int) crdt.Position {
	position := doc.Document.InsertAtIndex(val, index, doc.Site)
	val += ":" + strconv.Itoa(site)
	doc.RemoteDocument.InsertAtPosition(position, val)
	return position
}

// func (client *Client) Insert(val string, index int) {
// 	position := client.document.InsertAtIndex(val, index, client.site)
// 	client.clientServer.SendInsertRequest(position, val, client.site)
// }

func (doc *SynchedDocument) InsertAtPosition(pos crdt.Position, val string) {
	doc.Document.InsertAtPosition(pos, val)
	// TODO: send server an acknowledgement request
}

func (doc *SynchedDocument) DeleteAtIndex(index int) {
	position := doc.Document.DeleteAtIndex(index)
	doc.RemoteDocument.DeleteAtPosition(position)
}

// func (client *Client) Delete(index int) {
// 	position := client.document.DeleteAtIndex(index)
// 	client.clientServer.SendDeleteRequest(position, client.site)
// }

func (doc *SynchedDocument) DeleteAtPosition(pos crdt.Position) {
	doc.Document.DeleteAtPosition(pos)
	// TODO: send server an acknowledgement request
}


func (doc *SynchedDocument)PrintDocument() {
	fmt.Printf("Document for client site N %d is : \n", doc.Site)
	fmt.Println(doc.Document.ToString())
}