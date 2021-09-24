package model

import "fmt"

type GuessCoord struct {
	Coord
	Probability float64
}

func (g *GuessCoord) String() string {
	return fmt.Sprintf("%s[%4f]", g.Coord, g.Probability)
}

type GuessCoordSet []GuessCoord

func (g *GuessCoordSet) Set(c *Coord, p float64) {
	for i, guess := range *g {
		if guess.Hit(c) {
			(*g)[i].Probability = clamp(p)
			return
		}
	}

	*g = append(*g, GuessCoord{*c, p})
}

func (g *GuessCoordSet) Add(c *Coord, p float64) {
	for i, guess := range *g {
		if guess.Hit(c) {
			(*g)[i].Probability = clamp(guess.Probability + p)
			return
		}
	}

	*g = append(*g, GuessCoord{*c, p})
}

func (g *GuessCoordSet) Mult(c *Coord, p float64) {
	for i, guess := range *g {
		if guess.Hit(c) {
			(*g)[i].Probability = clamp(guess.Probability * p)
			return
		}
	}
}

func (g *GuessCoordSet) Clear(c *Coord) {
	for i, guess := range *g {
		if guess.Hit(c) {
			*g = append((*g)[:i], (*g)[i+1:]...)
			return
		}
	}
}

func (g GuessCoordSet) Prob(c *Coord) float64 {
	for _, guess := range g {
		if guess.Hit(c) {
			return guess.Probability
		}
	}

	return Impossible
}

func clamp(p float64) float64 {
	if p < Impossible {
		return Impossible
	} else if p > Certain {
		return Certain
	}

	return p
}

func (v GuessCoordSet) Clone() GuessCoordSet {
	clone := make(GuessCoordSet, len(v))

	for i := 0; i < len(v); i++ {
		clone[i] = v[i]
	}

	return clone
}

type SnakeVision []GuessCoordSet

func (v SnakeVision) Clone() SnakeVision {
	clone := make(SnakeVision, len(v))
	for x := 0; x < len(v); x++ {
		clone[x] = v[x].Clone()
	}

	return clone
}
