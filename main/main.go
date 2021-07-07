package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/core/network"
	"github.com/ajiku17/CollaborativeTextEditor/core/server"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

func main() {

	server.NewServer()
	
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	id := utils.GenerateNewID()
	manager := network.NewDocumentManager(id).(*network.DocumentManager)
	doc := synceddoc.New(manager, onChange, nil, nil).(*synceddoc.SyncedDocument)
	doc.Connect()
	doc.SetChangeListener(nil)

	id1 := utils.GenerateNewID()
	manager1 := network.NewDocumentManager(id1).(*network.DocumentManager)
	doc1 := synceddoc.New(manager1, onChange, nil, nil).(*synceddoc.SyncedDocument)
	doc1.Connect()
	doc1.SetChangeListener(nil)
	time.Sleep(5 * time.Second)

	doc.InsertAtIndex(doc.GetDocument().Length(), "H")
	doc1.InsertAtIndex(doc1.GetDocument().Length(), "e")
	doc.InsertAtIndex(doc.GetDocument().Length(), "l")
	// doc1.InsertAtIndex(doc1.GetDocument().Length(), "l")
	// doc1.InsertAtIndex(doc1.GetDocument().Length(), "o")
	// doc.InsertAtIndex(doc.GetDocument().Length(), " ")
	// doc.InsertAtIndex(doc.GetDocument().Length(), "W")
	// doc1.InsertAtIndex(doc1.GetDocument().Length(), "o")
	// doc1.InsertAtIndex(doc1.GetDocument().Length(), "r")
	// doc.InsertAtIndex(doc.GetDocument().Length(), "l")
	// doc.InsertAtIndex(doc.GetDocument().Length(), "d")
	time.Sleep(20 * time.Second)
	fmt.Printf("DOcument for id %s : %s\n", doc.GetID(), doc.GetDocument().ToString())
	fmt.Printf("DOcument for id %s : %s\n", doc1.GetID(), doc1.GetDocument().ToString())
}

func onChange(changeName string, change interface {}) {
}