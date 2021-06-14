package crdt

import (
	"math"
	"reflect"
	"testing"

	"github.com/utils"
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

	for i := 0; i < 10; i++ {
		position = append(position, Identifier{pos: i, site: i})
	}

	pref := prefix(position, 5)
	AssertTrue(t, len(pref) == 5)

	for i := 0; i < 5; i++ {
		AssertTrue(t, pref[i] == position[i].pos)
	}

	pref = prefix(position, 12)
	AssertTrue(t, len(pref) == 12)

	for i := 0; i < 10; i++ {
		AssertTrue(t, pref[i] == position[i].pos)
	}

	AssertTrue(t, pref[10] == 0)
	AssertTrue(t, pref[11] == 0)

	pref = prefix(position, 0)
	AssertTrue(t, len(pref) == 0)
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

func TestPositionIsGreaterOrEqual(t *testing.T) {
	var position1, position2 Position

	// #1
	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, !PositionIsGreaterOrEqual(position1, position2))
	AssertTrue(t, PositionIsGreaterOrEqual(position2, position1))

	// #2
	NumberSetBase(64)
	position1 = Position{}
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 7})

	position2 = Position{}
	position2 = append(position2, Identifier{3, 2})
	position2 = append(position2, Identifier{8, 7})

	AssertTrue(t, !PositionIsLessThan(position1, position2))
	AssertTrue(t, !PositionIsLessThan(position2, position1))

	position2 = append(position2, Identifier{1, 2})
	AssertTrue(t, PositionIsGreaterOrEqual(position2, position1))
	AssertTrue(t, !PositionIsGreaterOrEqual(position1, position2))
}

func TestPositionIsLessThan(t *testing.T) {
	var position1, position2 Position

	// #1
	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, PositionIsLessThan(position1, position2))
	AssertTrue(t, !PositionIsLessThan(position2, position1))

	// #2
	NumberSetBase(64)
	position1 = Position{}
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 7})

	position2 = Position{}
	position2 = append(position2, Identifier{3, 2})
	position2 = append(position2, Identifier{8, 7})

	AssertTrue(t, !PositionIsLessThan(position1, position2))
	AssertTrue(t, !PositionIsLessThan(position2, position1))

	position2 = append(position2, Identifier{1, 2})
	AssertTrue(t, PositionIsLessThan(position1, position2))
	AssertTrue(t, !PositionIsLessThan(position2, position1))
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

func TestConstructPosition(t *testing.T) {
	var r Number
	var prevPos, afterPos Position

	prevPos = append(prevPos, Identifier{3, 2})
	prevPos = append(prevPos, Identifier{9, 3})
	prevPos = append(prevPos, Identifier{1, 2})

	afterPos = append(afterPos, Identifier{3, 1})
	afterPos = append(afterPos, Identifier{9, 5})
	afterPos = append(afterPos, Identifier{2, 5})

	r = Number{3, 9, 1, 4}

	newPos := constructPosition(r, prevPos, afterPos, 8)
	AssertTrue(t, len(newPos) == 4)
	for i := 0; i < len(prevPos); i++ {
		AssertTrue(t, prevPos[i].pos == newPos[i].pos && prevPos[i].site == newPos[i].site)
	}
	AssertTrue(t, newPos[len(newPos) - 1].pos == 4 && newPos[len(newPos) - 1].site == 8)
}

func TestPositionAllocation(t *testing.T) {
	var prevPos, afterPos Position

	// #1
	NumberSetBase(10)
	prevPos = append(prevPos, Identifier{7, 2})

	afterPos = append(afterPos, Identifier{10, 1})

	newPosition := AllocPosition(prevPos, afterPos, 4)

	AssertTrue(t, len(newPosition) == 1)
	AssertTrue(t, IsLessThan(PositionToNumber(prevPos), PositionToNumber(newPosition)))
	AssertTrue(t, newPosition[len(newPosition) - 1].site == 4)
	AssertTrue(t, newPosition[len(newPosition) - 1].pos > 7)
	AssertTrue(t, newPosition[len(newPosition) - 1].pos < 10)

	// #2
	NumberSetBase(10)

	prevPos = Position{}
	afterPos = Position{}

	prevPos = append(prevPos, Identifier{3, 2})
	prevPos = append(prevPos, Identifier{9, 3})
	prevPos = append(prevPos, Identifier{1, 2})

	afterPos = append(afterPos, Identifier{3, 1})
	afterPos = append(afterPos, Identifier{9, 5})
	afterPos = append(afterPos, Identifier{2, 5})

	newPosition = AllocPosition(prevPos, afterPos, 9)

	AssertTrue(t, len(newPosition) == 4)
	AssertTrue(t, IsLessThan(PositionToNumber(prevPos), PositionToNumber(newPosition)))
	AssertTrue(t, newPosition[len(newPosition) - 1].site == 9)
	AssertTrue(t, newPosition[len(newPosition) - 1].pos < 10)
}

