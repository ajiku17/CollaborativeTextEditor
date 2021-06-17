package crdt

import (
	"log"
	"testing"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

func DocumentInsertAtTop(b *testing.B, document Document, manager PositionManager) {
	text := "Hello again!"
	for n := 0; n < b.N; n++ {	
		manager.PositionManagerInit()
		document.DocumentInit(manager)
		InsertAtTop(document, utils.Reverse(text))
		if document.ToString() != text {
			b.Helper()
			b.Errorf("assertion failed")
		}
	}
}

func DocumentInsertAtBottom(b *testing.B, document Document, manager PositionManager) {
	text := "Hi everyone!"

	for n := 0; n < b.N; n++ {
		manager.PositionManagerInit()
		document.DocumentInit(manager)
		InsertAtBottom(document, text)
		if document.ToString() != text {
			b.Helper()
			b.Errorf("assertion failed")
		}
	}
}

func BenchmarkDocument(b *testing.B) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	implementations :=  []struct {
		newDocumentInstance func () Document
		newPositionManagerInstance func () PositionManager
		name string
	} {
		{ 
			func() Document {
				return new(BasicDocument)
			}, 
			func() PositionManager {
				return new(BasicPositionManager)
			}, 
			"BasicDocument"},
	}

	for _, impl := range implementations {
		b.Run(impl.name, func (b *testing.B) {
			b.Run("BenchmarkInsertAtBottom", func (b* testing.B) {
				DocumentInsertAtBottom(b, impl.newDocumentInstance(), impl.newPositionManagerInstance())
			})
			b.Run("BenchmarkDocumentInsertAtTop", func (b* testing.B) {
				DocumentInsertAtTop(b, impl.newDocumentInstance(), impl.newPositionManagerInstance())
			})
		})
	}
}