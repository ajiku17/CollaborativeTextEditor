package crdt

import (
	"testing"

	"github.com/utils"
)

func BenchmarkInsertAtTop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		text := "Hello again!"
		document := InsertAtTop(utils.Reverse(text))
		if document.ToString() != text {
			b.Helper()
			b.Errorf("assertion failed")
		}
	}
}

func BenchmarkInsertAtBottom(b *testing.B) {
	for n := 0; n < b.N; n++ {
		text := "Hi everyone!"
		document := InsertAtBottom(text)
		if document.ToString() != text {
			b.Helper()
			b.Errorf("assertion failed")
		}
	}
}
