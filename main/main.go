package main

import (
	"fmt"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/crdt"
)

func main() {
	server := NewServer()

	doc1 := crdt.NewSynchedDoc(1)
	doc2 := crdt.NewSynchedDoc(2)
	server.ConnectWithClient(doc1)
	server.ConnectWithClient(doc2)
	
	fmt.Printf("INserting at 0, character \"H\"\n")
	doc1.InsertAtIndex("H", 0, 1)
	
	fmt.Printf("INserting at %d, character \"e\"\n", doc2.GetLastIndex())
	doc2.InsertAtIndex("e", doc2.GetLastIndex(), 2)
	
	fmt.Printf("INserting at %d, character \"l\"\n", doc1.GetLastIndex())
	doc1.InsertAtIndex("l", doc1.GetLastIndex(), 1)
	// doc1.InsertAtIndex("l", doc1.GetLastIndex(), 1)
	// doc2.InsertAtIndex("o", doc2.GetLastIndex(), 2)
	// doc2.InsertAtIndex(" ", doc2.GetLastIndex(), 2)
	// doc2.InsertAtIndex("W", doc2.GetLastIndex(), 2)
	// doc1.InsertAtIndex("o", doc1.GetLastIndex(), 1)
	// doc1.InsertAtIndex("r", doc1.GetLastIndex(), 1)
	// doc1.InsertAtIndex("l", doc1.GetLastIndex(), 1)
	// doc1.DeleteAtIndex(doc1.GetLastIndex())
	// doc2.InsertAtIndex("d", doc2.GetLastIndex(), 2)

	time.Sleep(10 * time.Second)
	doc1.PrintDocument()
	doc2.PrintDocument()
}
