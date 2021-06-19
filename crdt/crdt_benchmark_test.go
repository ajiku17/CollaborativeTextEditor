package crdt

import (
	"log"
	"testing"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

func DocumentInsertAtTop(b *testing.B, newDocumentInstance func () Document) {
	text := "Hello again!"
	for n := 0; n < b.N; n++ {	
		document := newDocumentInstance()
		InsertAtTop(document, utils.Reverse(text))
		if document.ToString() != text {
			b.Helper()
			b.Errorf("assertion failed")
		}
	}
}

func DocumentInsertAtBottom(b *testing.B, newDocumentInstance func () Document) {
	text := "Hi everyone!"

	for n := 0; n < b.N; n++ {
		document := newDocumentInstance()
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
		name string
	} {
		{ 
			func() Document {
				return NewBasicDocument(NewBasicPositionManager())
			},
			"BasicDocument",
		},
	}

	for _, impl := range implementations {
		b.Run(impl.name, func (b *testing.B) {
			b.Run("BenchmarkInsertAtBottom", func (b* testing.B) {
				DocumentInsertAtBottom(b, impl.newDocumentInstance)
			})
			b.Run("BenchmarkDocumentInsertAtTop", func (b* testing.B) {
				DocumentInsertAtTop(b, impl.newDocumentInstance)
			})
		})
	}
}