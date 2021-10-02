package predict_test

import (
	"testing"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
	"github.com/itzamna314/battlesnake/predict"
	"github.com/itzamna314/battlesnake/testdata"
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
		t.Errorf("Expected tail probability to be Certain (1.0), was %v", tail)
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

func TestMoveEnemies(t *testing.T) {
	// Pick a frame with some challenging enemy movement
	initial, _ := testdata.Frame("juke_early")

	var ps predict.State
	ps.Init(&initial)

	ps.MoveEnemies(&initial.You)

	// Assert moves on 'Untimely Neglected Wearable'
	var (
		unwIdx   int
		unwSnake *game.Battlesnake
	)
	for i, snake := range ps.Board.Snakes {
		if snake.Name == "Untimely Neglected Wearable" {
			unwIdx = i
			unwSnake = &ps.Board.Snakes[i]
			break
		}
	}

	t.Logf("Asserting movement for snake %s (%d)", ps.Board.Snakes[unwIdx].Name, unwIdx)

	head := ps.HeadGuesses[unwIdx]
	if len(head) != 1 {
		t.Fatalf("Expected 1 head guess for %s: (3,5)", unwSnake.Name)
	}

	hProb := head.Prob(&game.Coord{3, 5})

	if hProb != 1 {
		t.Errorf("Expected head guess (3,5) with p %.2f, was at %s\n", hProb, &head[0])
	}
}

func TestMoveEnemiesAroundYou(t *testing.T) {
	// Pick a frame with some challenging enemy movement
	initial, _ := testdata.Frame("tight_spot")

	var ps predict.State
	ps.Init(&initial)

	var (
		enemyIdx int
		prob     float64
	)
	for i, snake := range ps.Board.Snakes {
		if snake.ID == ps.You.ID {
			continue
		}

		enemyIdx = i
	}

	// First move - we go right
	ps.MoveEnemies(&ps.You)
	ps.Move(&ps.You, game.Right)

	// Second move - keep going right
	ps.MoveEnemies(&ps.You)
	ps.Move(&ps.You, game.Right)

	prob = ps.HeadGuesses[enemyIdx].Prob(&game.Coord{6, 1})
	if prob != 0 {
		t.Errorf("Move 2 expected probability 0 at (6,1), got %v", prob)
	}

	// Third move - keep going right
	ps.MoveEnemies(&ps.You)
	ps.Move(&ps.You, game.Right)

	prob = ps.HeadGuesses[enemyIdx].Prob(&game.Coord{7, 1})
	if prob != 0 {
		t.Errorf("Move 3 expected probability 0 at (7,1), got %v", prob)
	}

	// Fourth move - go down
	ps.MoveEnemies(&ps.You)
	ps.Move(&ps.You, game.Down)

	prob = ps.HeadGuesses[enemyIdx].Prob(&game.Coord{8, 1})
	if prob != 0 {
		t.Errorf("Move 4 expected probability 0 at (8,1), got %v", prob)
	}

	// Fifth move - move enemies
	// With this path, we blocked the enemy from possibly reaching (9,0)
	ps.MoveEnemies(&ps.You)
	ps.Move(&ps.You, game.Left)

	prob = ps.HeadGuesses[enemyIdx].Prob(&game.Coord{9, 1})
	if prob != 0 {
		t.Errorf("Move 5 expected probability 0 at (9,1), got %v", prob)
	}

	ps.MoveEnemies(&ps.You)

	prob = ps.HeadGuesses[enemyIdx].Prob(&game.Coord{9, 0})
	if prob != 0 {
		t.Errorf("Move 6 expected probability 0 at (9,0), got %v", prob)
	}
}

func TestEnemyAte(t *testing.T) {
	// Pick a frame where an enemy may eat in 2 turns
	initial, _ := testdata.Frame("enemy_ate2")

	var ps predict.State
	ps.Init(&initial)

	var rufio int
	for i, s := range ps.Board.Snakes {
		if s.Name == "Rufio the Tenacious" {
			rufio = i
			break
		}
	}

	// First move - we go down
	ps.MoveEnemies(&ps.You)

	lProb := ps.HeadGuesses[rufio].Prob(&game.Coord{6, 8})
	if lProb != float64(1)/float64(3) {
		t.Errorf("Expected 1/3 prob at (6,8), got %v", lProb)
	}

	uProb := ps.HeadGuesses[rufio].Prob(&game.Coord{7, 9})
	if uProb != float64(1)/float64(3) {
		t.Errorf("Expected 1/3 prob at (7,9), got %v", uProb)
	}

	rProb := ps.HeadGuesses[rufio].Prob(&game.Coord{8, 8})
	if rProb != float64(1)/float64(3) {
		t.Errorf("Expected 1/3 prob at (8,8), got %v", rProb)
	}

	tailProb := ps.BodyGuesses[rufio].Prob(&game.Coord{3, 7})
	if tailProb != guess.Impossible {
		t.Errorf("Expected tail to move when no food available")
	}

	ps.Move(&ps.You, game.Down)

	// Second move - we go left
	ps.MoveEnemies(&ps.You)

	ateTailProb := ps.BodyGuesses[rufio].Prob(&game.Coord{4, 7})
	if ateTailProb == guess.Impossible {
		t.Errorf("Expected non-zero tail when head maybe ate")
	}

	// TODO: head prob should go to 0 where we collided with head guess

	ps.Move(&ps.You, game.Left)

	// Third move - we are dead if rufio ate
	// We should see the same tail at (5,7) that we saw at (4,7) last turn
	ps.MoveEnemies(&ps.You)

	tailProb = ps.BodyGuesses[rufio].Prob(&game.Coord{5, 7})
	if tailProb != ateTailProb {
		t.Errorf("Did not shift possible tail (4,7) to (5,7). Expected %.2f got %.2f", ateTailProb, tailProb)
	}

	tailProb = ps.BodyGuesses[rufio].Prob(&game.Coord{4, 7})
	if tailProb != guess.Impossible {
		t.Errorf("Expected possible tail at (4,7) to clear")
	}
}