func TestDocInsert(t *testing.T) {
	document := new(Document)

	// Insert at the beginning
	document.docInsert(0, Element{"data", nil})
	AssertTrue(t, len(*document) == 1)

	document.docInsert(1, Element{"data1", Position{Identifier{0, 5}}})
	AssertTrue(t, len(*document) == 2)
	AssertTrue(t, (*document)[1].position[0].site == 5)

	// Insert in the middle
	document.docInsert(1, Element{"data2", Position{Identifier{0, 7}}})
	AssertTrue(t, len(*document) == 3)
	AssertTrue(t, (*document)[1].position[0].site == 7)
	AssertTrue(t, (*document)[2].position[0].site == 5)

	// Insert at the end
	document.docInsert(3, Element{"end", Position{Identifier{0, 7}}})
	AssertTrue(t, len(*document) == 4)
	AssertTrue(t, (*document)[3].data == "end")
}

func TestDocDelete(t *testing.T) {
	document := new(Document)

	// Insert at the beginning
	document.docInsert(0, Element{"begin", nil})
	AssertTrue(t, len(*document) == 1)

	document.docInsert(1, Element{"data1", Position{Identifier{0, 5}}})
	AssertTrue(t, len(*document) == 2)
	AssertTrue(t, (*document)[1].position[0].site == 5)

	// Insert in the middle
	document.docInsert(1, Element{"data2", Position{Identifier{0, 6}}})
	AssertTrue(t, len(*document) == 3)
	AssertTrue(t, (*document)[1].position[0].site == 6)
	AssertTrue(t, (*document)[2].position[0].site == 5)

	// Insert at the end
	document.docInsert(3, Element{"end", Position{Identifier{0, 7}}})
	AssertTrue(t, len(*document) == 4)
	AssertTrue(t, (*document)[3].data == "end")

	document.docDelete(2);
	AssertTrue(t, len(*document) == 3)
	AssertTrue(t, (*document)[2].data == "end")
	AssertTrue(t, len((*document)[2].position) == 1)
	AssertTrue(t, IdentifierEquals((*document)[2].position[0], Identifier{0, 7}))

	document.docDelete(0);
	AssertTrue(t, len(*document) == 2)
	AssertTrue(t, (*document)[0].data == "data2")
	AssertTrue(t, len((*document)[0].position) == 1)
	AssertTrue(t, IdentifierEquals((*document)[0].position[0], Identifier{0, 6}))

	document.docDelete(1);
	AssertTrue(t, len(*document) == 1)
	AssertTrue(t, (*document)[0].data == "data2")
	AssertTrue(t, len((*document)[0].position) == 1)
	AssertTrue(t, IdentifierEquals((*document)[0].position[0], Identifier{0, 6}))

	document.docDelete(0);
	AssertTrue(t, len(*document) == 0)
}

func TestDocCreation(t *testing.T) {
	document := NewDocument()

	AssertTrue(t, len(*document) == 2)
	AssertTrue(t, (*document)[0].position[0].pos == 0)
	AssertTrue(t, (*document)[1].position[0].pos == math.MaxInt32)
}

func InsertAtTop(text string) *Document {
	document := NewDocument()
	for _, character := range text {
		document.InsertAt(string(character), 0, utils.RandBetween(1, 5))
	}
	return document
}

func InsertAtBottom(text string) *Document {
	document := NewDocument()
	for index, character := range text {
		document.InsertAt(string(character), index, utils.RandBetween(1, 5))
	}
	return document
}

