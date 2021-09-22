package model

import "bytes"

type VisionCell struct {
	Food    float64
	Enemies []float64
}

type BoardVision [][]VisionCell

func (v BoardVision) Clone() BoardVision {
	clone := make(BoardVision, len(v))
	for x := 0; x < len(v); x++ {
		clone[x] = make([]VisionCell, len(v[x]))

		for y := 0; y < len(v[x]); y++ {
			clone[x][y] = v[x][y]
		}
	}

	return clone
}

func (v BoardVision) String() string {
	if len(v) == 0 {
		return ""
	}

	var (
		out    bytes.Buffer
		width  = len(v)
		height = len(v[0])
	)

	// Print top to bottom
	for y := height - 1; y >= 0; y-- {
	NextCoord:
		// Print left to right
		for x := 0; x < width; x++ {
			if v[x][y].Food != 0 {
				out.WriteString("+")
				continue
			}

			for _, e := range v[x][y].Enemies {
				if e != 0 {
					out.WriteString("X")
					continue NextCoord
				}
			}

			out.WriteString("O")
		}
		out.WriteString("\n")
	}

	return out.String()
}
