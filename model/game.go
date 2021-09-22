package model

type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Timeout int32   `json:"timeout"`
}

type Ruleset struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Battlesnake struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Health  int32   `json:"health"`
	Body    []Coord `json:"body"`
	Head    Coord   `json:"head"`
	Length  int32   `json:"length"`
	Latency string  `json:"latency"`

	// Used in non-standard game modes
	Shout string `json:"shout"`
	Squad string `json:"squad"`
}

func (b *Battlesnake) Clone() Battlesnake {
	clone := Battlesnake{
		ID:      b.ID,
		Name:    b.Name,
		Health:  b.Health,
		Head:    b.Head,
		Length:  b.Length,
		Latency: b.Latency,
		Body:    make([]Coord, len(b.Body)),
	}

	for i, bd := range b.Body {
		clone.Body[i] = bd
	}

	return clone
}
