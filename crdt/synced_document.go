package crdt

import (
	"fmt"
	"strconv"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type SyncedDocument struct {
	site int
	updateManager DocumentUpdateManager
	Document     Document
	toCall chan *utils.PackedDocument
	lastIndex int
}

func NewSynchedDoc(site int) *SyncedDocument {
	serverUrl := "localhost:8081"
	doc := NewBasicDocument(NewBasicPositionManager())
	synchedDoc := SyncedDocument{site, NewDocumentUpdateManager(serverUrl), doc, make(chan(*utils.PackedDocument)), 0}
	go synchedDoc.sync()
	go synchedDoc.handleActions()
	return &synchedDoc
}

func (doc *SyncedDocument) Close() {
	doc.toCall <- &utils.PackedDocument{"", "", "", "Done"}
}

func (doc *SyncedDocument) GetSite() int {
	return doc.site
}

func (doc *SyncedDocument) GetLastIndex() int {
	return doc.lastIndex
}

func (doc *SyncedDocument)InsertAtIndex(val string, index int, site int) Position {
	doc.lastIndex ++ 
	position := doc.Document.GetInsertPosition(index, doc.site)
	// position := doc.Document.InsertAtIndex(val, index, doc.site)
	fmt.Printf("LAST INDEX IS _ %d", doc.lastIndex)
	doc.toCall <- &utils.PackedDocument{strconv.Itoa(doc.site), BasicPositionToString(position.(BasicPosition)), val, "Insert"}
	doc.updateManager.Insert(position, val, site)
	return position
}


func (doc *SyncedDocument) InsertAtPosition(pos Position, val string) {
	doc.lastIndex ++
	// doc.Document.InsertAtPosition(pos, val)
	doc.toCall <- &utils.PackedDocument{strconv.Itoa(doc.site), BasicPositionToString(pos.(BasicPosition)), val, "Insert"}
	// TODO: send server an acknowledgement request
}

func (doc *SyncedDocument) DeleteAtIndex(index int) {
	doc.lastIndex --
	position := doc.Document.GetDeletePosition(index)
	// position := doc.Document.DeleteAtIndex(index)
	doc.toCall <- &utils.PackedDocument{strconv.Itoa(doc.site), BasicPositionToString(position.(BasicPosition)), "", "Delete"}
	doc.updateManager.Delete(position, doc.site)
}


func (doc *SyncedDocument) DeleteAtPosition(pos Position) {
	doc.lastIndex --
	// doc.Document.DeleteAtPosition(pos)
	doc.toCall <- &utils.PackedDocument{strconv.Itoa(doc.site), BasicPositionToString(pos.(BasicPosition)), "", "Delete"}
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

func (doc *SyncedDocument)handleActions() {
	
	for {
		packedDocument, ok := <-doc.toCall
		if ok {
			fmt.Printf("FROM CHANEL _ %s", packedDocument)
			if packedDocument.Action == "Done" {
				return
			} else if packedDocument.Action=="Insert" {
				doc.Document.InsertAtPosition(ToBasicPosition(packedDocument.Position), packedDocument.Value)
			} else {
				doc.Document.DeleteAtPosition(ToBasicPosition(packedDocument.Position))
			}
		}
	}
}