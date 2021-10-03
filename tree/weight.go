package tree

import "context"

func weightWorker(ctx context.Context, weight <-chan *Node, expand chan<- *Node) {
	for {
		select {
		case <-ctx.Done():
			return
		case nd, ok := <-weight:
			if !ok {
				return
			}

			nd.Brain.MoveEnemies(nd.SnakeID)

			nd.Weight = nd.Brain.Weight(nd.Coord, nd.SnakeID)

			if nd.Parent != nil {
				nd.Weight += nd.Parent.Weight
			}

			nd.Brain.Move(nd.SnakeID, nd.Direction)

			select {
			case <-ctx.Done():
				return
			case expand <- nd:
			}
		}
	}
}
