package test

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"reflect"
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
	d := synceddoc.New("1")

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

	d = synceddoc.New("1")
	text := "hello everybody"
	for i := len(text) - 1; i >= 0; i-- {
		d.LocalInsert(0, string(text[i]))
	}

	AssertTrue(t, d.ToString() == text)

	d = synceddoc.New("1")
	text = "hey everyone"
	for i := 0; i < len(text); i++ {
		d.LocalInsert(i, string(text[i]))
	}

	AssertTrue(t, d.ToString() == text)
}

// calls DeleteAt and checks if peer documents are identical.
func TestLocalDelete(t *testing.T) {
	d := synceddoc.New("1")
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
	d := synceddoc.New("1")
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
	d := synceddoc.New("1")

	AssertTrue(t, d.ToString() == "")

	siteId := utils.UUID("site1")
	doc := crdt.NewBasicDocument(crdt.NewBasicPositionManager(siteId))

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 1,
		Cmd:         crdt.OpInsert {
				Pos: doc.InsertAtIndex("h", 0),
				Val: "h",
			},
		}, nil)

	AssertTrue(t, d.ToString() == "h")

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 2,
		Cmd:         crdt.OpInsert {
			Pos: doc.InsertAtIndex("e", 1),
			Val: "e",
		},
	}, nil)
	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 3,
		Cmd:         crdt.OpInsert {
			Pos: doc.InsertAtIndex("l", 2),
			Val: "l",
		},
	},nil)
	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 4,
		Cmd:         crdt.OpInsert {
			Pos: doc.InsertAtIndex("l", 3),
			Val: "l",
		},
	}, nil)
	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 5,
		Cmd:         crdt.OpInsert {
			Pos: doc.InsertAtIndex("o", 4),
			Val: "o",
		},
	}, nil)

	AssertTrue(t, d.ToString() == "hello")

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 6,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(0),
		},
	}, nil)

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 7,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(2),
		},
	}, nil)

	AssertTrue(t, d.ToString() == "elo")

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 8,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(1),
		},
	}, nil)

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 9,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(1),
		},
	}, nil)

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 10,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(0),
		},
	}, nil)

	AssertTrue(t, d.ToString() == "")
}

// calls serialize on the document.
// returned value should later be deserialized into a valid document.
func TestSerialize(t *testing.T) {
	d := synceddoc.New("1")
	text := "hello everybody"
	for i := len(text) - 1; i >= 0; i-- {
		d.LocalInsert(0, string(text[i]))
	}
	docId := d.GetID()

	serialized, err := d.Serialize()
	AssertTrue(t, err == nil)

	nd, err := synceddoc.Load("1", serialized)
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
	d := synceddoc.New("1")

	s := []rune {}

	siteId := utils.UUID("site1")
	doc := crdt.NewBasicDocument(crdt.NewBasicPositionManager(siteId))

	d.ConnectSignals(onChangeTest, nil, nil)

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 1,
		Cmd:         crdt.OpInsert {
				Pos: doc.InsertAtIndex("h", 0),
				Val: "h",
			},
	}, &s)

	AssertTrue(t, string(s) == "h")

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 2,
		Cmd:         crdt.OpInsert {
				Pos: doc.InsertAtIndex("e", 1),
				Val: "e",
			},
	}, &s)
	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 3,
		Cmd:         crdt.OpInsert {
				Pos: doc.InsertAtIndex("l", 2),
				Val: "l",
			},
	}, &s)
	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 4,
		Cmd:         crdt.OpInsert {
				Pos: doc.InsertAtIndex("l", 3),
				Val: "l",
			},
	}, &s)
	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 5,
		Cmd:         crdt.OpInsert {
				Pos: doc.InsertAtIndex("o", 4),
				Val: "o",
			},
	}, &s)

	AssertTrue(t, string(s) == "hello")

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 6,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(0),
		},
	}, &s)

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 7,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(2),
		},
	}, &s)

	AssertTrue(t, string(s) == "elo")

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 8,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(1),
		},
	}, &s)

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 9,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(1),
		},
	}, &s)

	d.ApplyRemoteOp(synceddoc.Op{
		PeerId:      siteId,
		PeerOpIndex: 10,
		Cmd:         crdt.OpDelete {
			Pos: doc.DeleteAtIndex(0),
		},
	}, &s)

	AssertTrue(t, string(s) == "")
}


// calls SetCursor and checks if peer documents are identical.
func TestSetCursor(t *testing.T) {

}

