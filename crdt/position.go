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

func IdentifierEquals(id1, id2 Identifier) bool {
	return id1.pos == id2.pos && id1.site == id2.site
}

func IdentifierIsGreaterOrEqual(id1, id2 Identifier) bool {
	if id1.pos == id2.pos {
		return id1.site >= id2.site
	}

	return id1.pos >= id2.pos
}

func IdentifierIsLessThan(id1, id2 Identifier) bool {
	return !IdentifierIsGreaterOrEqual(id1, id2)
}
 
func PositionEquals(pos1, pos2 Position) bool {
	if (len(pos1) != len(pos2)) {
		return false
	}

	for i := 0; i < len(pos1); i++ {
		if !IdentifierEquals(pos1[i], pos2[i]) {
			return false
		}
	}

	return true
}

func PositionIsGreaterOrEqual(pos1, pos2 Position) bool {
	for i := 0; i < utils.Max(len(pos1), len(pos2)); i++ {
		var id1, id2 Identifier

		if i >= len(pos1) {
			id1 = Identifier{}
		} else {
			id1 = pos1[i]
		}

		if i >= len(pos2) {
			id2 = Identifier{}
		} else {
			id2 = pos2[i]
		}

		if !IdentifierEquals(id1, id2) {
			return IdentifierIsGreaterOrEqual(id1, id2)
		}
	}

	return true
}

func PositionIsLessThan(pos1, pos2 Position) bool {
	return !PositionIsGreaterOrEqual(pos1, pos2)
}

func PositionSubtract(pos1, pos2 Position) Number {
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


