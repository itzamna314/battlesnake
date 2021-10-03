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
		{"leave_hazard", 7, []game.Direction{game.Up, game.Left}},
		{"over_chase", 7, []game.Direction{game.Right}},
		{"corner_crash", 7, []game.Direction{game.Left}},
		{"pessimistic", 7, []game.Direction{game.Left, game.Up}},
		{"hungry_hazard", 8, []game.Direction{game.Right}},
		{"wont_eat", 8, []game.Direction{game.Up}},
		{"risky_food", 7, []game.Direction{game.Left, game.Down}},
		{"juke_early", 7, []game.Direction{game.Down}},
		{"wont_eat2", 3, []game.Direction{game.Left}},
		{"tight_spot", 8, []game.Direction{game.Right}},
		{"enemy_ate2", 6, []game.Direction{game.Right}},
		{"get_long", 7, []game.Direction{game.Up}},
		{"scared_hazard", 8, []game.Direction{game.Right, game.Up}},
		{"wheres_my_butt", 9, []game.Direction{game.Left, game.Up}},
		{"enemy_behind_butt", 8, []game.Direction{game.Down}},
		{"corner_of_doom", 8, []game.Direction{game.Up}},
		{"wont_eat_start", 10, []game.Direction{game.Up}},
		{"reckless_thirst", 11, []game.Direction{game.Left}},
	}

	for _, tt := range testCases {
		t.Run(tt.frame, func(t *testing.T) {
			input, ok := testdata.Frame(tt.frame)
			if !ok {
				t.Fatalf("Failed to find frame %s", tt.frame)
			}

			mv, meta := tree.Search(testTimeout(),
				&input,
				&input.You,
				&predict.State{},
				tree.ConfigMaxDepth(tt.depth),
			)

			t.Logf("Search meta %+v\n", meta)

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
