package crdt

import (
	"log"
	"math"
)

const base = 10

type Element struct {
	data     string
	position Position
}

type Document []Element

func NewDocument() *Document {
	doc := new(Document)
	NumberSetBase(base)
	
	doc.docInsert(0, Element{"", Position{Identifier{math.MaxInt32, 0}}})
	doc.docInsert(0, Element{"", Position{Identifier{0, 0}}})
	return doc
}

func (doc *Document) GetLength() int {
	return len(*doc) - 2
}

func (doc *Document) docInsert(index int, elem Element) {
	if index < 0 || index > len(*doc) {
		log.Fatal("Document: invalid insert index")
	}

	copyDoc := Document{}

	copyDoc = append(copyDoc, (*doc)[:index]...)
	copyDoc = append(copyDoc, elem)
	copyDoc = append(copyDoc, (*doc)[index:]...)
	
	*doc = copyDoc[:]
}

func (doc *Document) docDelete(index int) Position {
	if index < 0 || index > len(*doc) {
		log.Fatal("Document: invalid insert index")
	}

	copyDoc := Document{}

	copyDoc = append(copyDoc, (*doc)[:index]...)
	copyDoc = append(copyDoc, (*doc)[index + 1:]...)
	removedPos := (*doc)[index].position

	*doc = copyDoc[:]

	return removedPos
}

func (doc *Document) InsertAtPos(pos Position, val string) {
	var index int
	copyDoc := Document{}

	for i, e := range *doc {
		if (PositionIsLessThan(e.position, pos)) {
			index = i
		} else {
			break
		}
	}

	copyDoc = append(copyDoc, (*doc)[:index + 1]...)
	copyDoc = append(copyDoc, Element{val, pos})
	copyDoc = append(copyDoc, (*doc)[index + 1:]...)
	
	*doc = copyDoc[:]
}

func (doc *Document) DeleteAtPos(pos Position) Position {
	var index int
	copyDoc := Document{}

	for i, e := range *doc {
		if (PositionEquals(e.position, pos)) {
			index = i
			break
		}
	}

	copyDoc = append(copyDoc, (*doc)[:index]...)
	copyDoc = append(copyDoc, (*doc)[index + 1:]...)
	removedPos := (*doc)[index].position

	*doc = copyDoc[:]

	return removedPos
}

func (doc *Document) InsertAt(val string, index, site int) Position {
	if index < 0 || index > len(*doc) - 2 {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	if len(*doc) < 2 {
		log.Fatal("Document: invalid document")
	}

	prevPos := ((*doc)[index]).position
	afterPos := ((*doc)[index + 1]).position
	position := AllocPosition(prevPos, afterPos, site)
	doc.docInsert(index + 1, Element{val, position})

	return position
}

func (doc *Document) DeleteAt(index int) Position {
	if index < 0 || index > len(*doc) - 2 {
		log.Fatalf("Document: invalid delete index %v", index)
	}

	return doc.docDelete(index + 1)
}

func (elem *Element) ToString() string {
	return elem.data
}

func (doc *Document) ToString() string {
	res := ""
	for i := 0; i < len(*doc); i++ {
		res += (*doc)[i].data
	}
	return res
}
