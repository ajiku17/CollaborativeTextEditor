package crdt

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Element struct {
	Data     string
	Position Position
}

type BasicDocument struct {
	Elems           []Element
	PositionManager PositionManager
}

func NewBasicDocument(positionManager PositionManager) *BasicDocument {
	doc := new(BasicDocument)

	doc.Elems = []Element{}
	doc.PositionManager = positionManager

	doc.DocInsert(0, Element{"", doc.PositionManager.GetMaxPosition()})
	doc.DocInsert(0, Element{"", doc.PositionManager.GetMinPosition()})

	return doc
}

func (doc *BasicDocument) Length() int {
	return len(doc.Elems) - 2
}

func (doc *BasicDocument) InsertAtIndex(val string, index int) Position {
	if index < 0 || index > len(doc.Elems) - 2 {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	if len(doc.Elems) < 2 {
		log.Fatal("Document: invalid document")
	}

	prevPos := (doc.Elems[index]).Position
	afterPos := (doc.Elems[index + 1]).Position
	position := doc.PositionManager.AllocPositionBetween(prevPos, afterPos)
	doc.DocInsert(index + 1, Element{val, position})

	return position
}

func (doc *BasicDocument) DeleteAtIndex(index int) Position {
	if index < 0 || index > len(doc.Elems) - 2 {
		log.Fatalf("Document: invalid delete index %v", index)
	}

	return doc.DocDelete(index + 1)
}

func (doc *BasicDocument) ToString() string {
	res := ""
	for i := 0; i < len(doc.Elems); i++ {
		res += doc.Elems[i].Data
	}
	return res
}

func (doc *BasicDocument) DocInsert(index int, elem Element) {
	if index < 0 || index > len(doc.Elems) {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	var copyElems []Element

	copyElems = append(copyElems, doc.Elems[:index]...)
	copyElems = append(copyElems, elem)
	copyElems = append(copyElems, doc.Elems[index:]...)
	
	doc.Elems = copyElems[:]
}

func (doc *BasicDocument) DocDelete(index int) Position {
	if index < 0 || index > len(doc.Elems) {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	var copyElems []Element

	copyElems = append(copyElems, doc.Elems[:index]...)
	copyElems = append(copyElems, doc.Elems[index + 1:]...)
	removedPos := doc.Elems[index].Position

	doc.Elems = copyElems[:]

	return removedPos
}

func (doc *BasicDocument) InsertAtPosition(pos Position, val string) int {
	var index int
	var copyDoc []Element

	for i, e := range doc.Elems {
		if doc.PositionManager.PositionIsLessThan(e.Position, pos) {
			index = i
		} else {
			break
		}
	}

	copyDoc = append(copyDoc, doc.Elems[:index + 1]...)
	copyDoc = append(copyDoc, Element{val, pos})
	copyDoc = append(copyDoc, doc.Elems[index + 1:]...)
	
	doc.Elems = copyDoc[:]

	return index
}

func (doc *BasicDocument) DeleteAtPosition (pos Position) int {
	var index int
	var copyDoc []Element

	for i, e := range doc.Elems {
		if doc.PositionManager.PositionsEqual(e.Position, pos) {
			index = i
			break
		}
	}

	copyDoc = append(copyDoc, doc.Elems[:index]...)
	copyDoc = append(copyDoc, doc.Elems[index + 1:]...)

	doc.Elems = copyDoc[:]

	return index
}

func (doc *BasicDocument) Serialize() ([]byte, error) {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	err := e.Encode(doc.Elems)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (doc *BasicDocument) Deserialize(data []byte) error {
	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)

	err := d.Decode(&doc.Elems)
	if err != nil {
		return err
	}

	return nil
}