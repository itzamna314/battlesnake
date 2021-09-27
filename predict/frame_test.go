package predict_test

import (
	"testing"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/predict"
	"github.com/itzamna314/battlesnake/testdata"
	"github.com/itzamna314/battlesnake/tree"
)

func TestFrames(t *testing.T) {
	testCases := []struct {
		frame        string
		allowedMoves []game.Direction
	}{
		{"afraid_to_eat", []game.Direction{game.Up}},
		{"no_mercy", []game.Direction{game.Right}},
		{"enemy_ate", []game.Direction{game.Down}},
		{"leave_hazard", []game.Direction{game.Left}},
		{"over_chase", []game.Direction{game.Right}},
		{"corner_crash", []game.Direction{game.Left}},
		{"bad_joust", []game.Direction{game.Down}},
		{"pessimistic", []game.Direction{game.Left, game.Up}},
		{"hungry_hazard", []game.Direction{game.Right}},
	}

	for _, tt := range testCases {
		t.Run(tt.frame, func(t *testing.T) {
			input, ok := testdata.Frame(tt.frame)
			if !ok {
				t.Fatalf("Failed to find frame %s", tt.frame)
			}

			mv := tree.Search(testTimeout(), &input, &input.You, &predict.State{})

			for _, allowed := range tt.allowedMoves {
				if allowed == mv {
					t.Logf("Made allowed move %s", mv)
					return
				}
			}

			t.Errorf("Made disallowed move %s. Allowed: %v", mv, tt.allowedMoves)
		})
	}
}
