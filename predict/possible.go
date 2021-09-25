package predict

import (
	"github.com/itzamna314/battlesnake/game"
)

type MoveCoord struct {
	game.Coord
	Weight float64
	Shout  string
}

type PossibleMoves [4]*MoveCoord

func Options(myHead *game.Coord) PossibleMoves {
	var opts [4]*MoveCoord

	opts[game.Up] = &MoveCoord{
		Coord: game.Coord{myHead.X, myHead.Y + 1},
	}
	opts[game.Down] = &MoveCoord{
		Coord: game.Coord{myHead.X, myHead.Y - 1},
	}
	opts[game.Left] = &MoveCoord{
		Coord: game.Coord{myHead.X - 1, myHead.Y},
	}
	opts[game.Right] = &MoveCoord{
		Coord: game.Coord{myHead.X + 1, myHead.Y},
	}

	return opts
}
