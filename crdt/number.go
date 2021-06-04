package crdt

import (
	"fmt"

	"github.com/utils"
)

var BASE int

type Number []int

func (number Number) ToString() {
	for _, n := range number {
		fmt.Printf("%d ", n)
	}
	fmt.Printf("\n")
}

func NumberSetBase(base int) {
	BASE = base
}

func prependZeroes(n Number, length int) Number {
	for len(n) < length {
		n = append(Number{0}, n...)
	}
	return n
}

/*
 *
 */
func IsLessThan(num1, num2 Number) bool {
	if BASE == 0 {
		panic("assertion failed: BASE can not be zero. Use NumberSetBase(base)")
	}

	if len(num1) < len(num2) {
		return true
	}

	if len(num1) > len(num2) {
		return false
	}

	for i := 0; i < len(num1); i++ {
		if num1[i] != num2[i] {
			return num1[i] < num2[i]
		}
	}

	return false
}

/*
 *	precondition: pos1 > pos2 so we don't have negative numbers
 *  pos1 and pos2 can have different lengths
 */
func NumberSubtract(num1, num2 Number) Number {
	if BASE == 0 {
		panic("assertion failed: BASE can not be zero. Use NumberSetBase(base)")
	}

	if IsLessThan(num1, num2) {
		panic("assertion failed: num1 is less than num2")
	}

	num2 = prependZeroes(num2, len(num1))

	result := make(Number, len(num1))
	carry := 0
	for i := len(num1) - 1; i >= 0; i-- {
		if num1[i] >= num2[i]+carry {
			result[i] = num1[i] - num2[i] - carry
			carry = 0
		} else {
			result[i] = BASE + num1[i] - (num2[i] + carry)
			carry = 1
		}
	}

	return result
}

func NumberAdd(num1, num2 Number) Number {
	if BASE == 0 {
		panic("assertion failed: BASE can not be zero. Use NumberSetBase(base)")
	}

	size := utils.Max(len(num1), len(num2))

	num1 = prependZeroes(num1, size)
	num2 = prependZeroes(num2, size)

	result := make(Number, size)
	carry := 0
	for i := size - 1; i >= 0; i-- {
		result[i] = (num1[i] + num2[i] + carry) % BASE
		carry = (num1[i] + num2[i] + carry) / BASE
	}

	if carry > 0 {
		result = append(Number{carry}, result...)
	}

	return result
}

func NumberToInt(n Number) int {
	if BASE == 0 {
		panic("assertion failed: BASE can not be zero. Use NumberSetBase(base)")
	}

	res := 0
	for _, ni := range n {
		res = res*BASE + ni
	}

	return res
}
