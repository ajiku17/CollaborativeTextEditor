package crdt

const BASE = 10

type Identifier struct {
	pos  int
	site int
}

type Position []Identifier

func PositionToNumber(pos Position) Number {
	num := make(Number, len(pos))
	for i := 0; i < len(pos); i++ {
		num = append(num, pos[i].pos)
	}
	return num
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
	return pos
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

func AllocPosition(prevPos Position, afterPos Position) Position {
	index := 0
	interval := 0
	for interval < 1 {
		index++

		interval = PositionSubtract(prefix(afterPos, index), prefix(prevPos, index)) - 1
	}
	step := min(BASE, interval)

	position := PositionAddInt(prefix(prevPos, index), randBetween(0, step)+1)

	return position
}

func TestPrefix() {
	var position Position
	zero := Identifier{0, 0}

	for i := 0; i < 10; i++ {
		position = append(position, Identifier{pos: i, site: i})
	}

	pref := prefix(position, 5)
	if len(pref) != 5 {
		panic("invalid size")
	}

	for i := 0; i < 5; i++ {
		if pref[i] != position[i] {
			panic("aaaaa!")
		}
	}

	pref = prefix(position, 12)
	if len(pref) != 12 {
		panic("invalid size")
	}

	for i := 0; i < 10; i++ {
		if pref[i] != position[i] {
			panic("aaaaa!")
		}
	}

	if pref[10] != zero {
		panic("aaaaa!")
	}

	if pref[11] != zero {
		panic("aaaaa!")
	}

	pref = prefix(position, 0)
	if len(pref) != 0 {
		panic("invalid size")
	}
}
