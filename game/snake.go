package game

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

func (b *Battlesnake) MoveDet(dir Direction, ate bool) {
	b.moveBody(ate)

	// Step the head in direction, and copy to body
	b.Head = b.Head.Step(dir)
	b.Body[0] = b.Head
}

func (b *Battlesnake) MoveProb(ate bool) {
	b.moveBody(ate)

	// We don't know where the head went
	// Remove from deterministic structure
	b.Body = b.Body[1:]
}

func (b *Battlesnake) moveBody(ate bool) {
	// If ate, grow tail
	if ate {
		b.Body = append(b.Body, b.Body[len(b.Body)-1])
	}

	// Copy each body segment to next
	// Head will remain copied into neck
	for i := len(b.Body) - 1; i > 0; i-- {
		next := i - 1
		b.Body[i] = b.Body[next]
	}
}
