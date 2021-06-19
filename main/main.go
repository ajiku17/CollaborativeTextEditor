package main

import (
	"log"
	"sync"

	"github.com/ajiku17/CollaborativeTextEditor/crdt"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var mu sync.Mutex
	mu.Lock()
	server := NewServer()
	go server.HandleRequests(&mu)

	mu.Lock()
	
	doc1 := crdt.NewSynchedDoc(1)
	doc2 := crdt.NewSynchedDoc(2)
	server.ConnectWithClient(doc1)
	server.ConnectWithClient(doc2)

	doc1.InsertAtIndex("H", 0, 1)
	doc2.InsertAtIndex("e", 1, 2)
	// doc1.InsertAtIndex("l", 2, 1)
	// doc1.InsertAtIndex("l", 3, 1)
	// doc2.InsertAtIndex("o", 4, 2)
	// doc2.InsertAtIndex(" ", 5, 2)
	// doc2.InsertAtIndex("W", 6, 2)
	// doc1.InsertAtIndex("o", 7, 1)
	// doc1.InsertAtIndex("r", 8, 1)
	// doc1.InsertAtIndex("l", 9, 1)
	// doc1.DeleteAtIndex(9)
	// doc2.InsertAtIndex("d", 9, 2)

	doc1.PrintDocument()
	doc2.PrintDocument()
}
