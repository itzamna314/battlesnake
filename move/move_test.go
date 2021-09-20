package move_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/itzamna314/battlesnake/api"
	"github.com/itzamna314/battlesnake/model"
	"github.com/itzamna314/battlesnake/move"
)

func TestEatOne(t *testing.T) {
	// Arrange
	me := model.Battlesnake{
		// Length 3, facing right
		Head: model.Coord{X: 2, Y: 0},
		Body: []model.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := model.GameState{
		Board: model.Board{
			Snakes: []model.Battlesnake{me},
			Food: []model.Coord{
				{2, 1},
			},
		},
		You: me,
	}

	nextMove := move.Next(state)
	if nextMove.Move != "up" {
		t.Errorf("snake did not eat food, went %s", nextMove.Move)
	}
}

func TestEatTwo(t *testing.T) {
	// Arrange
	me := model.Battlesnake{
		// Length 3, facing right
		Head: model.Coord{X: 2, Y: 0},
		Body: []model.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := model.GameState{
		Board: model.Board{
			Snakes: []model.Battlesnake{me},
			Food: []model.Coord{
				{2, 1},
				{3, 0},
				{4, 0},
			},
		},
		You: me,
	}

	nextMove := move.Next(state)
	if nextMove.Move != "right" {
		t.Errorf("snake did not eat 2 food, went %s", nextMove.Move)
	}
}

func TestNoCrash(t *testing.T) {
	testServer := server()
	defer testServer.Close()

	result, err := play(t, testServer.URL)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if result.Turn < 25 {
		t.Errorf("Expected to live minimum 25 turns, lived for %d", result.Turn)
	}
}

func server() *httptest.Server {
	hnd := api.Build()
	svr := httptest.NewServer(hnd)

	return svr
}

type playResult struct {
	Turn int
}

// play runs a solo game on a 5x5 grid, with no food
func play(t *testing.T, url string) (*playResult, error) {
	cmd := exec.Command("battlesnake", "play", "-W", "5", "-H", "5", "--name", "test", "--url", url, "-g", "solo", "--foodSpawnChance", "0")

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	// Execute the command
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	output := stdBuffer.String()
	t.Log(output)

	bufLines := strings.Split(output, "\n")
	outLine := bufLines[len(bufLines)-2]

	re := regexp.MustCompile(`\[DONE\]: Game completed after (\d+) turns.`)
	matches := re.FindStringSubmatch(outLine)
	if len(matches) != 2 {
		return nil, fmt.Errorf("Unexpected regexp match %v", matches)
	}

	numTurns, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse num turns %s as int: %s", matches[1], err)
	}
	return &playResult{numTurns}, nil
}
