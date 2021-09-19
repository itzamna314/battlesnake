package model

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c *Coord) Hit(other *Coord) bool {
	return c.X == other.X && c.Y == other.Y
}

type MoveCoord struct {
	Coord
	Safe bool
}

func Options(myHead *Coord) [4]MoveCoord {
	var opts [4]MoveCoord

	opts[Up].Coord.X = myHead.X
	opts[Up].Coord.Y = myHead.Y + 1
	opts[Up].Safe = true

	opts[Down].Coord.X = myHead.X
	opts[Down].Coord.Y = myHead.Y - 1
	opts[Down].Safe = true

	opts[Left].Coord.X = myHead.X - 1
	opts[Left].Coord.Y = myHead.Y
	opts[Left].Safe = true

	opts[Right].Coord.X = myHead.X + 1
	opts[Right].Coord.Y = myHead.Y
	opts[Right].Safe = true

	return opts
}
