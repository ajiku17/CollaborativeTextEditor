package crdt

import (
	"log"
	"testing"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

func InsertAtTop(doc Document, text string) {
	for _, character := range text {
		doc.InsertAtIndex(string(character), 0, utils.RandBetween(1, 5))
	}
}

func InsertAtBottom(doc Document, text string) {
	for index, character := range text {
		doc.InsertAtIndex(string(character), index, utils.RandBetween(1, 5))
	}
}

func DocumentInsertAtIndex(t *testing.T, newDocumentInstance func () Document) {
	
	// #1
	document := newDocumentInstance()

	text := "Hi everyone!"
	InsertAtBottom(document, text)
	AssertTrue(t, document.ToString() == text)

	// #2
	document = newDocumentInstance()
	
	text = "Hello again!"
	InsertAtTop(document, utils.Reverse(text))
	AssertTrue(t, document.ToString() == text)

	// #3
	document = newDocumentInstance()

	text = "Hello!"
	document.InsertAtIndex("e", 0, 1)
	document.InsertAtIndex("l", 1, 4)
	document.InsertAtIndex("o", 2, 3)
	document.InsertAtIndex("l", 1, 1)
	document.InsertAtIndex("!", 4, 2)
	document.InsertAtIndex("H", 0, 4)
	AssertTrue(t, document.ToString() == text)
}

func DocumentDeleteAtIndex(t *testing.T, newDocumentInstance func () Document) {

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

	document.InsertAtIndex("H", 0, 4)
	document.InsertAtIndex("e", 1, 1)
	document.InsertAtIndex("l", 2, 4)
	document.InsertAtIndex("l", 3, 1)
	document.InsertAtIndex("o", 4, 1)
	AssertTrue(t, document.Length() == 15)
	AssertTrue(t, document.ToString() == "Hello everyone!")

	// #2
	document = newDocumentInstance()

	document.InsertAtIndex("H", 0, 1)
	document.InsertAtIndex("i", 1, 4)
	document.InsertAtIndex(" ", 2, 1)
	document.InsertAtIndex("e", 3, 4)
	document.InsertAtIndex("v", 4, 1)
	document.InsertAtIndex("e", 5, 4)
	document.InsertAtIndex("r", 6, 1)
	document.InsertAtIndex("y", 7, 4)
	document.InsertAtIndex("o", 8, 1)
	document.InsertAtIndex("n", 9, 4)
	document.InsertAtIndex("e", 10, 1)
	document.InsertAtIndex("!", 11, 1)
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

	document.InsertAtIndex("f", 3, 4)
	document.InsertAtIndex("o", 4, 1)
	document.InsertAtIndex("l", 5, 4)
	document.InsertAtIndex("k", 6, 1)
	document.InsertAtIndex("s", 7, 1)
	AssertTrue(t, document.Length() == 9)
	AssertTrue(t, document.ToString() == "Hi folks!")
}

func DocInsertAtPosition(t *testing.T, newDocumentInstance func () Document, newManagerInstance func () PositionManager) {
	// #1
	document := newDocumentInstance()
	manager := newManagerInstance()

	text := "Hi everyone!"
	prev := manager.GetMinPosition()
	next := manager.GetMaxPosition()
	positions := []struct {
		pos Position
		val string
	}{}

	for _, c := range text {
		pos := manager.AllocPositionBetween(prev, next, utils.RandBetween(1, 5))
		positions = append(positions, struct {
			pos Position
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

func DocDeleteAtPos(t *testing.T, newDocumentInstance func () Document, newManagerInstance func () PositionManager) {
	document := newDocumentInstance()
	manager := newManagerInstance()

	// #1
	text := "Hi everyone!"
	prev := manager.GetMinPosition()
	next := manager.GetMaxPosition()
	positions := []struct {
		pos Position
		val string
	}{}

	for _, c := range text {
		pos := manager.AllocPositionBetween(prev, next, utils.RandBetween(1, 5))
		positions = append(positions, struct {
			pos Position
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

func TestDocument(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	implementations :=  []struct {
		newDocumentInstance func () Document
		newPositionManagerInstance func () PositionManager
		name string
	} {
		{ 
			func() Document {
				return NewBasicDocument(NewBasicPositionManager())
			},
			func() PositionManager {
				return NewBasicPositionManager()
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
		})
	}
}