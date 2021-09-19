package move

import (
	"log"
	"math/rand"

	"github.com/itzamna314/battlesnake/model"
)

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func Next(state model.GameState) model.BattlesnakeMoveResponse {
	var (
		myBody        = state.You.Body
		myHead        = myBody[0]
		possibleMoves = model.Options(&myHead)
	)

	// Step 1 - Don't hit walls.
	// Use information in GameState to prevent your Battlesnake from moving beyond the boundaries of the board.
	// boardWidth := state.Board.Width
	// boardHeight := state.Board.Height
	if myHead.X-1 < 0 {
		possibleMoves[model.Left].Safe = false
	} else if myHead.X+1 == state.Board.Width {
		possibleMoves[model.Right].Safe = false
	}

	if myHead.Y-1 < 0 {
		possibleMoves[model.Down].Safe = false
	} else if myHead.Y+1 == state.Board.Height {
		possibleMoves[model.Up].Safe = false
	}

	// Step 2 - Don't hit yourself.
	// Use information in GameState to prevent your Battlesnake from colliding with itself.
	for dir, poss := range possibleMoves {
		if !poss.Safe {
			continue
		}

		for _, body := range myBody {
			if poss.Coord.Hit(&body) {
				possibleMoves[dir].Safe = false
			}
		}
	}

	// Step 3 - Don't collide with others.
	// Use information in GameState to prevent your Battlesnake from colliding with others.
Enemies:
	for _, enemy := range state.Board.Snakes {
		for eIdx, eBody := range enemy.Body {
			isHead := eIdx == 0
			for dir, poss := range possibleMoves {
				if !poss.Safe {
					continue
				}

				// KILL
				if isHead && enemy.Length < state.You.Length {
					possibleMoves[dir].Safe = true
					possibleMoves[dir].Weight = 1
					break Enemies
				}

				if poss.Hit(&eBody) {
					possibleMoves[dir].Safe = false
				}
			}
		}
	}

	// Step 4 - Find food.
	// Use information in GameState to seek out and find food.
	var (
		closestFood model.Coord
		minDist     int
	)
	for _, food := range state.Board.Food {
		dist := myHead.Dist(&food)

		if minDist == 0 || dist < minDist {
			minDist = dist
			closestFood = food
		}
	}

	// Prefer to move toward the nearest food
	if minDist > 0 {
		step := myHead.StepToward(&closestFood)
		possibleMoves[step].Weight = 0.75
	}

	// Finally, choose a move from the available safe moves.
	// TODO: Step 5 - Select a move to make based on strategy, rather than random.
	var (
		nextMove  string
		safeMoves []string
		maxWeight float64
	)
	for dir, coord := range possibleMoves {
		if !coord.Safe {
			continue
		}

		if coord.Weight > maxWeight {
			safeMoves = []string{model.Direction(dir).String()}
		} else if coord.Weight == maxWeight {
			safeMoves = append(safeMoves, model.Direction(dir).String())
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
		log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	}
	return model.BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