func TestGetIntersecting(t *testing.T) {
	intervals := []synceddoc.Interval{{20, 35}, {50, 52}, {56, 70}, {75, 78}, {81, 83} ,{85, 90}}
	interval := synceddoc.Interval{45, 56}

	intersecting := synceddoc.GetIntersecting(interval, intervals)
	AssertTrue(t, reflect.DeepEqual(intersecting, []synceddoc.Interval{{50, 52}, {56, 70}}))

	interval = synceddoc.Interval{35, 45}

	intersecting = synceddoc.GetIntersecting(interval, intervals)
	AssertTrue(t, reflect.DeepEqual(intersecting, []synceddoc.Interval{{20, 35}}))

	interval = synceddoc.Interval{36, 45}

	intersecting = synceddoc.GetIntersecting(interval, intervals)
	AssertTrue(t, len(intersecting) == 0)

	interval = synceddoc.Interval{0, 100}

	intersecting = synceddoc.GetIntersecting(interval, intervals)
	AssertTrue(t, reflect.DeepEqual(intersecting, intervals))
}

func TestFindMissingIndices(t *testing.T) {
	intervals1 := []synceddoc.Interval{{20, 35}, {50, 52}}
	intervals2 := []synceddoc.Interval{{20, 35}, {50, 52}}

	missing := synceddoc.FindMissingIndices(intervals1, intervals2)
	fmt.Println(missing)
	AssertTrue(t, len(missing) == 0)

	intervals1 = []synceddoc.Interval{{30, 35}, {38, 50}, {58, 79}, {81, 85} ,{87, 90}}
	intervals2 = []synceddoc.Interval{{20, 35}, {50, 52}, {56, 70}, {75, 78}, {81, 83} ,{85, 90}}

	missing = synceddoc.FindMissingIndices(intervals1, intervals2)
	fmt.Println(missing)
	AssertTrue(t, reflect.DeepEqual(missing, []int{20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 51, 52, 56, 57, 86}))

	intervals1 = []synceddoc.Interval{}
	intervals2 = []synceddoc.Interval{{20, 21}, {50, 52}}

	missing = synceddoc.FindMissingIndices(intervals1, intervals2)
	fmt.Println(missing)
	AssertTrue(t, reflect.DeepEqual(missing, []int{20, 21, 50, 51, 52}))

	intervals1 = []synceddoc.Interval{ {1, 12},  {16, 20}}
	intervals2 = []synceddoc.Interval{{10, 21}}

	missing = synceddoc.FindMissingIndices(intervals1, intervals2)
	fmt.Println(missing)
	AssertTrue(t, reflect.DeepEqual(missing, []int{13, 14, 15, 21}))

	intervals1 = []synceddoc.Interval{ {1, 12},  {16, 20}}
	intervals2 = []synceddoc.Interval{{10, 21}, {50, 53}}

	missing = synceddoc.FindMissingIndices(intervals1, intervals2)
	fmt.Println(missing)
	AssertTrue(t, reflect.DeepEqual(missing, []int{13, 14, 15, 21, 50, 51, 52, 53}))

	intervals1 = []synceddoc.Interval{ {1, 12},  {16, 20}, {22, 50}}
	intervals2 = []synceddoc.Interval{{10, 24}, {50, 53}}

	missing = synceddoc.FindMissingIndices(intervals1, intervals2)
	fmt.Println(missing)
	AssertTrue(t, reflect.DeepEqual(missing, []int{13, 14, 15, 21, 51, 52, 53}))

	intervals1 = []synceddoc.Interval{}
	intervals2 = []synceddoc.Interval{}

	missing = synceddoc.FindMissingIndices(intervals1, intervals2)
	fmt.Println(missing)
	AssertTrue(t, len(missing) == 0)

	intervals1 = []synceddoc.Interval{{20, 21}, {50, 52}}
	intervals2 = []synceddoc.Interval{}

	missing = synceddoc.FindMissingIndices(intervals1, intervals2)
	fmt.Println(missing)
	AssertTrue(t, len(missing) == 0)
}



