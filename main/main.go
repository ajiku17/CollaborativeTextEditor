package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
)

func main() {

	NewServer()
	
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// doc := crdt.NewBasicDocument(crdt.NewBasicPositionManager())
	doc := synceddoc.New().(*synceddoc.SyncedDocument)
	doc1 := synceddoc.New().(*synceddoc.SyncedDocument)
	doc.Connect()
	go doc.SetChangeListener(nil)

	doc.InsertAtIndex(doc.LocalDocument.Length(), "H")
	doc.InsertAtIndex(doc.LocalDocument.Length(), "e")
	doc.InsertAtIndex(doc.LocalDocument.Length(), "l")
	doc.InsertAtIndex(doc.LocalDocument.Length(), "l")
	doc1.Connect()
	go doc1.SetChangeListener(nil)
	time.Sleep(5* time.Second)
	doc1.InsertAtIndex(doc1.LocalDocument.Length(), "o")
	doc.InsertAtIndex(doc.LocalDocument.Length(), " ")
	doc.InsertAtIndex(doc.LocalDocument.Length(), "W")
	doc1.InsertAtIndex(doc1.LocalDocument.Length(), "o")
	doc1.InsertAtIndex(doc1.LocalDocument.Length(), "r")
	doc.InsertAtIndex(doc.LocalDocument.Length(), "l")
	doc.InsertAtIndex(doc.LocalDocument.Length(), "d")
	time.Sleep(10 * time.Second)
	fmt.Printf("DOcument for id %s : %s\n", doc.GetID(), doc.LocalDocument.ToString())
	fmt.Printf("DOcument for id %s : %s\n", doc1.GetID(), doc1.LocalDocument.ToString())
}
