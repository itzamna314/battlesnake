package predict

import (
	"github.com/itzamna314/battlesnake/game"
)

func (p *State) moveSnakeBody(snake *game.Battlesnake, ate bool) {
	// If ate, grow tail
	if ate {
		snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])
	}

	// Copy each body segment to next
	// Head will remain copied into neck
	for i := len(snake.Body) - 1; i > 0; i-- {
		next := i - 1
		snake.Body[i] = snake.Body[next]
	}
}
