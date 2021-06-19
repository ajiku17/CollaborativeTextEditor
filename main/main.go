package main

import (
	"log"

	"github.com/ajiku17/CollaborativeTextEditor/client"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server := NewServer()
	go server.HandleRequests()

	doc1 := client.NewSynchedDoc(1)
	doc2 := client.NewSynchedDoc(2)
	server.ConnectWithClient(doc1)
	server.ConnectWithClient(doc2)

	doc1.InsertAtIndex("H", 0, doc1.Site)
	doc2.InsertAtIndex("e", 1, doc2.Site)
	doc1.InsertAtIndex("l", 2, doc1.Site)
	doc1.InsertAtIndex("l", 3, doc1.Site)
	doc2.InsertAtIndex("o", 4, doc2.Site)
	doc2.InsertAtIndex(" ", 5, doc2.Site)
	doc2.InsertAtIndex("W", 6, doc2.Site)
	doc1.InsertAtIndex("o", 7, doc1.Site)
	doc1.InsertAtIndex("r", 8, doc1.Site)
	doc1.InsertAtIndex("l", 9, doc1.Site)
	doc1.DeleteAtIndex(9)
	doc2.InsertAtIndex("d", 9, doc2.Site)

	doc1.PrintDocument()
	doc2.PrintDocument()
}
