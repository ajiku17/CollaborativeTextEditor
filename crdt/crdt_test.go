package crdt

import (
	"fmt"
	"reflect"
	"testing"
)

func LongsEqual(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Helper()
		t.Errorf("expected %v actual %v!", expected, actual)
	}
}

func AssertTrue(t *testing.T, condition bool) {
	if !condition {
		t.Helper()
		t.Errorf("assertion failed")
	}
}

func TestNumberToInt(t *testing.T) {
	NumberSetBase(10)

	AssertTrue(t, NumberToInt(Number{1, 2, 3}) == 123)
}

func TestIsLessThan(t *testing.T) {
	NumberSetBase(64)

	AssertTrue(t, IsLessThan(Number{3, 2}, Number{3, 8}))
	AssertTrue(t, IsLessThan(Number{3, 2}, Number{1, 5, 8}))
	AssertTrue(t, !IsLessThan(Number{1, 5, 8}, Number{3, 2}))
	AssertTrue(t, IsLessThan(Number{3, 8}, Number{5, 6}))
}

func TestNumberAdd(t *testing.T) {
	NumberSetBase(10)

	// #1
	num1 := Number{3, 8}
	num2 := Number{9, 8}

	sum := NumberAdd(num1, num2)
	AssertTrue(t, NumberToInt(sum) == NumberToInt(Number{1, 3, 6}))

	// #2
	NumberSetBase(64)
	num1 = Number{3, 8} // 200
	num2 = Number{9, 8} // 584

	sum = NumberAdd(num1, num2) // 784 = 12 * 64 + 16
	AssertTrue(t, NumberToInt(sum) == NumberToInt(Number{12, 16}))

	// #3
	NumberSetBase(10)
	num1 = Number{3, 2}
	num2 = Number{1, 5, 8}

	sum = NumberAdd(num2, num1)
	AssertTrue(t, NumberToInt(sum) == NumberToInt(Number{1, 9, 0}))

	// #4
	NumberSetBase(10)
	num1 = Number{1, 5, 8}
	num2 = Number{3, 2}

	sum = NumberAdd(num2, num1)
	AssertTrue(t, NumberToInt(sum) == NumberToInt(Number{1, 9, 0}))
}

func TestNumberSubtract(t *testing.T) {
	NumberSetBase(10)

	// #1
	num1 := Number{3, 8}
	num2 := Number{9, 8}

	sum := NumberSubtract(num2, num1)
	AssertTrue(t, NumberToInt(sum) == NumberToInt(Number{6, 0}))

	// #2
	NumberSetBase(64)
	num1 = Number{3, 8} // 200
	num2 = Number{9, 8} // 584

	sum = NumberSubtract(num2, num1) // 384 = 6 * 64
	AssertTrue(t, NumberToInt(sum) == NumberToInt(Number{6, 0}))

	// #3
	NumberSetBase(64)
	num1 = Number{3, 8} // 200
	num2 = Number{5, 6} // 326

	sum = NumberSubtract(num2, num1) // 126 = 1 * 64 + 62
	AssertTrue(t, NumberToInt(sum) == NumberToInt(Number{1, 62}))

	// #4
	NumberSetBase(10)
	num1 = Number{5, 1}
	num2 = Number{3, 2}

	sum = NumberSubtract(num1, num2)
	AssertTrue(t, NumberToInt(sum) == NumberToInt(Number{1, 9}))

	// #5
	NumberSetBase(10)
	num1 = Number{3, 2}
	num2 = Number{1, 5, 8}

	sum = NumberSubtract(num2, num1)
	AssertTrue(t, NumberToInt(sum) == NumberToInt(Number{1, 2, 6}))
}

func TestPositionPrefix(t *testing.T) {
	var position Position
	zero := Identifier{0, 0}

	for i := 0; i < 10; i++ {
		position = append(position, Identifier{pos: i, site: i})
	}

	pref := prefix(position, 5)
	AssertTrue(t, len(pref) == 5)

	for i := 0; i < 5; i++ {
		AssertTrue(t, pref[i] == position[i])
	}

	pref = prefix(position, 12)
	AssertTrue(t, len(pref) == 12)

	for i := 0; i < 10; i++ {
		AssertTrue(t, pref[i] == position[i])
	}

	AssertTrue(t, pref[10] == zero)
	AssertTrue(t, pref[11] == zero)

	pref = prefix(position, 0)
	AssertTrue(t, 0 == len(pref))
}

