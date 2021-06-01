package main

import (
	"CollaborativeTextEditor/crdt"
	"fmt"
)

func main() {
	fmt.Println("Hello world")

	crdt.TestPrefix()

	doc := crdt.NewDocument()

	fmt.Print(doc.toString())
}