func TestCreatePatch(t *testing.T) {
	d := synceddoc.New("1")

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

	d1, err := synceddoc.Open("2", string(d.GetID()))
	AssertTrue(t, err == nil)

	d1State := d1.GetCurrentState()

	patch := d.CreatePatch(d1State)

	d1.ApplyPatch(patch)

	AssertTrue(t, d.GetID() == d1.GetID())
	fmt.Println("d string:", d.ToString())
	fmt.Println("d1 string:", d1.ToString())
	AssertTrue(t, d.ToString() == d1.ToString())

	d1.ApplyPatch(patch) // apply patch once more

	AssertTrue(t, d.GetID() == d1.GetID())
	fmt.Println("d string:", d.ToString())
	fmt.Println("d1 string:", d1.ToString())
	AssertTrue(t, d.ToString() == d1.ToString())

	d2, err := synceddoc.Open("2", string(d.GetID()))
	AssertTrue(t, err == nil)

	d2State := d2.GetCurrentState()

	patch = d1.CreatePatch(d2State)

	d2.ApplyPatch(patch)

	AssertTrue(t, d1.GetID() == d2.GetID())
	fmt.Println("d1 string:", d1.ToString())
	fmt.Println("d2 string:", d2.ToString())
	AssertTrue(t, d1.ToString() == d2.ToString())
}

func TestPatchApply(t *testing.T) {
	d1 := synceddoc.New("1")

	d1.LocalInsert(0, "h")
	d1.LocalInsert(1, "e")
	d1.LocalInsert(2, "l")
	d1.LocalInsert(3, "o")
	d1.LocalInsert(4, "w")
	d1.LocalInsert(5, "o")
	d1.LocalInsert(6, "r")
	d1.LocalInsert(7, "l")
	d1.LocalInsert(8, "d")
	d1.LocalInsert(3, "l")
	d1.LocalInsert(5, " ")

	d2, err := synceddoc.Open("2", string(d1.GetID()))
	AssertTrue(t, err == nil)

	d2.LocalInsert(0, "ბ")
	d2.LocalInsert(1, "a")
	d2.LocalInsert(2, "d")
	d2.LocalInsert(3, "o")
	d2.LocalInsert(4, "R")
	d2.LocalInsert(5, "o")
	d2.LocalInsert(6, "q")
	d2.LocalInsert(7, "l")
	d2.LocalInsert(8, "d")
	d2.LocalInsert(3, "p")
	d2.LocalInsert(5, " ")

	fmt.Println(d1.ToString())
	fmt.Println(d2.ToString())

	d1State := d1.GetCurrentState()
	d2State := d2.GetCurrentState()

	d1Patch := d2.CreatePatch(d1State)
	d2Patch := d1.CreatePatch(d2State)

	d1.ApplyPatch(d1Patch)
	d2.ApplyPatch(d2Patch)

	AssertTrue(t, d1.GetID() == d2.GetID())
	fmt.Println(d1.ToString())
	fmt.Println(d2.ToString())
	AssertTrue(t, d1.ToString() == d2.ToString())
}

func TestOverlappingPatchApply(t *testing.T) {
	d1 := synceddoc.New("1")

	d1.LocalInsert(0, "h")
	d1.LocalInsert(1, "e")
	d1.LocalInsert(2, "l")
	d1.LocalInsert(3, "o")
	d1.LocalInsert(4, "w")
	d1.LocalInsert(5, "o")
	d1.LocalInsert(6, "r")
	d1.LocalInsert(7, "l")
	d1.LocalInsert(8, "d")
	d1.LocalInsert(3, "l")
	d1.LocalInsert(5, " ")

	d2, err := synceddoc.Open("2", string(d1.GetID()))
	AssertTrue(t, err == nil)

	d2.LocalInsert(0, "ბ")
	d2.LocalInsert(1, "a")
	d2.LocalInsert(2, "d")
	d2.LocalInsert(3, "o")
	d2.LocalInsert(4, "R")
	d2.LocalInsert(5, "o")
	d2.LocalInsert(6, "q")
	d2.LocalInsert(7, "l")
	d2.LocalInsert(8, "d")
	d2.LocalInsert(3, "p")
	d2.LocalInsert(5, " ")

	fmt.Println(d1.ToString())
	fmt.Println(d2.ToString())

	d1State := d1.GetCurrentState()
	d2State := d2.GetCurrentState()

	d1Patch := d2.CreatePatch(d1State)

	d1.ApplyPatch(d1Patch)

	// add modifications to the document
	d1.LocalInsert(8, "a")
	d1.LocalInsert(3, "b")
	d1.LocalInsert(5, "c")

	d2Patch := d1.CreatePatch(d2State)

	d2.ApplyPatch(d2Patch)

	AssertTrue(t, d1.GetID() == d2.GetID())
	fmt.Println(d1.ToString())
	fmt.Println(d2.ToString())
	AssertTrue(t, d1.ToString() == d2.ToString())
}

