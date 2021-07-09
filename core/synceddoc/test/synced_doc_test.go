package test

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"testing"
)

func AssertTrue(t *testing.T, condition bool) {
	if !condition {
		t.Helper()
		t.Errorf("assertion failed")
	}
}

// calls InsertAt and checks if peer documents are identical.
func TestLocalInsert(t *testing.T) {
	d := synceddoc.New()

	d.LocalInsert(0, "h")
	d.LocalInsert(1, "e")
	d.LocalInsert(2, "l")
	d.LocalInsert(3, "o")
	d.LocalInsert(4, "w")
	d.LocalInsert(5, "o")
	d.LocalInsert(6, "r")
	d.LocalInsert(7, "l")
	d.LocalInsert(8, "d")
	d.LocalInsert(3, "l")
	d.LocalInsert(5, " ")

	AssertTrue(t, d.ToString() == "hello world")

	d = synceddoc.New()
	text := "hello everybody"
	for i := len(text) - 1; i >= 0; i-- {
		d.LocalInsert(0, string(text[i]))
	}

	AssertTrue(t, d.ToString() == text)

	d = synceddoc.New()
	text = "hey everyone"
	for i := 0; i < len(text); i++ {
		d.LocalInsert(i, string(text[i]))
	}

	AssertTrue(t, d.ToString() == text)
}

// calls DeleteAt and checks if peer documents are identical.
func TestLocalDelete(t *testing.T) {
	d := synceddoc.New()
	text := "hello everybody"
	for i := len(text) - 1; i >= 0; i-- {
		d.LocalInsert(0, string(text[i]))
	}

	AssertTrue(t, d.ToString() == text)

	d.LocalDelete(5)

	AssertTrue(t, d.ToString() == "helloeverybody")

	d.LocalDelete(13)
	d.LocalDelete(12)
	d.LocalDelete(9)
	d.LocalDelete(9)
	d.LocalDelete(9)
	d.LocalDelete(5)
	d.LocalDelete(5)
	d.LocalDelete(5)

	AssertTrue(t, d.ToString() == "hellor")

	d.LocalDelete(0)

	AssertTrue(t, d.ToString() == "ellor")

	d.LocalDelete(0)

	AssertTrue(t, d.ToString() == "llor")

	d.LocalDelete(0)
	d.LocalDelete(0)
	d.LocalDelete(0)
	d.LocalDelete(0)

	AssertTrue(t, d.ToString() == "")
}

func TestLocalInsertDelete(t *testing.T) {
	d := synceddoc.New()
	text := "helow everybody"
	for i := len(text) - 1; i >= 0; i-- {
		d.LocalInsert(0, string(text[i]))
	}

	AssertTrue(t, d.ToString() == text)

	d.LocalInsert(3, "l")

	AssertTrue(t, d.ToString() == "hellow everybody")

	d.LocalDelete(5)

	AssertTrue(t, d.ToString() == "hello everybody")

	d.LocalInsert(5, "!")

	AssertTrue(t, d.ToString() == "hello! everybody")

	d.LocalInsert(16, "?")

	AssertTrue(t, d.ToString() == "hello! everybody?")
}

func TestApplyRemoteOp(t *testing.T) {
	d := synceddoc.New()

	AssertTrue(t, d.ToString() == "")

	siteId := utils.UUID("site1")
	doc := crdt.NewBasicDocument(crdt.NewBasicPositionManager(siteId))

	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("h", 0),
		Val: "h",
	}, nil)

	AssertTrue(t, d.ToString() == "h")

	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("e", 1),
		Val: "e",
	}, nil)
	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("l", 2),
		Val: "l",
	}, nil)
	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("l", 3),
		Val: "l",
	}, nil)
	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("o", 4),
		Val: "o",
	}, nil)

	AssertTrue(t, d.ToString() == "hello")

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(0),
	}, nil)

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(2),
	}, nil)

	AssertTrue(t, d.ToString() == "elo")

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(1),
	}, nil)

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(1),
	}, nil)

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(0),
	}, nil)

	AssertTrue(t, d.ToString() == "")
}

// calls serialize on the document.
// returned value should later be deserialized into a valid document.
func TestSerialize(t *testing.T) {
	d := synceddoc.New()
	text := "hello everybody"
	for i := len(text) - 1; i >= 0; i-- {
		d.LocalInsert(0, string(text[i]))
	}
	docId := d.GetID()

	serialized, err := d.Serialize()
	AssertTrue(t, err == nil)

	nd, err := synceddoc.Load(serialized)
	AssertTrue(t, err == nil)

	AssertTrue(t, nd.GetID() == docId)
	AssertTrue(t, nd.ToString() == text)
}

func onChangeTest(changeName string, change interface {}, aux interface{}) {
	h := aux.(*[]rune)

	switch change.(type) {
	case synceddoc.MessageInsert:
		ch := change.(synceddoc.MessageInsert)
		newH := (*h)[:ch.Index]
		newH = append(newH, rune(ch.Value[0]))
		*h = append(newH, (*h)[ch.Index:]...)
	case synceddoc.MessageDelete:
		ch := change.(synceddoc.MessageDelete)
		newH := (*h)[:ch.Index]
		*h = append(newH, (*h)[ch.Index + 1:]...)
	}
}

// make changes on the document offline, and later call connect.
// peers should receive those changes after connect is called.
func TestConnectSignals(t *testing.T) {
	d := synceddoc.New()

	s := []rune {}

	siteId := utils.UUID("site1")
	doc := crdt.NewBasicDocument(crdt.NewBasicPositionManager(siteId))

	d.ConnectSignals(onChangeTest, nil, nil)

	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("h", 0),
		Val: "h",
	}, &s)

	AssertTrue(t, string(s) == "h")

	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("e", 1),
		Val: "e",
	}, &s)
	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("l", 2),
		Val: "l",
	}, &s)
	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("l", 3),
		Val: "l",
	}, &s)
	d.ApplyRemoteOp(siteId, crdt.OpInsert {
		Pos: doc.InsertAtIndex("o", 4),
		Val: "o",
	}, &s)

	AssertTrue(t, string(s) == "hello")

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(0),
	}, &s)

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(2),
	}, &s)

	AssertTrue(t, string(s) == "elo")

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(1),
	}, &s)

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(1),
	}, &s)

	d.ApplyRemoteOp(siteId, crdt.OpDelete {
		Pos: doc.DeleteAtIndex(0),
	}, &s)

	AssertTrue(t, string(s) == "")
}


// calls SetCursor and checks if peer documents are identical.
func TestSetCursor(t *testing.T) {

}