package crdt

import (
	"fmt"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type SyncedDocument struct {
	site int
	updateManager DocumentUpdateManager
	Document     Document
	toCall chan *utils.PackedDocument
}

func NewSynchedDoc(site int) *SyncedDocument {
	serverUrl := "localhost:8081"
	doc := NewBasicDocument(NewBasicPositionManager())
	synchedDoc := SyncedDocument{site, NewDocumentUpdateManager(serverUrl), doc, make(chan(*utils.PackedDocument))}
	synchedDoc.updateManager.ConnectWithServer(site)
	// go synchedDoc.handleActions()
	return &synchedDoc
}

func (doc *SyncedDocument) Start() {
	go doc.sync()
}

func (doc *SyncedDocument) Close() {
	doc.toCall <- &utils.PackedDocument{"", "", "", "Done"}
}

func (doc *SyncedDocument) GetSite() int {
	return doc.site
}

func (doc *SyncedDocument) GetLastIndex() int {
	return doc.Document.Length()
}

func (doc *SyncedDocument)InsertAtIndex(val string, index int, site int) Position {
	// position := doc.Document.GetInsertPosition(index, doc.site)
	position := doc.Document.InsertAtIndex(val, index, doc.site)
	// fmt.Printf("LAST INDEX IS _ %d", doc.lastIndex)
	// doc.toCall <- &utils.PackedDocument{strconv.Itoa(doc.site), BasicPositionToString(position.(BasicPosition)), val, "Insert"}
	doc.updateManager.Insert(position, val, site)
	return position
}


func (doc *SyncedDocument) InsertAtPosition(pos Position, val string) {
	doc.Document.InsertAtPosition(pos, val)
	// doc.toCall <- &utils.PackedDocument{strconv.Itoa(doc.site), BasicPositionToString(pos.(BasicPosition)), val, "Insert"}
	// TODO: send server an acknowledgement request
}

func (doc *SyncedDocument) DeleteAtIndex(index int) {
	// position := doc.Document.GetDeletePosition(index)
	position := doc.Document.DeleteAtIndex(index)
	// doc.toCall <- &utils.PackedDocument{strconv.Itoa(doc.site), BasicPositionToString(position.(BasicPosition)), "", "Delete"}
	doc.updateManager.Delete(position, doc.site)
}


func (doc *SyncedDocument) DeleteAtPosition(pos Position) {
	doc.Document.DeleteAtPosition(pos)
	// doc.toCall <- &utils.PackedDocument{strconv.Itoa(doc.site), BasicPositionToString(pos.(BasicPosition)), "", "Delete"}
	// TODO: send server an acknowledgement request
}


func (doc *SyncedDocument)PrintDocument() {
	fmt.Printf("Document for client site N %d is : \n", doc.site)
	fmt.Println(doc.Document.ToString())
}

func (doc *SyncedDocument)sync() {
	notified := make(map[utils.PackedDocument]struct{})
	for { 	
		packedDocument := doc.updateManager.Notify()
		// fmt.Printf("Notified channel %d, %s\n", doc.site, packedDocument)
		if packedDocument != nil {
			if _, ok := notified[*packedDocument]; ok {
				continue
			}
			notified[*packedDocument] = struct{}{}
			if packedDocument.Action == "Insert" {
				doc.InsertAtPosition(ToBasicPosition(packedDocument.Position), packedDocument.Value)
			} else {
				doc.DeleteAtPosition(ToBasicPosition(packedDocument.Position))
			}
			
			time.Sleep(time.Second)
		}
	}
}

func (doc *SyncedDocument)handleActions() {	
	for {
		packedDocument, ok := <-doc.toCall
		if ok {
			// fmt.Printf("from channel site= %d %s\n", doc.site, packedDocument)
			// if packedDocument.Action == "Done" {
			// 	return
			// } else 
			if packedDocument.Action=="Insert" {
				doc.Document.InsertAtPosition(ToBasicPosition(packedDocument.Position), packedDocument.Value)
			} else {
				doc.Document.DeleteAtPosition(ToBasicPosition(packedDocument.Position))
			}
		}
	}
}