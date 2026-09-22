package analyzer

import (
	"fmt"
)

type Op string

const (
	OpDependents   Op = "dependents"   // кто импортирует Target
	OpDependencies Op = "dependencies" // что импортирует Target
	OpPath         Op = "path"         // путь From → To
	OpNeighbors    Op = "neighbors"    // Target + соседи на 1 hop
)

type Query struct {
	Op            Op       `json:"op"`
	Target        NodeID   `json:"target"`
	From          NodeID   `json:"from"`
	To            NodeID   `json:"to"`
	ExcludePrefix []string `json:"exclude_prefix"`
}

func Execute(g *Graph, q Query) (*Graph, error) {
	pg := PackageGraph(g)

	switch q.Op {
	case OpDependents, OpDependencies:
		err := includesTarget(pg, q.Target)
		if err != nil {
			return nil, err
		}

		g := createGraph(q.Op, pg, q.Target)
		return g, nil

	case OpPath:
		queue := []NodeID{q.From}
		seen := map[NodeID]struct{}{q.From: {}}
		prev := map[NodeID]NodeID{}

	LOOP:
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			for _, e := range pg.Edges {
				if e.From != current {
					continue
				}

				next := NodeID(e.To)
				if _, ok := seen[next]; ok {
					continue
				}
				seen[next] = struct{}{}
				prev[next] = current
				if next == q.To {
					break LOOP
				}
				queue = append(queue, next)
			}

		}

		if _, ok := seen[q.To]; !ok {
			return nil, fmt.Errorf("Empty graph. No import FROM: %s, TO: %s", q.From, q.To)
		}

		newNodes := make([]Node, 0, len(g.Nodes))
		newEdges := make([]Edge, 0, len(g.Edges))

		nodeSet := map[NodeID]struct{}{q.To: {}}
		cur := q.To

		for cur != q.From {
			cur = prev[cur]
			nodeSet[cur] = struct{}{}
		}

		for _, n := range pg.Nodes {
			if _, ok := nodeSet[n.ID]; ok {
				newNodes = append(newNodes, n)
			}
		}

		for _, e := range pg.Edges {
			_, okFrom := nodeSet[e.From]
			_, okTo := nodeSet[NodeID(e.To)]
			if okFrom && okTo {
				newEdges = append(newEdges, e)
			}
		}

		return &Graph{
			Nodes: newNodes,
			Edges: newEdges,
		}, nil

	case OpNeighbors:
		err := includesTarget(pg, q.Target)
		if err != nil {
			return nil, err
		}
		g := createGraph(q.Op, pg, q.Target)
		return g, nil
	}

	return nil, fmt.Errorf("Operation %s is not valid", q.Op)
}

func includesTarget(g *Graph, t NodeID) error {
	for _, n := range g.Nodes {
		if n.ID == t {
			return nil
		}
	}
	return fmt.Errorf("target %s not found", t)
}

func createGraph(op Op, g *Graph, t NodeID) *Graph {
	nodeSet := make(map[NodeID]struct{})
	nodeSet[t] = struct{}{}

	for _, e := range g.Edges {
		switch op {
		case OpDependents:
			if e.To == string(t) {
				nodeSet[e.From] = struct{}{}
			}
		case OpDependencies:
			if e.From == t {
				nodeSet[NodeID(e.To)] = struct{}{}
			}
		case OpNeighbors:
			{
				if e.To == string(t) {
					nodeSet[e.From] = struct{}{}
				}
				if e.From == t {
					nodeSet[NodeID(e.To)] = struct{}{}
				}
			}
		}

	}

	newNodes := make([]Node, 0, len(g.Nodes))
	newEdges := make([]Edge, 0, len(g.Edges))

	for _, n := range g.Nodes {
		if _, ok := nodeSet[n.ID]; ok {
			newNodes = append(newNodes, n)
		}
	}

	for _, e := range g.Edges {
		_, okFrom := nodeSet[e.From]
		_, okTo := nodeSet[NodeID(e.To)]
		if okFrom && okTo {
			newEdges = append(newEdges, e)
		}
	}

	return &Graph{
		Nodes: newNodes,
		Edges: newEdges,
	}
}
