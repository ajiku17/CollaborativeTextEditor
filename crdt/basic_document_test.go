package crdt

import (
	"testing"
)

func TestDocInsert(t *testing.T) {
	manager := new(BasicPositionManager)
	document := new(BasicDocument)

	// Insert at the beginning
	document.docInsert(0, Element{"data", nil})
	AssertTrue(t, len(document.elems) == 1)

	document.docInsert(1, Element{"data1", BasicPosition{Identifier{0, 5}}})
	AssertTrue(t, len(document.elems) == 2)
	AssertTrue(t, manager.PositionsEqual(document.elems[1].position, BasicPosition{Identifier{0, 5}}))

	// Insert in the middle
	document.docInsert(1, Element{"data2", BasicPosition{Identifier{0, 7}}})
	AssertTrue(t, len(document.elems) == 3)
	AssertTrue(t, manager.PositionsEqual(document.elems[1].position, BasicPosition{Identifier{0, 7}}))
	AssertTrue(t, manager.PositionsEqual(document.elems[2].position, BasicPosition{Identifier{0, 5}}))

	// Insert at the end
	document.docInsert(3, Element{"end", BasicPosition{Identifier{0, 7}}})
	AssertTrue(t, len(document.elems) == 4)
	AssertTrue(t, (document.elems)[3].data == "end")
}

func TestDocDelete(t *testing.T) {
	manager := new(BasicPositionManager)
	document := new(BasicDocument)

	// Insert at the beginning
	document.docInsert(0, Element{"data", nil})
	AssertTrue(t, len(document.elems) == 1)

	document.docInsert(1, Element{"data1", BasicPosition{Identifier{0, 5}}})
	AssertTrue(t, len(document.elems) == 2)
	AssertTrue(t, manager.PositionsEqual(document.elems[1].position, BasicPosition{Identifier{0, 5}}))

	// Insert in the middle
	document.docInsert(1, Element{"data2", BasicPosition{Identifier{0, 6}}})
	AssertTrue(t, len(document.elems) == 3)
	AssertTrue(t, manager.PositionsEqual(document.elems[1].position, BasicPosition{Identifier{0, 6}}))
	AssertTrue(t, manager.PositionsEqual(document.elems[2].position, BasicPosition{Identifier{0, 5}}))

	// Insert at the end
	document.docInsert(3, Element{"end", BasicPosition{Identifier{0, 7}}})
	AssertTrue(t, len(document.elems) == 4)
	AssertTrue(t, (document.elems)[3].data == "end")

	document.docDelete(2);
	AssertTrue(t, len(document.elems) == 3)
	AssertTrue(t, (document.elems)[2].data == "end")
	AssertTrue(t, manager.PositionsEqual(document.elems[2].position, BasicPosition{Identifier{0, 7}}))

	document.docDelete(0);
	AssertTrue(t, len(document.elems) == 2)
	AssertTrue(t, (document.elems)[0].data == "data2")
	AssertTrue(t, manager.PositionsEqual(document.elems[0].position, BasicPosition{Identifier{0, 6}}))

	document.docDelete(1);
	AssertTrue(t, len(document.elems) == 1)
	AssertTrue(t, (document.elems)[0].data == "data2")
	AssertTrue(t, manager.PositionsEqual(document.elems[0].position, BasicPosition{Identifier{0, 6}}))

	document.docDelete(0);
	AssertTrue(t, len(document.elems) == 0)
}

func TestDocInit(t *testing.T) {
	manager := new(BasicPositionManager)
	document := new(BasicDocument)

	document.DocumentInit(manager)

	AssertTrue(t, len(document.elems) == 2)
	AssertTrue(t, manager.PositionsEqual(document.elems[0].position, manager.GetMinPosition()))
	AssertTrue(t, manager.PositionsEqual(document.elems[1].position, manager.GetMaxPosition()))
}