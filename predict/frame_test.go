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
		depth        int
		allowedMoves []game.Direction
	}{
		{"afraid_to_eat", 7, []game.Direction{game.Up}},
		{"no_mercy", 7, []game.Direction{game.Right}},
		{"enemy_ate", 7, []game.Direction{game.Down}},
		{"leave_hazard", 7, []game.Direction{game.Left}},
		{"over_chase", 7, []game.Direction{game.Right}},
		{"corner_crash", 7, []game.Direction{game.Left}},
		{"bad_joust", 7, []game.Direction{game.Down}},
		{"pessimistic", 7, []game.Direction{game.Left}},
		{"hungry_hazard", 7, []game.Direction{game.Right}},
		{"wont_eat", 7, []game.Direction{game.Up}},
		{"risky_food", 7, []game.Direction{game.Down}},
	}

	for _, tt := range testCases {
		t.Run(tt.frame, func(t *testing.T) {
			input, ok := testdata.Frame(tt.frame)
			if !ok {
				t.Fatalf("Failed to find frame %s", tt.frame)
			}

			mv := tree.Search(testTimeout(),
				&input,
				&input.You,
				&predict.State{},
				tree.ConfigMaxDepth(tt.depth),
			)

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