func TestDocumentInsertAt(t *testing.T) {

	// #1
	text := "Hi everyone!"
	document := InsertAtBottom(text)
	AssertTrue(t, document.ToString() == text)

	// #2
	text = "Hello again!"
	document = InsertAtTop(utils.Reverse(text))
	AssertTrue(t, document.ToString() == text)

	// #3
	text = "Hello!"
	document = NewDocument()
	document.InsertAt("e", 0, 1)
	document.InsertAt("l", 1, 4)
	document.InsertAt("o", 2, 3)
	document.InsertAt("l", 1, 1)
	document.InsertAt("!", 4, 2)
	document.InsertAt("H", 0, 4)
	AssertTrue(t, document.ToString() == text)
}

func TestDocumentDeleteAt(t *testing.T) {

	// #1
	text := "Hi everyone!"
	document := InsertAtBottom(text)
	AssertTrue(t, document.GetLength() == 12)
	AssertTrue(t, document.ToString() == text)

	document.DeleteAt(0)
	document.DeleteAt(0)
	AssertTrue(t, document.GetLength() == 10)
	AssertTrue(t, document.ToString() == " everyone!")

	document.InsertAt("H", 0, 4)
	document.InsertAt("e", 1, 1)
	document.InsertAt("l", 2, 4)
	document.InsertAt("l", 3, 1)
	document.InsertAt("o", 4, 1)
	AssertTrue(t, document.GetLength() == 15)
	AssertTrue(t, document.ToString() == "Hello everyone!")

	// #2
	document = NewDocument()
	document.InsertAt("H", 0, 1)
	document.InsertAt("i", 1, 4)
	document.InsertAt(" ", 2, 1)
	document.InsertAt("e", 3, 4)
	document.InsertAt("v", 4, 1)
	document.InsertAt("e", 5, 4)
	document.InsertAt("r", 6, 1)
	document.InsertAt("y", 7, 4)
	document.InsertAt("o", 8, 1)
	document.InsertAt("n", 9, 4)
	document.InsertAt("e", 10, 1)
	document.InsertAt("!", 11, 1)
	AssertTrue(t, document.GetLength() == 12)
	AssertTrue(t, document.ToString() == "Hi everyone!")

	document.DeleteAt(3)
	document.DeleteAt(3)
	document.DeleteAt(3)
	document.DeleteAt(3)
	document.DeleteAt(3)
	document.DeleteAt(3)
	document.DeleteAt(3)
	document.DeleteAt(3)
	AssertTrue(t, document.GetLength() == 4)
	AssertTrue(t, document.ToString() == "Hi !")

	document.InsertAt("f", 3, 4)
	document.InsertAt("o", 4, 1)
	document.InsertAt("l", 5, 4)
	document.InsertAt("k", 6, 1)
	document.InsertAt("s", 7, 1)
	AssertTrue(t, document.GetLength() == 9)
	AssertTrue(t, document.ToString() == "Hi folks!")
}

func TestDocInsertAtPos(t *testing.T) {
	// #1
	document := NewDocument()

	document.InsertAtPos(Position{Identifier{1, 1}}, "H")
	document.InsertAtPos(Position{Identifier{2, 1}}, "i")
	document.InsertAtPos(Position{Identifier{3, 1}}, " ")
	document.InsertAtPos(Position{Identifier{4, 1}}, "e")
	document.InsertAtPos(Position{Identifier{5, 1}}, "v")
	document.InsertAtPos(Position{Identifier{6, 1}}, "e")
	document.InsertAtPos(Position{Identifier{7, 1}}, "r")
	document.InsertAtPos(Position{Identifier{8, 1}}, "y")
	document.InsertAtPos(Position{Identifier{9, 1}}, "o")
	document.InsertAtPos(Position{Identifier{10, 1}}, "n")
	document.InsertAtPos(Position{Identifier{11, 1}}, "e")
	document.InsertAtPos(Position{Identifier{12, 1}}, "!")
	AssertTrue(t, document.GetLength() == 12)
	AssertTrue(t, document.ToString() == "Hi everyone!")

	// #2
	document = NewDocument()

	document.InsertAtPos(Position{Identifier{4, 1}}, "e")
	document.InsertAtPos(Position{Identifier{1, 1}}, "H")
	document.InsertAtPos(Position{Identifier{8, 1}}, "y")
	document.InsertAtPos(Position{Identifier{9, 1}}, "o")
	document.InsertAtPos(Position{Identifier{7, 1}}, "r")
	document.InsertAtPos(Position{Identifier{6, 1}}, "e")
	document.InsertAtPos(Position{Identifier{2, 1}}, "i")
	document.InsertAtPos(Position{Identifier{3, 1}}, " ")
	document.InsertAtPos(Position{Identifier{5, 1}}, "v")
	document.InsertAtPos(Position{Identifier{12, 1}}, "!")
	document.InsertAtPos(Position{Identifier{10, 1}}, "n")
	document.InsertAtPos(Position{Identifier{11, 1}}, "e")
	AssertTrue(t, document.GetLength() == 12)
	AssertTrue(t, document.ToString() == "Hi everyone!")

	// #3
	document = NewDocument()

	document.InsertAtPos(Position{Identifier{1, 1}}, "H")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}}, "e")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "l")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "l")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "o")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "!")
	AssertTrue(t, document.GetLength() == 6)
	AssertTrue(t, document.ToString() == "Hello!")

	// #4
	document = NewDocument()

	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "l")
	document.InsertAtPos(Position{Identifier{1, 1}}, "H")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "!")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "l")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}}, "e")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "o")
	AssertTrue(t, document.GetLength() == 6)
	AssertTrue(t, document.ToString() == "Hello!")
}

