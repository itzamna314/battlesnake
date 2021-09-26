package guess

import "github.com/itzamna314/battlesnake/game"

type CoordSet []Coord

func (g *CoordSet) Set(c *game.Coord, p float64) {
	for i, guess := range *g {
		if guess.Hit(c) {
			(*g)[i].Probability = clamp(p)
			return
		}
	}

	*g = append(*g, Coord{*c, p})
}

func (g *CoordSet) Add(c *game.Coord, p float64) {
	for i, guess := range *g {
		if guess.Hit(c) {
			(*g)[i].Probability = clamp(guess.Probability + p)
			return
		}
	}

	*g = append(*g, Coord{*c, p})
}

func (g *CoordSet) Mult(c *game.Coord, p float64) float64 {
	for i, guess := range *g {
		if guess.Hit(c) {
			(*g)[i].Probability = clamp(guess.Probability * p)
			return (*g)[i].Probability
		}
	}

	return Impossible
}

func (g *CoordSet) Clear(c *game.Coord) {
	for i, guess := range *g {
		if guess.Hit(c) {
			*g = append((*g)[:i], (*g)[i+1:]...)
			return
		}
	}
}

func (g CoordSet) Prob(c *game.Coord) float64 {
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

func (v CoordSet) Clone() CoordSet {
	clone := make(CoordSet, len(v))

	for i := 0; i < len(v); i++ {
		clone[i] = v[i]
	}

	return clone
}
