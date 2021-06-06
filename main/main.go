package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/crdt"
)

func main() {
	doc := crdt.NewDocument()

	doc.InsertAt("H", 0, 1)
	doc.InsertAt("e", 1, 1)
	doc.InsertAt("l", 2, 1)
	doc.InsertAt("l", 3, 1)
	doc.InsertAt("o", 4, 1)
	doc.InsertAt(" ", 5, 1)
	doc.InsertAt("W", 6, 1)
	doc.InsertAt("o", 7, 1)
	doc.InsertAt("r", 8, 1)
	doc.InsertAt("l", 9, 1)
	doc.InsertAt("d", 10, 1)

	fmt.Println(doc.ToString())

	result := testing.Benchmark(func(parentB *testing.B) {
		parentB.Run("example", func(b *testing.B) {
			fmt.Printf("ok")
		})
	})

	fmt.Printf("%v", result)

	empty_benchmark := testing.BenchmarkResult{}
	if reflect.DeepEqual(result, empty_benchmark) {
		fmt.Printf("aaa")
	}
}
