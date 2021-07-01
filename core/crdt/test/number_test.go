package test

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"testing"
)

func TestNumberToInt(t *testing.T) {
	crdt.NumberSetBase(10)

	AssertTrue(t, crdt.NumberToInt(crdt.Number{1, 2, 3}) == 123)
}

func TestIsLessThan(t *testing.T) {
	crdt.NumberSetBase(64)

	AssertTrue(t, crdt.IsLessThan(crdt.Number{3, 2}, crdt.Number{3, 8}))
	AssertTrue(t, crdt.IsLessThan(crdt.Number{3, 2}, crdt.Number{1, 5, 8}))
	AssertTrue(t, !crdt.IsLessThan(crdt.Number{1, 5, 8}, crdt.Number{3, 2}))
	AssertTrue(t, crdt.IsLessThan(crdt.Number{3, 8}, crdt.Number{5, 6}))
}

func TestNumberAdd(t *testing.T) {
	crdt.NumberSetBase(10)

	// #1
	num1 := crdt.Number{3, 8}
	num2 := crdt.Number{9, 8}

	sum := crdt.NumberAdd(num1, num2)
	AssertTrue(t, crdt.NumberToInt(sum) == crdt.NumberToInt(crdt.Number{1, 3, 6}))

	// #2
	crdt.NumberSetBase(64)
	num1 = crdt.Number{3, 8} // 200
	num2 = crdt.Number{9, 8} // 584

	sum = crdt.NumberAdd(num1, num2) // 784 = 12 * 64 + 16
	AssertTrue(t, crdt.NumberToInt(sum) == crdt.NumberToInt(crdt.Number{12, 16}))

	// #3
	crdt.NumberSetBase(10)
	num1 = crdt.Number{3, 2}
	num2 = crdt.Number{1, 5, 8}

	sum = crdt.NumberAdd(num2, num1)
	AssertTrue(t, crdt.NumberToInt(sum) == crdt.NumberToInt(crdt.Number{1, 9, 0}))

	// #4
	crdt.NumberSetBase(10)
	num1 = crdt.Number{1, 5, 8}
	num2 = crdt.Number{3, 2}

	sum = crdt.NumberAdd(num2, num1)
	AssertTrue(t, crdt.NumberToInt(sum) == crdt.NumberToInt(crdt.Number{1, 9, 0}))
}

func TestNumberSubtract(t *testing.T) {
	crdt.NumberSetBase(10)

	// #1
	num1 := crdt.Number{3, 8}
	num2 := crdt.Number{9, 8}

	sum := crdt.NumberSubtract(num2, num1)
	AssertTrue(t, crdt.NumberToInt(sum) == crdt.NumberToInt(crdt.Number{6, 0}))

	// #2
	crdt.NumberSetBase(64)
	num1 = crdt.Number{3, 8} // 200
	num2 = crdt.Number{9, 8} // 584

	sum = crdt.NumberSubtract(num2, num1) // 384 = 6 * 64
	AssertTrue(t, crdt.NumberToInt(sum) == crdt.NumberToInt(crdt.Number{6, 0}))

	// #3
	crdt.NumberSetBase(64)
	num1 = crdt.Number{3, 8} // 200
	num2 = crdt.Number{5, 6} // 326

	sum = crdt.NumberSubtract(num2, num1) // 126 = 1 * 64 + 62
	AssertTrue(t, crdt.NumberToInt(sum) == crdt.NumberToInt(crdt.Number{1, 62}))

	// #4
	crdt.NumberSetBase(10)
	num1 = crdt.Number{5, 1}
	num2 = crdt.Number{3, 2}

	sum = crdt.NumberSubtract(num1, num2)
	AssertTrue(t, crdt.NumberToInt(sum) == crdt.NumberToInt(crdt.Number{1, 9}))

	// #5
	crdt.NumberSetBase(10)
	num1 = crdt.Number{3, 2}
	num2 = crdt.Number{1, 5, 8}

	sum = crdt.NumberSubtract(num2, num1)
	AssertTrue(t, crdt.NumberToInt(sum) == crdt.NumberToInt(crdt.Number{1, 2, 6}))
}
