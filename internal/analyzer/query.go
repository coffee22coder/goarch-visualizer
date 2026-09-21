package analyzer

import "fmt"

type Op string

const (
	OpDependents   Op = "dependents"   // кто импортирует Target
	OpDependencies Op = "dependencies" // что импортирует Target
	OpPath         Op = "path"         // путь From → To
	OpNeighbors    Op = "neighbors"    // Target + соседи на 1 hop
)

type Query struct {
	Op            Op
	Target        NodeID
	From, To      NodeID
	ExcludePrefix []string
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
	case OpNeighbors:
	default:
		return nil, fmt.Errorf("Operation %s is not valid", q.Op)
	}

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
