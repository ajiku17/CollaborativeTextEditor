package test

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"log"
	"testing"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

func InsertAtTop(doc crdt.Document, text string) {
	for _, character := range text {
		doc.InsertAtIndex(string(character), 0)
	}
}

func InsertAtBottom(doc crdt.Document, text string) {
	for index, character := range text {
		doc.InsertAtIndex(string(character), index)
	}
}

func DocumentInsertAtIndex(t *testing.T, newDocumentInstance func () crdt.Document) {
	
	// #1
	document := newDocumentInstance()

	text := "Hi everyone!"
	InsertAtBottom(document, text)
	AssertTrue(t, document.Length() == 12)
	AssertTrue(t, document.ToString() == text)

	// #2
	document = newDocumentInstance()
	
	text = "Hello again!"
	InsertAtTop(document, utils.Reverse(text))
	AssertTrue(t, document.Length() == 12)
	AssertTrue(t, document.ToString() == text)

	// #3
	document = newDocumentInstance()

	document.InsertAtIndex("e", 0)
	document.InsertAtIndex("l", 1)
	document.InsertAtIndex("o", 2)
	document.InsertAtIndex("l", 1)
	document.InsertAtIndex("!", 4)
	document.InsertAtIndex("H", 0)
	AssertTrue(t, document.Length() == 6)
	AssertTrue(t, document.ToString() == "Hello!")
}

func DocumentDeleteAtIndex(t *testing.T, newDocumentInstance func () crdt.Document) {

	// #1
	document := newDocumentInstance()

	text := "Hi everyone!"
	InsertAtBottom(document, text)
	AssertTrue(t, document.Length() == 12)
	AssertTrue(t, document.ToString() == text)

	document.DeleteAtIndex(0)
	document.DeleteAtIndex(0)
	AssertTrue(t, document.Length() == 10)
	AssertTrue(t, document.ToString() == " everyone!")

	document.InsertAtIndex("H", 0)
	document.InsertAtIndex("e", 1)
	document.InsertAtIndex("l", 2)
	document.InsertAtIndex("l", 3)
	document.InsertAtIndex("o", 4)
	AssertTrue(t, document.Length() == 15)
	AssertTrue(t, document.ToString() == "Hello everyone!")

	// #2
	document = newDocumentInstance()

	document.InsertAtIndex("H", 0)
	document.InsertAtIndex("i", 1)
	document.InsertAtIndex(" ", 2)
	document.InsertAtIndex("e", 3)
	document.InsertAtIndex("v", 4)
	document.InsertAtIndex("e", 5)
	document.InsertAtIndex("r", 6)
	document.InsertAtIndex("y", 7)
	document.InsertAtIndex("o", 8)
	document.InsertAtIndex("n", 9)
	document.InsertAtIndex("e", 10)
	document.InsertAtIndex("!", 11)
	AssertTrue(t, document.Length() == 12)
	AssertTrue(t, document.ToString() == "Hi everyone!")

	document.DeleteAtIndex(3)
	document.DeleteAtIndex(3)
	document.DeleteAtIndex(3)
	document.DeleteAtIndex(3)
	document.DeleteAtIndex(3)
	document.DeleteAtIndex(3)
	document.DeleteAtIndex(3)
	document.DeleteAtIndex(3)
	AssertTrue(t, document.Length() == 4)
	AssertTrue(t, document.ToString() == "Hi !")

	document.InsertAtIndex("f", 3)
	document.InsertAtIndex("o", 4)
	document.InsertAtIndex("l", 5)
	document.InsertAtIndex("k", 6)
	document.InsertAtIndex("s", 7)
	AssertTrue(t, document.Length() == 9)
	AssertTrue(t, document.ToString() == "Hi folks!")
}

func DocInsertAtPosition(t *testing.T, newDocumentInstance func () crdt.Document, newManagerInstance func () crdt.PositionManager) {
	// #1
	document := newDocumentInstance()
	manager := newManagerInstance()

	text := "Hi everyone!"
	prev := manager.GetMinPosition()
	next := manager.GetMaxPosition()
	positions := []struct {
		pos crdt.Position
		val string
	}{}

	for _, c := range text {
		pos := manager.AllocPositionBetween(prev, next)
		positions = append(positions, struct {
			pos crdt.Position
			val string
		}{pos, string(c)})
		prev = pos
	}


	for _, e := range positions {
		document.InsertAtPosition(e.pos, e.val)
	}

	AssertTrue(t, document.Length() == 12)
	AssertTrue(t, document.ToString() == "Hi everyone!")

	// #2
	document = newDocumentInstance()

	shuffled := positions[:]
	for i := range shuffled {
		j := utils.RandBetween(0, len(shuffled) - 1)
		tmp := shuffled[i]
		shuffled[i] = shuffled[j]
		shuffled[j] = tmp
	}
	
	for _, e := range shuffled {
		document.InsertAtPosition(e.pos, e.val)
	}

	AssertTrue(t, document.Length() == 12)
	AssertTrue(t, document.ToString() == "Hi everyone!")
}

