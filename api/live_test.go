package api_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/itzamna314/battlesnake/api"
)

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
	cmd := exec.Command("battlesnake", "play", "-W", "5", "-H", "5", "--name", "test", "--url", url, "-g", "solo", "--foodSpawnChance", "0", "-v")

	var buf bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = &buf

	// Execute the command
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	output := buf.String()
	t.Logf(output)

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
