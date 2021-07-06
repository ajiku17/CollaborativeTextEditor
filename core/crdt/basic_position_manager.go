package crdt

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

const base = math.MaxInt32


type Identifier struct {
	Pos  int
	Site utils.UUID
}

type BasicPositionManager struct {
	SiteID utils.UUID
	Base   int
}

type BasicPosition []Identifier

func NewBasicPositionManager(siteId utils.UUID) *BasicPositionManager {
	gob.Register(BasicPosition{})

	manager := new(BasicPositionManager)
	manager.Base = base
	manager.SiteID = siteId
	NumberSetBase(manager.Base)

	return manager
}

func (manager *BasicPositionManager) GetMaxPosition() Position {
	return BasicPosition{Identifier{manager.Base, manager.SiteID}}
}

func (manager *BasicPositionManager) GetMinPosition() Position {
	return BasicPosition{Identifier{0, manager.SiteID}}
}

func PositionToNumber(pos BasicPosition) Number {
	num := make(Number, len(pos))
	for i := 0; i < len(pos); i++ {
		num[i] = pos[i].Pos
	}
	return num
}

func IdentifierEquals(id1, id2 Identifier) bool {
	return id1.Pos == id2.Pos && id1.Site == id2.Site
}

func IdentifierIsGreaterOrEqual(id1, id2 Identifier) bool {
	if id1.Pos == id2.Pos {
		return id1.Site >= id2.Site
	}

	return id1.Pos >= id2.Pos
}

func IdentifierIsLessThan(id1, id2 Identifier) bool {
	return !IdentifierIsGreaterOrEqual(id1, id2)
}
 
func (manager *BasicPositionManager) PositionsEqual(pos1, pos2 Position) bool {
	basicPos1, ok1 := pos1.(BasicPosition)
	basicPos2, ok2 := pos2.(BasicPosition)
	if ok1 && ok2 {

		if len(basicPos1) != len(basicPos2) {
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

	return false
}

func (manager *BasicPositionManager)PositionIsGreaterOrEqual(pos1, pos2 Position) bool {
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

	return false
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
	return append(pos, Identifier{identifier.Pos + val, identifier.Site})
}

func Prefix(position BasicPosition, index int) Number {
	var numberCopy Number

	for i := 0; i < index; i++ {
		if i < len(position) {
			numberCopy = append(numberCopy, position[i].Pos)
		} else {
			numberCopy = append(numberCopy, 0)
		}
	}
	return numberCopy
}

func ConstructPosition(r Number, prevPos, afterPos BasicPosition, site utils.UUID) BasicPosition {
	var res BasicPosition

	for i, digit := range r {
		var s utils.UUID
		
		if i == len(r) - 1 {
			s = site
		} else if i < len(prevPos) && digit == prevPos[i].Pos {
			s = prevPos[i].Site
		} else if i < len(afterPos) && digit == afterPos[i].Pos{
			s = afterPos[i].Site
		} else {
			s = site
		}

		res = append(res, Identifier{digit, s})
	}

	return res
}

func (manager *BasicPositionManager) AllocPositionBetween(prevPos, afterPos Position) Position {
	prevBasicPos, ok1 := prevPos.(BasicPosition)
	afterBasicPos, ok2 := afterPos.(BasicPosition)
	if ok1 && ok2 {
		index := 0
		interval := 0
		for interval < 1 {
			index++
			interval = NumberToInt(NumberSubtract(Prefix(afterBasicPos, index), Prefix(prevBasicPos, index))) - 1
		}
		step := utils.Min(BASE, interval)

		r := Prefix(prevBasicPos, index)
		position := ConstructPosition(NumberAdd(r, Number{utils.RandBetween(0, step) + 1}), prevBasicPos, afterBasicPos, manager.SiteID)

		return position
	} else {
		log.Fatalf("BasicPositionManager: Invalid position types %T and %T", prevPos, afterBasicPos)
	}

	return nil
}

func ToBasicPosition(position string) Position  {
	result_position := new(BasicPosition)
	i := 0
	for i < len(position) {
		if position[i] == '(' {
			pos, _ := strconv.Atoi(position[1:strings.Index(position, ",")])
			position = position[strings.Index(position, ",") + 1:]
			site, _ := strconv.Atoi(position[:strings.Index(position, ")")])
			position = position[strings.Index(position, ")"):]
			identifier := Identifier{pos, utils.UUID(site)}
			*result_position = append(*result_position, identifier)
			position = position[1:]
		} else {
			i ++
		}
	}
	return *result_position
}

func BasicPositionToString(pos BasicPosition) string {
	res := ""
	for _, identifier := range pos {
		res += fmt.Sprintf("(%d,%d)", identifier.Pos, identifier.Site)
	}
	return res
}