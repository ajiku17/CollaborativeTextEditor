package crdt

import (
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

func (doc *BasicDocument) Length() int {
	return len(doc.Elems) - 2
}

func NewBasicDocument(positionManager PositionManager) *BasicDocument {
	doc := new(BasicDocument)
	doc.Elems = []Element{}
	doc.PositionManager = positionManager

	doc.DocInsert(0, Element{"", doc.PositionManager.GetMaxPosition()})
	doc.DocInsert(0, Element{"", doc.PositionManager.GetMinPosition()})

	return doc
}


func (doc *BasicDocument)GetInsertPosition(index int, site int) Position {
	// if index < 0 || index > len(doc.Elems) - 2 {
	// 	log.Fatalf("Document: invalid insert index %v", index)
	// }

	// if len(doc.Elems) < 2 {
	// 	log.Fatal("Document: invalid document")
	// }

	prevPos := (doc.Elems[index]).Position
	afterPos := (doc.Elems[index + 1]).Position
	position := doc.PositionManager.AllocPositionBetween(prevPos, afterPos, site)	
	return position
}

func (doc *BasicDocument)GetDeletePosition(index int) Position {
	// if index < 0 || index > len(doc.Elems) - 2 {
	// 	log.Fatalf("Document: invalid delete index %v", index)
	// }

	res := doc.Elems[index + 1].Position
	return res
}

func (doc *BasicDocument) InsertAtIndex(val string, index, site int) Position {
	// if index < 0 || index > len(doc.Elems) - 2 {
	// 	log.Fatalf("Document: invalid insert index %v", index)
	// }

	// if len(doc.Elems) < 2 {
	// 	log.Fatal("Document: invalid document")
	// }

	// prevPos := (doc.Elems[index]).Position
	// afterPos := (doc.Elems[index + 1]).Position
	// position := doc.PositionManager.AllocPositionBetween(prevPos, afterPos, site)
	position := doc.GetInsertPosition(index, site)
	doc.DocInsert(index + 1, Element{val, position})

	return position
}

func (doc *BasicDocument) DeleteAtIndex(index int) Position {
	// if index < 0 || index > len(doc.Elems) - 2 {
	// 	log.Fatalf("Document: invalid delete index %v", index)
	// }

	// res := doc.Elems[index + 1].Position
	res := doc.GetDeletePosition(index)
	doc.DocDelete(index + 1)
	return res
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

func (doc *BasicDocument) InsertAtPosition(pos Position, val string) {
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
}

func (doc *BasicDocument) DeleteAtPosition (pos Position) {
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
}