package tree

type ConfigFn func(*Tree)

func ConfigMaxDepth(depth int) func(*Tree) {
	return func(t *Tree) {
		t.MaxDepth = depth
	}
}
