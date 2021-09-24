package testdata

import (
	"embed"
	"encoding/json"
	"io/fs"
	"strings"

	"github.com/itzamna314/battlesnake/game"
)

var (
	//go:embed frames
	frameFiles   embed.FS
	frameLibrary map[string]game.GameState
)

func init() {
	frameLibrary = make(map[string]game.GameState)

	err := fs.WalkDir(frameFiles, ".", func(path string, de fs.DirEntry, err error) error {
		if de.IsDir() {
			return nil
		}

		name := strings.Split(de.Name(), ".")
		if len(name) != 2 {
			return nil
		}

		bb, err := fs.ReadFile(frameFiles, path)
		if err != nil {
			return err
		}

		var gs game.GameState
		if err := json.Unmarshal(bb, &gs); err != nil {
			return err
		}

		frameLibrary[name[0]] = gs
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func Frame(name string) (game.GameState, bool) {
	gs, ok := frameLibrary[name]
	return gs, ok
}
