package testdata

import (
	"embed"
	"encoding/json"
	"io/fs"
	"strings"

	"github.com/itzamna314/battlesnake/model"
)

var (
	//go:embed frames
	frameFiles   embed.FS
	frameLibrary map[string]model.GameState
)

func init() {
	frameLibrary = make(map[string]model.GameState)

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

		var gs model.GameState
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

func Frame(name string) (model.GameState, bool) {
	gs, ok := frameLibrary[name]
	return gs, ok
}
