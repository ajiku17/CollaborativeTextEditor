package main

import (
	"fmt"

	"github.com/crdt"
)

func main() {
	fmt.Println("Hello world")

	crdt.TestPrefix()

	doc := crdt.NewDocument()

	fmt.Print(doc.ToString())
}