func TestPositionToNumber(t *testing.T) {
	var position Position

	for i := 0; i < 6; i++ {
		position = append(position, Identifier{i, i})
	}

	number := PositionToNumber(position)
	expected := Number{0, 1, 2, 3, 4, 5}
	AssertTrue(t, reflect.DeepEqual(number, expected))
}

func TestPositionSubtract(t *testing.T) {
	var position1, position2 Position

	// #1
	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, NumberToInt(PositionSubtract(position2, position1)) == NumberToInt(Number{6, 0}))

	// #2
	NumberSetBase(64)
	position1 = Position{}
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 7})

	position2 = Position{}
	position2 = append(position2, Identifier{5, 5})
	position2 = append(position2, Identifier{6, 1})

	AssertTrue(t, NumberToInt(PositionSubtract(position2, position1)) == NumberToInt(Number{1, 62}))
}

func TestPositionAdd(t *testing.T) {
	var position1, position2 Position

	// #2
	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{9, 3})
	position1 = append(position1, Identifier{2, 2})

	position2 = append(position2, Identifier{3, 1})
	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{1, 5})

	AssertTrue(t, NumberToInt(PositionAdd(position1, position2)) == 783)

	// #2
	NumberSetBase(64)

	position1 = Position{}
	position2 = Position{}

	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, NumberToInt(PositionAdd(position1, position2)) == NumberToInt(Number{12, 16}))
}

func TestPositionAllocation(t *testing.T) {
	var prevPos, afterPos Position

	// #2
	NumberSetBase(10)
	prevPos = append(prevPos, Identifier{3, 2})
	prevPos = append(prevPos, Identifier{9, 3})
	prevPos = append(prevPos, Identifier{1, 2})

	afterPos = append(afterPos, Identifier{3, 1})
	afterPos = append(afterPos, Identifier{9, 5})
	afterPos = append(afterPos, Identifier{2, 5})

	newPosition := AllocPosition(prevPos, afterPos, 9)

	AssertTrue(t, len(newPosition) == 4)
	AssertTrue(t, IsLessThan(PositionToNumber(prevPos), PositionToNumber(newPosition)))
	AssertTrue(t, newPosition[len(newPosition)-1].site == 9)
	AssertTrue(t, newPosition[len(newPosition)-1].pos < 10)
}

func TestDocCreation(t *testing.T) {
	document := NewDocument()
	AssertTrue(t, document.GetLength() == 2)
	AssertTrue(t, (*document)[0].position[0].pos == 0)

}

func TestDocInsert(t *testing.T) {
	document := new(Document)

	// Insert at the beginning
	document = document.docInsert(0, Element{"data", nil})
	AssertTrue(t, document.GetLength() == 1)

	document = document.docInsert(1, Element{"data1", Position{Identifier{0, 5}}})
	AssertTrue(t, document.GetLength() == 2)
	AssertTrue(t, (*document)[1].position[0].site == 5)

	// Insert in the middle
	document = document.docInsert(1, Element{"data2", Position{Identifier{0, 7}}})
	AssertTrue(t, document.GetLength() == 3)
	AssertTrue(t, (*document)[1].position[0].site == 7)
	AssertTrue(t, (*document)[2].position[0].site == 5)

	// Insert at the end
	document = document.docInsert(3, Element{"end", Position{Identifier{0, 7}}})
	AssertTrue(t, document.GetLength() == 4)
	AssertTrue(t, (*document)[3].data == "end")
}

func TestDocument(t *testing.T) {
	document := NewDocument()
	document.InsertAt("H", 0, 1)
	// document.InsertAt("i", 0, 4)
	fmt.Printf("%v", document)
	// document.InsertAt("H", 0, 1)
	// document.InsertAt("H", 0, 1)
	// document.InsertAt("H", 0, 1)
	// document.InsertAt("H", 0, 1)
	// document.InsertAt("H", 0, 1)
}
