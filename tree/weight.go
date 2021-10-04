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

			// TODO: Combine move enemies and move snake
			nd.Brain.MoveEnemies(nd.SnakeID)

			// Remove coord from weight
			// Allow Move to step toward coord
			nd.Weight = nd.Brain.Weight(nd)

			// TODO: Move to multiverse brain
			/*
				if nd.Parent != nil {
					nd.Weight += nd.Parent.Weight
				}
			*/

			nd.Brain.Move(nd.SnakeID, nd.Direction)

			select {
			case <-ctx.Done():
				return
			case expand <- nd:
			}
		}
	}
}
