package main

import (
	"fmt"
	"log"

	"github.com/crdt"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	doc := new(crdt.BasicDocument)
	positionManager := new(crdt.BasicPositionManager)

	positionManager.PositionManagerInit()
	doc.DocumentInit(positionManager)

	doc.InsertAtIndex("H", 0, 1)
	doc.InsertAtIndex("e", 1, 1)
	doc.InsertAtIndex("l", 2, 1)
	doc.InsertAtIndex("l", 3, 1)
	doc.InsertAtIndex("o", 4, 1)
	doc.InsertAtIndex(" ", 5, 1)
	doc.InsertAtIndex("W", 6, 1)
	doc.InsertAtIndex("o", 7, 1)
	doc.InsertAtIndex("r", 8, 1)
	doc.InsertAtIndex("l", 9, 1)
	doc.InsertAtIndex("d", 10, 1)

	fmt.Println(doc.ToString())
}
