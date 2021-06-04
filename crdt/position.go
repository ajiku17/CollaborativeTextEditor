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

func prefix(position Position, index int) Position {
	var positionCopy Position

	for i := 0; i < index; i++ {
		if i < len(position) {
			positionCopy = append(positionCopy, position[i])
		} else {
			positionCopy = append(positionCopy, Identifier{})
		}
	}
	return positionCopy
}

func AllocPosition(prevPos Position, afterPos Position, site int) Position {
	index := 0
	interval := 0
	for interval < 1 {
		index++
		interval = NumberToInt(PositionSubtract(prefix(afterPos, index), prefix(prevPos, index))) - 1
	}
	step := utils.Min(BASE, interval)

	position := prefix(prevPos, index)
	position[len(position)-1].pos = utils.RandBetween(0, step) + 1
	position[len(position)-1].site = site

	return position
}
