package model

type GameTree struct {
	Root *TreeNode
}

type TreeNode struct {
	State    *GameState
	Moves    PossibleMoves
	Children [4]*TreeNode
}
