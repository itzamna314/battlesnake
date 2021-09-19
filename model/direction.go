package model

type Direction uint8

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "up"
	case Down:
		return "down"
	case Left:
		return "left"
	case Right:
		return "right"
	}

	return "down"
}
