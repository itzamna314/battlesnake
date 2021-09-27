package predict

import "github.com/itzamna314/battlesnake/game"

const (
	// Certain death means we will crash into a wall,
	// or ourself, or another snake's
	// Use 1 less bit than possible, so that
	// CertainDeath + CertainDeath doesn't wrap around
	CertainDeath int32 = -1 << 30

	// This is the maximum weight we'll put into avoiding enemies
	AvoidEnemies int32 = -1 << 20

	// Neutral means we have no opinion about this coord
	Neutral int32 = 0

	// Apply this weight to food we would eat if we don't want to eat
	FoodNotHungry int32 = 0 // -1 << 4

	// Apply this weight to all food on the board if we're hungry
	FoodHungry int32 = 1 << 14

	// Apply this weight to all food on the board if we're starving
	FoodStarving int32 = 1 << 16

	// EnemyKill means we will eliminate an enemy, but not win the game
	// Increasing this weight will increase our aggression
	EnemyKill int32 = 1 << 12

	// CertainWin means we will collide with the final,
	// shorter opponent's head and they cannot avoid us
	// Use 1 less bit than possible, so that
	// CertainWin + CertainWin doesn't wrap around
	CertainWin int32 = 1 << 30
)

func (s *State) Weight(coord *game.Coord, snake *game.Battlesnake) int32 {
	if SnakeWillDie(s, coord, snake) {
		return CertainDeath
	}

	weight := Neutral

	enemy := s.weightEnemies(coord, snake)
	if enemy <= CertainDeath {
		return CertainDeath
	}
	weight += enemy

	/*
		hazard := s.weightHazard(coord, snake)
		weight += hazard
		if weight < CertainDeath {
			return CertainDeath
		}
	*/

	food := s.weightFood(coord, snake)
	weight += food

	return weight
}

func (s *State) Abort(weight int32) bool {
	return weight <= CertainDeath
}
