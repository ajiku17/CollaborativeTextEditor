package crdt

import "testing"

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
