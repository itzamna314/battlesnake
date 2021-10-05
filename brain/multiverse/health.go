package multiverse

func (s *State) weightHealth(snake *Snake) FloatWeight {
	return Base * FloatWeight(snake.Health) / 100
}
