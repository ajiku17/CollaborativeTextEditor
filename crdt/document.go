package crdt

import (
	"math/rand"
)

type Element struct {
	data     string
	position Position
}

type Document []Element

func NewDocument() *Document {
	doc := new(Document)

	return doc
}


func randBetween(low, high int) int {
	return rand.Intn(high-low) + low
}

func min(x, y int) int {
	if x > y {
		return y
	}

	return x
}

func (doc *Document) InsertAt(val string, index int) {
	prevPos := ((*doc)[index]).position
	afterPos := ((*doc)[index + 1]).position
	position := AllocPosition(prevPos, afterPos)
	element := Element{val, position}
	(*doc)[index] = element
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
