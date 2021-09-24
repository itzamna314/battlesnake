package game

import "fmt"

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

func (c *Coord) Step(dir Direction) Coord {
	out := *c

	switch dir {
	case Up:
		out.Y += 1
	case Down:
		out.Y -= 1
	case Left:
		out.X -= 1
	case Right:
		out.X += 1
	}

	return out
}

func (c *Coord) StepToward(other *Coord) Direction {
	var (
		xDiff, yDiff       = other.X - c.X, other.Y - c.Y
		xDistRaw, yDistRaw = xDiff * xDiff, yDiff * yDiff
	)

	// Move in the farthest dimension first
	if xDistRaw > yDistRaw {
		if xDiff < 0 {
			return Left
		} else {
			return Right
		}
	} else {
		if yDiff < 0 {
			return Down
		} else {
			return Up
		}
	}
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

type MoveCoord struct {
	Coord
	Weight float64
	Shout  string
}

type PossibleMoves [4]*MoveCoord

func Options(myHead *Coord) PossibleMoves {
	var opts [4]*MoveCoord

	opts[Up] = &MoveCoord{
		Coord: Coord{myHead.X, myHead.Y + 1},
	}
	opts[Down] = &MoveCoord{
		Coord: Coord{myHead.X, myHead.Y - 1},
	}
	opts[Left] = &MoveCoord{
		Coord: Coord{myHead.X - 1, myHead.Y},
	}
	opts[Right] = &MoveCoord{
		Coord: Coord{myHead.X + 1, myHead.Y},
	}

	return opts
}
