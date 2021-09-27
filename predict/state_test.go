package predict_test

import (
	"testing"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
	"github.com/itzamna314/battlesnake/predict"
)

func TestMoveGameState(t *testing.T) {
	me := game.Battlesnake{
		// Length 3, facing right
		ID:     "me",
		Head:   game.Coord{X: 2, Y: 0},
		Body:   []game.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 80,
	}
	enemy := game.Battlesnake{
		// Length 3, facing down
		ID:     "enemy",
		Head:   game.Coord{X: 8, Y: 8},
		Body:   []game.Coord{{X: 8, Y: 8}, {X: 9, Y: 8}, {X: 9, Y: 9}},
		Health: 80,
	}
	input := game.GameState{
		Board: game.Board{
			Height: 10,
			Width:  10,
			Snakes: []game.Battlesnake{me, enemy},
			Food: []game.Coord{
				{2, 1},
			},
		},
		You: me,
	}

	// Clone input to initialize data structures
	var state predict.State
	state.Init(&input)

	// Project me moving up
	state.Move(&me, game.Up)

	// Assert that I'm in the right place
	assertHit(t, &game.Coord{2, 1}, &state.You.Head)
	assertHit(t, &game.Coord{2, 1}, &state.You.Body[0])
	assertHit(t, &game.Coord{2, 0}, &state.You.Body[1])
	assertHit(t, &game.Coord{1, 0}, &state.You.Body[2])

	// Assert that I ate the food
	for _, food := range state.Board.Food {
		if food.Hit(&game.Coord{2, 1}) {
			t.Errorf("Expected food at (2,1) to be eaten by me")
		}
	}

	// Move the enemies
	state.MoveEnemies(&me)

	// Certain remaining body segments
	enemy = state.Board.Snakes[1]
	if len(enemy.Body) != 2 {
		t.Fatalf("Expected enemy to be length 2 due to uncertainty, was %d", len(enemy.Body))
	}

	assertHit(t, &game.Coord{8, 8}, &enemy.Body[0])
	assertHit(t, &game.Coord{9, 8}, &enemy.Body[1])

	// Validate guesses, including certain ones
	if len(state.BodyGuesses) != 2 {
		t.Fatalf("2 snakes required for easy lookup")
	}

	guesses := state.BodyGuesses[1]
	if len(guesses) != 2 {
		t.Errorf("Expected 2 certain body segments. Got %d",
			len(guesses))
	}

	neck := guesses.Prob(&game.Coord{8, 8})
	if neck != guess.Certain {
		t.Errorf("Expected neck probability to be Certain (1.0), was %v", neck)
	}
	tail := guesses.Prob(&game.Coord{9, 8})
	if tail != guess.Certain {
		t.Errorf("Expected tail probability to be Certain (1.0), was %v", neck)
	}

	oneThird := 1.0 / 3.0

	guesses = state.HeadGuesses[1]
	if len(guesses) != 3 {
		t.Errorf("Expected 3 possible head guesses. Got %d",
			len(guesses))
	}

	headUp := guesses.Prob(&game.Coord{8, 9})
	if headUp != oneThird {
		t.Errorf("Expected head up probability to be 1/3, was %v", headUp)
	}
	headLeft := guesses.Prob(&game.Coord{7, 8})
	if headUp != oneThird {
		t.Errorf("Expected head left probability to be 1/3, was %v", headLeft)
	}
	headDown := guesses.Prob(&game.Coord{8, 7})
	if headDown != oneThird {
		t.Errorf("Expected head down probability to be 1/3, was %v", headDown)
	}
}

func assertHit(t *testing.T, exp, act *game.Coord) bool {
	t.Helper()

	hits := exp.Hit(act)
	if !hits {
		t.Errorf("%s at unexpected location. Expected %s", act, exp)
	}

	return hits
}
