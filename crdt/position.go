package crdt

import (
	"github.com/utils"
)

type Identifier struct {
	pos  int
	site int
}

type Position []Identifier

func PositionToNumber(pos Position) Number {
	num := make(Number, len(pos))
	for i := 0; i < len(pos); i++ {
		num[i] = pos[i].pos
	}
	return num
}

func PositionSubtract(pos1, pos2 Position) Number {
	// fmt.Printf("%v\n", pos1)
	// fmt.Printf("%v\n", pos2)
	num1 := PositionToNumber(pos1)
	num2 := PositionToNumber(pos2)

	return NumberSubtract(num1, num2)
}

func PositionAdd(pos1, pos2 Position) Number {
	num1 := PositionToNumber(pos1)
	num2 := PositionToNumber(pos2)

	return NumberAdd(num1, num2)
}

func PositionAddInt(pos Position, val int) Position {
	identifier := pos[len(pos)-1]
	return append(pos, Identifier{identifier.pos + val, identifier.site})
}

func prefix(position Position, index int) Number {
	var numberCopy Number

	for i := 0; i < index; i++ {
		if i < len(position) {
			numberCopy = append(numberCopy, position[i].pos)
		} else {
			numberCopy = append(numberCopy, 0)
		}
	}
	return numberCopy
}

func constructPosition(r Number, prevPos, afterPos Position, site int) Position {
	var res Position;

	for i, digit := range r {
		var s int
		
		if i == len(r) - 1 {
			s = site
		} else if i < len(prevPos) && digit == prevPos[i].pos {
			s = prevPos[i].site
		} else if i < len(afterPos) && digit == afterPos[i].pos{
			s = afterPos[i].site
		} else {
			s = site
		}

		res = append(res, Identifier{digit, s})
	}

	return res
}

func AllocPosition(prevPos Position, afterPos Position, site int) Position {
	index := 0
	interval := 0
	for interval < 1 {
		index++
		interval = NumberToInt(NumberSubtract(prefix(afterPos, index), prefix(prevPos, index))) - 1
	}
	step := utils.Min(BASE, interval)

	r := prefix(prevPos, index)
	position := constructPosition(NumberAdd(r, Number{utils.RandBetween(0, step) + 1}), prevPos, afterPos, site)

	return position
}
