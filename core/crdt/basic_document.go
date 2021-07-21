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

type OpInsert struct {
	Pos Position
	Val string
}

type OpDelete struct {
	Pos Position
}

type BasicDocument struct {
	Elems           []Element
	PositionManager PositionManager
	History         []interface{}
}

func NewBasicDocument(positionManager PositionManager) *BasicDocument {
	gob.Register(OpInsert{})
	gob.Register(OpDelete{})

	doc := new(BasicDocument)

	doc.Elems = []Element{}
	doc.PositionManager = positionManager
	doc.History = []interface{}{}

	doc.DocInsert(0, Element{"", doc.PositionManager.GetMaxPosition()})
	doc.DocInsert(0, Element{"", doc.PositionManager.GetMinPosition()})

	return doc
}

func (d *BasicDocument) Length() int {
	return len(d.Elems) - 2
}

func (d *BasicDocument) DocInsert(index int, elem Element) {
	if index < 0 || index > len(d.Elems) {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	var copyElems []Element

	copyElems = append(copyElems, d.Elems[:index]...)
	copyElems = append(copyElems, elem)
	copyElems = append(copyElems, d.Elems[index:]...)

	d.Elems = copyElems[:]
}

func (d *BasicDocument) DocDelete(index int) Position {
	if index < 0 || index > len(d.Elems) {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	var copyElems []Element

	copyElems = append(copyElems, d.Elems[:index]...)
	copyElems = append(copyElems, d.Elems[index + 1:]...)
	removedPos := d.Elems[index].Position

	d.Elems = copyElems[:]

	return removedPos
}

func (d *BasicDocument) pushBackHistory(op interface {}) {
	d.History = append(d.History, op)
}

func (d *BasicDocument) InsertAtIndex(val string, index int) Position {
	if index < 0 || index > len(d.Elems) - 2 {
		log.Fatalf("Document: invalid insert index %v", index)
	}

	if len(d.Elems) < 2 {
		log.Fatal("Document: invalid document")
	}

	prevPos := (d.Elems[index]).Position
	afterPos := (d.Elems[index + 1]).Position
	position := d.PositionManager.AllocPositionBetween(prevPos, afterPos)

	d.DocInsert(index + 1, Element{val, position})
	d.pushBackHistory(OpInsert {Pos: position, Val: val})

	return position
}

func (d *BasicDocument) DeleteAtIndex(index int) Position {
	if index < 0 || index > len(d.Elems) - 2 {
		log.Fatalf("Document: invalid delete index %v", index)
	}

	position := d.DocDelete(index + 1)
	d.pushBackHistory(OpDelete {Pos: position})

	return position
}

func (d *BasicDocument) ToString() string {
	res := ""
	for i := 0; i < len(d.Elems); i++ {
		res += d.Elems[i].Data
	}
	return res
}

func (d *BasicDocument) InsertAtPosition(pos Position, val string) int {
	var index int

	for i, e := range d.Elems {
		if d.PositionManager.PositionIsLessThan(e.Position, pos) {
			index = i
		} else {
			break
		}
	}

	d.DocInsert(index + 1, Element{Position: pos, Data: val})
	d.pushBackHistory(OpInsert {Pos: pos, Val: val})

	return index
}

func (d *BasicDocument) DeleteAtPosition (pos Position) int {
	var index int

	for i, e := range d.Elems {
		if d.PositionManager.PositionsEqual(e.Position, pos) {
			index = i
			break
		}
	}

	d.DocDelete(index)
	d.pushBackHistory(OpDelete {Pos: pos})

	return index - 1
}

func (d *BasicDocument) Serialize() ([]byte, error) {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	err := e.Encode(d.Elems)
	if err != nil {
		return nil, err
	}

	err = e.Encode(d.History)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (d *BasicDocument) Deserialize(data []byte) error {
	r := bytes.NewBuffer(data)
	dec := gob.NewDecoder(r)

	err := dec.Decode(&d.Elems)
	if err != nil {
		return err
	}

	err = dec.Decode(&d.History)
	if err != nil {
		return err
	}

	return nil
}

func (d *BasicDocument)GetNextHistoryData(index int) interface{} {
	if index >= len(d.History) {
		return nil
	}
	return d.History[index]
}

func (d *BasicDocument) GetHistory() []interface{} {
	return d.History
}