package crdt

import "log"

type Element struct {
	data     string
	position Position
}

type BasicDocument struct {
	elems           []Element
	positionManager   PositionManager
}

func (doc *BasicDocument) Length() int {
	return len(doc.elems) - 2
}

func NewBasicDocument(positionManager PositionManager) *BasicDocument {
	doc := new(BasicDocument)
	doc.elems = []Element{}
	doc.positionManager = positionManager

	doc.docInsert(0, Element{"", doc.positionManager.GetMaxPosition()})
	doc.docInsert(0, Element{"", doc.positionManager.GetMinPosition()})

	return doc
}


func (doc *BasicDocument) InsertAtIndex(val string, index, site int) Position {
	if index < 0 || index > len(doc.elems) - 2 {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	if len(doc.elems) < 2 {
		log.Fatal("Document: invalid document")
	}

	prevPos := (doc.elems[index]).position
	afterPos := (doc.elems[index + 1]).position
	position := doc.positionManager.AllocPositionBetween(prevPos, afterPos, site)
	doc.docInsert(index + 1, Element{val, position})

	return position
}

func (doc *BasicDocument) DeleteAtIndex(index int) {
	if index < 0 || index > len(doc.elems) - 2 {
		log.Fatalf("Document: invalid delete index %v", index)
	}

	doc.docDelete(index + 1)
}

func (doc *BasicDocument) ToString() string {
	res := ""
	for i := 0; i < len(doc.elems); i++ {
		res += doc.elems[i].data
	}
	return res
}

func (doc *BasicDocument) docInsert(index int, elem Element) {
	if index < 0 || index > len(doc.elems) {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	copyElems := []Element{}

	copyElems = append(copyElems, doc.elems[:index]...)
	copyElems = append(copyElems, elem)
	copyElems = append(copyElems, doc.elems[index:]...)
	
	doc.elems = copyElems[:]
}

func (doc *BasicDocument) docDelete(index int) Position {
	if index < 0 || index > len(doc.elems) {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	copyElems := []Element{}

	copyElems = append(copyElems, doc.elems[:index]...)
	copyElems = append(copyElems, doc.elems[index + 1:]...)
	removedPos := doc.elems[index].position

	doc.elems = copyElems[:]

	return removedPos
}

func (doc *BasicDocument) InsertAtPosition(pos Position, val string) {
	var index int
	copyDoc := []Element{}

	for i, e := range doc.elems {
		if (doc.positionManager.PositionIsLessThan(e.position, pos)) {
			index = i
		} else {
			break
		}
	}

	copyDoc = append(copyDoc, doc.elems[:index + 1]...)
	copyDoc = append(copyDoc, Element{val, pos})
	copyDoc = append(copyDoc, doc.elems[index + 1:]...)
	
	doc.elems = copyDoc[:]
}

func (doc *BasicDocument) DeleteAtPosition (pos Position) {
	var index int
	copyDoc := []Element{}

	for i, e := range doc.elems {
		if (doc.positionManager.PositionsEqual(e.position, pos)) {
			index = i
			break
		}
	}

	copyDoc = append(copyDoc, doc.elems[:index]...)
	copyDoc = append(copyDoc, doc.elems[index + 1:]...)

	doc.elems = copyDoc[:]
}