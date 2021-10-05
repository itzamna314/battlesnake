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

			nd.Brain.Move(nd.SnakeID, nd.Direction)

			// Remove coord from weight
			// Allow Move to step toward coord
			nd.Weight = nd.Brain.Weight(nd)

			select {
			case <-ctx.Done():
				return
			case expand <- nd:
			}
		}
	}
}
