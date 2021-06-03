package crdt

const BASE = 10
type Number []int

func prependZeroes(n Number, length int) {
	for len(n) < length {
		n = append(Number{0}, n...)
	}
}

/*
 *	precondition: pos1 > pos2 so we don't have negative numbers
 *  pos1 and pos2 can have different lengths
 */
func NumberSubtract(num1, num2 Number) Number {
	prependZeroes(num2, len(num1))

	result := make(Number, len(num1))
	carry := 0
	for i := len(num1) - 1; i >= 0; i-- {
		if num1[i] >= num2[i]+carry {
			result[i] = num1[i] - num2[i]
			carry = 0
		} else {
			result[i] = BASE + num1[i] - (num2[i] + carry)
			carry = 1
		}
	}

	return result
}

func NumberAdd(num1, num2 Number) Number {
	return num1
}
