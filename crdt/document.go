package crdt

import (
	"math"
)

type Element struct {
	data     string
	position Position
}

type Document []Element

func NewDocument() *Document {
	doc := new(Document)
	doc.docInsert(0, Element{"", Position{Identifier{0, math.MaxInt64}}})
	doc.docInsert(0, Element{"", Position{Identifier{0, 0}}})
	return doc
}

func (doc *Document) docInsert(index int, elem Element) {
	copyDoc := append(*doc, Element{})
	doc = new(Document)
	(*doc) = append(*doc, copyDoc[:index]...)
	(*doc) = append(*doc, elem)
	(*doc) = append(*doc, copyDoc[index:]...)
}

func (doc *Document) InsertAt(val string, index, site int) {
	prevPos := ((*doc)[index]).position
	afterPos := ((*doc)[index+1]).position
	position := AllocPosition(prevPos, afterPos, site)
	// fmt.Printf("%v", position)
	doc.docInsert(index, Element{val, position})

}

func (doc *Document) DeleteAt(index int) {

}

func (elem *Element) ToString() string {
	return elem.data
}

func (doc *Document) ToString() string {
	res := "Document : {"
	for i := 0; i < len(*doc); i++ {
		res += "[element - " + (*doc)[i].ToString() + "]"
	}
	res += "}"
	return res
}
