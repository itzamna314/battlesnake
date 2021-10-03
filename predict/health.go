package predict

func (s *State) weightHealth(snake *Snake) float64 {
	return Base * snake.Health / 100
}
