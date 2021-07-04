package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/crdt"
)

func main() {
	server := NewServer()

	fmt.Printf("%s", server.ConnectedSockets)

	doc1 := crdt.NewSynchedDoc(1)
	doc2 := crdt.NewSynchedDoc(2)

	time.Sleep(5 * time.Second)
	for !server.IsConnected(strconv.Itoa(1)) || !server.IsConnected(strconv.Itoa(2)) {}
	doc1.Start()
	doc2.Start()

	fmt.Printf("INserting at 0, site 1, character \"H\"\n")
	doc1.InsertAtIndex("H", 0, 1)
	
	fmt.Printf("INserting at %d, site 2, character \"e\"\n", doc2.GetLastIndex())
	doc2.InsertAtIndex("e", doc2.GetLastIndex(), 2)
	
	fmt.Printf("INserting at %d, site 1, character \"l\"\n", doc1.GetLastIndex())
	doc1.InsertAtIndex("l", doc1.GetLastIndex(), 1)
	fmt.Printf("INserting at %d, site 1, character \"L\"\n", doc1.GetLastIndex())
	doc1.InsertAtIndex("L"	, doc1.GetLastIndex(), 1)
	fmt.Printf("INserting at %d, site 2, character \"o\"\n", doc2.GetLastIndex()) 
	doc2.InsertAtIndex("o", doc2.GetLastIndex(), 2)
	// fmt.Printf("INserting at %d, site 2, character \" \"\n", doc2.GetLastIndex())
	// doc2.InsertAtIndex(" ", doc2.GetLastIndex(), 2)
	// fmt.Printf("INserting at %d, site 2, character \"W\"\n", doc2.GetLastIndex())
	// doc2.InsertAtIndex("W", doc2.GetLastIndex(), 2)
	// fmt.Printf("INserting at %d, site 1,  character \"O\"\n", doc1.GetLastIndex())
	// doc1.InsertAtIndex("O", doc1.GetLastIndex(), 1)
	// fmt.Printf("INserting at %d, site 1, character \"r\"\n", doc1.GetLastIndex())
	// doc1.InsertAtIndex("r", doc1.GetLastIndex(), 1)
	// doc1.InsertAtIndex("l", doc1.GetLastIndex(), 1)
	// doc1.DeleteAtIndex(doc1.GetLastIndex())
	// doc2.InsertAtIndex("d", doc2.GetLastIndex(), 2)

	time.Sleep(15 * time.Second)
	doc1.PrintDocument()
	doc2.PrintDocument()
	// doc1.Close()
	// doc2.Close()
}
