package model

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

func (c *Coord) StepToward(other *Coord) Direction {
	var (
		xDiff, yDiff       = other.X - c.X, other.Y - c.Y
		xDistRaw, yDistRaw = xDiff * xDiff, yDiff * yDiff
	)

	fmt.Printf("Step from %s toward %s\n", *c, *other)
	fmt.Printf("XDiff: %d, XDist: %d, YDiff: %d, YDist: %d\n",
		xDiff, xDistRaw, yDiff, yDistRaw)

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
	Safe   bool
	Weight float64
	Shout  string
}

type PossibleMoves [4]*MoveCoord

func Options(myHead *Coord) PossibleMoves {
	var opts [4]*MoveCoord

	opts[Up] = &MoveCoord{
		Coord: Coord{myHead.X, myHead.Y + 1},
		Safe:  true,
	}
	opts[Down] = &MoveCoord{
		Coord: Coord{myHead.X, myHead.Y - 1},
		Safe:  true,
	}
	opts[Left] = &MoveCoord{
		Coord: Coord{myHead.X - 1, myHead.Y},
		Safe:  true,
	}
	opts[Right] = &MoveCoord{
		Coord: Coord{myHead.X + 1, myHead.Y},
		Safe:  true,
	}

	return opts
}