func DocDeleteAtPos(t *testing.T, newDocumentInstance func () crdt.Document, newManagerInstance func () crdt.PositionManager) {
	document := newDocumentInstance()
	manager := newManagerInstance()

	// #1
	text := "Hi everyone!"
	prev := manager.GetMinPosition()
	next := manager.GetMaxPosition()
	positions := []struct {
		pos crdt.Position
		val string
	}{}

	for _, c := range text {
		pos := manager.AllocPositionBetween(prev, next)
		positions = append(positions, struct {
			pos crdt.Position
			val string
		}{pos, string(c)})
		prev = pos
	}


	for _, e := range positions {
		document.InsertAtPosition(e.pos, e.val)
	}

	AssertTrue(t, document.Length() == 12)
	AssertTrue(t, document.ToString() == "Hi everyone!")

	shuffled := positions[2:len(positions) - 1]
	for i := range shuffled {
		j := utils.RandBetween(0, len(shuffled) - 1)
		tmp := shuffled[i]
		shuffled[i] = shuffled[j]
		shuffled[j] = tmp
	}
	
	for i := 0; i < len(shuffled); i++ {
		document.DeleteAtPosition(shuffled[i].pos)
	}
	AssertTrue(t, document.Length() == 3)
	AssertTrue(t, document.ToString() == "Hi!")
}

func DocSerializeDeserialize(t *testing.T, newDocumentInstance func () crdt.Document) {
	document := newDocumentInstance()

	document.InsertAtIndex("H", 0)
	document.InsertAtIndex("e", 1)
	document.InsertAtIndex("l", 2)
	document.InsertAtIndex("l", 3)
	document.InsertAtIndex("o", 4)
	document.InsertAtIndex(" ", 5)
	document.InsertAtIndex("W", 6)
	document.InsertAtIndex("o", 7)
	document.InsertAtIndex("r", 8)
	document.InsertAtIndex("l", 9)
	document.InsertAtIndex("d", 10)

	serialized, err := document.Serialize()
	AssertTrue(t, err == nil)
	AssertTrue(t, document.ToString() == "Hello World")

	deserializedDoc := newDocumentInstance()

	err = deserializedDoc.Deserialize(serialized)
	AssertTrue(t, err == nil)

	AssertTrue(t, deserializedDoc.Length() == document.Length())
	AssertTrue(t, deserializedDoc.ToString() == document.ToString())

	deserializedDoc.InsertAtIndex("!", 11)
	AssertTrue(t, deserializedDoc.ToString() == "Hello World!")
}

func TestDocument(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	implementations :=  []struct {
		newDocumentInstance func () crdt.Document
		newPositionManagerInstance func () crdt.PositionManager
		name string
	} {
		{ 
			func() crdt.Document {
				return crdt.NewBasicDocument(crdt.NewBasicPositionManager("1"))
			},
			func() crdt.PositionManager {
				return crdt.NewBasicPositionManager("1")
			},
			"BasicDocument",
		},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func (t *testing.T) {
			t.Run("TestDocumentInsertAtIndex", func (t* testing.T) {
				DocumentInsertAtIndex(t, impl.newDocumentInstance)
			})
			t.Run("TestDocumentDeleteAtIndex", func (t* testing.T) {
				DocumentDeleteAtIndex(t, impl.newDocumentInstance)
			})
			t.Run("TestDocInsertAtPosition", func (t* testing.T) {
				DocInsertAtPosition(t, impl.newDocumentInstance, impl.newPositionManagerInstance)
			})
			t.Run("TestDocDeleteAtPos", func (t* testing.T) {
				DocDeleteAtPos(t, impl.newDocumentInstance, impl.newPositionManagerInstance)
			})
			t.Run("TestDocSerializeDeserialize", func (t* testing.T) {
				DocSerializeDeserialize(t, impl.newDocumentInstance)
			})
		})
	}
}