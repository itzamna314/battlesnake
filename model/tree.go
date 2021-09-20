package model

type GameTree struct {
	Root *TreeNode
}

type TreeNode struct {
	State    *GameState
	Children [4]*TreeNode
	Weight   float64
}