func TestDocDeleteAtPos(t *testing.T) {
	document := NewDocument()

	// #1
	document.InsertAtPos(Position{Identifier{1, 1}}, "H")
	document.InsertAtPos(Position{Identifier{2, 1}}, "i")
	document.InsertAtPos(Position{Identifier{3, 1}}, " ")
	document.InsertAtPos(Position{Identifier{4, 1}}, "e")
	document.InsertAtPos(Position{Identifier{5, 1}}, "v")
	document.InsertAtPos(Position{Identifier{6, 1}}, "e")
	document.InsertAtPos(Position{Identifier{7, 1}}, "r")
	document.InsertAtPos(Position{Identifier{8, 1}}, "y")
	document.InsertAtPos(Position{Identifier{9, 1}}, "o")
	document.InsertAtPos(Position{Identifier{10, 1}}, "n")
	document.InsertAtPos(Position{Identifier{11, 1}}, "e")
	document.InsertAtPos(Position{Identifier{12, 1}}, "!")
	AssertTrue(t, document.GetLength() == 12)
	AssertTrue(t, document.ToString() == "Hi everyone!")

	document.DeleteAtPos(Position{Identifier{4, 1}})
	document.DeleteAtPos(Position{Identifier{6, 1}})
	document.DeleteAtPos(Position{Identifier{10, 1}})
	document.DeleteAtPos(Position{Identifier{7, 1}})
	document.DeleteAtPos(Position{Identifier{5, 1}})
	document.DeleteAtPos(Position{Identifier{3, 1}})
	document.DeleteAtPos(Position{Identifier{11, 1}})
	document.DeleteAtPos(Position{Identifier{8, 1}})
	document.DeleteAtPos(Position{Identifier{9, 1}})
	AssertTrue(t, document.GetLength() == 3)
	AssertTrue(t, document.ToString() == "Hi!")

	// #2
	document = NewDocument()

	document.InsertAtPos(Position{Identifier{1, 1}}, "H")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}}, "e")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "l")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "l")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "o")
	document.InsertAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}}, "!")
	AssertTrue(t, document.GetLength() == 6)
	AssertTrue(t, document.ToString() == "Hello!")

	document.DeleteAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}})
	document.DeleteAtPos(Position{Identifier{1, 1}})
	document.DeleteAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}})
	AssertTrue(t, document.GetLength() == 3)
	AssertTrue(t, document.ToString() == "elo")

	document.DeleteAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}})
	document.DeleteAtPos(Position{Identifier{1, 1}, Identifier{1, 1}})
	document.DeleteAtPos(Position{Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}, Identifier{1, 1}})
	AssertTrue(t, document.GetLength() == 0)
	AssertTrue(t, document.ToString() == "")
}

func TestPositionConvertion(t *testing.T) {
	str := "(1298498082,1)"
	position := ToPosition(str)
	AssertTrue(t, position.ToString() == str)

	
	str = "(1298498082,1)(123,6)"
	position = ToPosition(str)
	AssertTrue(t, position.ToString() == str)

	
	str = ""
	position = ToPosition(str)
	AssertTrue(t, position.ToString() == str)
}