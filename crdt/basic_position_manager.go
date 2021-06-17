package crdt

import (
	"log"
	"math"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

const base = math.MaxInt32


type Identifier struct {
	pos  int
	site int
}

type BasicPositionManager struct {
	base int
}

type BasicPosition []Identifier

func (manager *BasicPositionManager) PositionManagerInit() {
	manager.base = base
	NumberSetBase(manager.base)
}

func (manager *BasicPositionManager) GetMaxPosition() Position {
	return BasicPosition{Identifier{manager.base, 0}}
}

func (manager *BasicPositionManager) GetMinPosition() Position {
	return BasicPosition{Identifier{0, 0}}
}

func PositionToNumber(pos BasicPosition) Number {
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
 
func (manager *BasicPositionManager) PositionsEqual(pos1, pos2 Position) bool {
	basicPos1, ok1 := pos1.(BasicPosition)
	basicPos2, ok2 := pos2.(BasicPosition)
	if ok1 && ok2 {

		if (len(basicPos1) != len(basicPos2)) {
			return false
		}

		for i := 0; i < len(basicPos1); i++ {
			if !IdentifierEquals(basicPos1[i], basicPos2[i]) {
				return false
			}
		}

		return true
	} else {
		log.Fatalf("BasicPositionManager: Invalid position types %T and %T", pos1, pos2)
	}

	return false;
}

func (mamanger *BasicPositionManager)PositionIsGreaterOrEqual(pos1, pos2 Position) bool {
	basicPos1, ok1 := pos1.(BasicPosition)
	basicPos2, ok2 := pos2.(BasicPosition)
	if ok1 && ok2 {
		for i := 0; i < utils.Max(len(basicPos1), len(basicPos2)); i++ {
			var id1, id2 Identifier

			if i >= len(basicPos1) {
				id1 = Identifier{}
			} else {
				id1 = basicPos1[i]
			}

			if i >= len(basicPos2) {
				id2 = Identifier{}
			} else {
				id2 = basicPos2[i]
			}

			if !IdentifierEquals(id1, id2) {
				return IdentifierIsGreaterOrEqual(id1, id2)
			}
		}

		return true
	} else {
		log.Fatalf("BasicPositionManager: Invalid position types %T and %T", pos1, pos2)
	}

	return false;
}

func (manager *BasicPositionManager)PositionIsLessThan(pos1, pos2 Position) bool {
	return !manager.PositionIsGreaterOrEqual(pos1, pos2)
}

func PositionSubtract(pos1, pos2 BasicPosition) Number {
	num1 := PositionToNumber(pos1)
	num2 := PositionToNumber(pos2)

	return NumberSubtract(num1, num2)
}

func PositionAdd(pos1, pos2 BasicPosition) Number {
	num1 := PositionToNumber(pos1)
	num2 := PositionToNumber(pos2)

	return NumberAdd(num1, num2)
}

func PositionAddInt(pos BasicPosition, val int) BasicPosition {
	identifier := pos[len(pos)-1]
	return append(pos, Identifier{identifier.pos + val, identifier.site})
}

func prefix(position BasicPosition, index int) Number {
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

func constructPosition(r Number, prevPos, afterPos BasicPosition, site int) BasicPosition {
	var res BasicPosition;

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

func (manager *BasicPositionManager) AllocPositionBetween(prevPos, afterPos Position, site int) Position {
	prevBasicPos, ok1 := prevPos.(BasicPosition)
	afterBasicPos, ok2 := afterPos.(BasicPosition)
	if ok1 && ok2 {
		index := 0
		interval := 0
		for interval < 1 {
			index++
			interval = NumberToInt(NumberSubtract(prefix(afterBasicPos, index), prefix(prevBasicPos, index))) - 1
		}
		step := utils.Min(BASE, interval)

		r := prefix(prevBasicPos, index)
		position := constructPosition(NumberAdd(r, Number{utils.RandBetween(0, step) + 1}), prevBasicPos, afterBasicPos, site)

		return position
	} else {
		log.Fatalf("BasicPositionManager: Invalid position types %T and %T", prevPos, afterBasicPos)
	}

	return nil
}


