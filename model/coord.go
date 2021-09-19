package model

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c *Coord) Hit(other *Coord) bool {
	return c.X == other.X && c.Y == other.Y
}

func (c *Coord) Dist(other *Coord) int {
	xDiff := c.X - other.X
	yDiff := c.Y - other.Y

	return (xDiff * xDiff) + (yDiff * yDiff)
}

func (c *Coord) StepToward(other *Coord) Direction {
	var (
		xDiff = other.X - c.X
		yDiff = other.Y - c.Y
	)

	// Move in the farthest dimension first
	if (xDiff * xDiff) > (yDiff * yDiff) {
		if xDiff < 0 {
			return Left
		} else {
			return Right
		}
	} else {
		if yDiff < 0 {
			return Down
		} else {
			return Right
		}
	}
}

type MoveCoord struct {
	Coord
	Safe   bool
	Weight float64
	Scream string
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
