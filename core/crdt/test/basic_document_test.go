package test

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"testing"
)

func TestDocInsert(t *testing.T) {
	manager := new(crdt.BasicPositionManager)
	document := new(crdt.BasicDocument)

	// Insert at the beginning
	document.DocInsert(0, crdt.Element{"data", nil})
	AssertTrue(t, len(document.Elems) == 1)

	document.DocInsert(1, crdt.Element{"data1", crdt.BasicPosition{crdt.Identifier{0, 5}}})
	AssertTrue(t, len(document.Elems) == 2)
	AssertTrue(t, manager.PositionsEqual(document.Elems[1].Position, crdt.BasicPosition{crdt.Identifier{0, 5}}))

	// Insert in the middle
	document.DocInsert(1, crdt.Element{"data2", crdt.BasicPosition{crdt.Identifier{0, 7}}})
	AssertTrue(t, len(document.Elems) == 3)
	AssertTrue(t, manager.PositionsEqual(document.Elems[1].Position, crdt.BasicPosition{crdt.Identifier{0, 7}}))
	AssertTrue(t, manager.PositionsEqual(document.Elems[2].Position, crdt.BasicPosition{crdt.Identifier{0, 5}}))

	// Insert at the end
	document.DocInsert(3, crdt.Element{"end", crdt.BasicPosition{crdt.Identifier{0, 7}}})
	AssertTrue(t, len(document.Elems) == 4)
	AssertTrue(t, (document.Elems)[3].Data == "end")
}

func TestDocDelete(t *testing.T) {
	manager := new(crdt.BasicPositionManager)
	document := new(crdt.BasicDocument)

	// Insert at the beginning
	document.DocInsert(0, crdt.Element{"data", nil})
	AssertTrue(t, len(document.Elems) == 1)

	document.DocInsert(1, crdt.Element{"data1", crdt.BasicPosition{crdt.Identifier{0, 5}}})
	AssertTrue(t, len(document.Elems) == 2)
	AssertTrue(t, manager.PositionsEqual(document.Elems[1].Position, crdt.BasicPosition{crdt.Identifier{0, 5}}))

	// Insert in the middle
	document.DocInsert(1, crdt.Element{"data2", crdt.BasicPosition{crdt.Identifier{0, 6}}})
	AssertTrue(t, len(document.Elems) == 3)
	AssertTrue(t, manager.PositionsEqual(document.Elems[1].Position, crdt.BasicPosition{crdt.Identifier{0, 6}}))
	AssertTrue(t, manager.PositionsEqual(document.Elems[2].Position, crdt.BasicPosition{crdt.Identifier{0, 5}}))

	// Insert at the end
	document.DocInsert(3, crdt.Element{"end", crdt.BasicPosition{crdt.Identifier{0, 7}}})
	AssertTrue(t, len(document.Elems) == 4)
	AssertTrue(t, (document.Elems)[3].Data == "end")

	document.DocDelete(2);
	AssertTrue(t, len(document.Elems) == 3)
	AssertTrue(t, (document.Elems)[2].Data == "end")
	AssertTrue(t, manager.PositionsEqual(document.Elems[2].Position, crdt.BasicPosition{crdt.Identifier{0, 7}}))

	document.DocDelete(0);
	AssertTrue(t, len(document.Elems) == 2)
	AssertTrue(t, (document.Elems)[0].Data == "data2")
	AssertTrue(t, manager.PositionsEqual(document.Elems[0].Position, crdt.BasicPosition{crdt.Identifier{0, 6}}))

	document.DocDelete(1);
	AssertTrue(t, len(document.Elems) == 1)
	AssertTrue(t, (document.Elems)[0].Data == "data2")
	AssertTrue(t, manager.PositionsEqual(document.Elems[0].Position, crdt.BasicPosition{crdt.Identifier{0, 6}}))

	document.DocDelete(0);
	AssertTrue(t, len(document.Elems) == 0)
}

func TestDocInit(t *testing.T) {
	manager := crdt.NewBasicPositionManager()
	document := crdt.NewBasicDocument(manager)

	AssertTrue(t, len(document.Elems) == 2)
	AssertTrue(t, manager.PositionsEqual(document.Elems[0].Position, manager.GetMinPosition()))
	AssertTrue(t, manager.PositionsEqual(document.Elems[1].Position, manager.GetMaxPosition()))
}