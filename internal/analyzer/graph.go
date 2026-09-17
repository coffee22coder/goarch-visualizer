package analyzer

type NodeType string
type EdgeType string

type NodeID string
type EdgeID string

const (
	NodePackage  NodeType = "package"
	NodeFunction NodeType = "function"
	NodeTypeDecl NodeType = "type" // struct + interface
)

const EdgeImport EdgeType = "import"
const EdgeCall EdgeType = "call"
const EdgeImplement EdgeType = "implement"

type Node struct {
	ID   NodeID
	Name string
	Type NodeType
	File string
	Line int

	ParentID    NodeID
	ChildrenIDs []NodeID
}

type Edge struct {
	ID   EdgeID
	From NodeID
	To   string
	Kind EdgeType
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}

func PackageGraph(g *Graph) *Graph {
	packageIDs := make(map[NodeID]struct{})
	for _, n := range g.Nodes {
		if n.Type == NodePackage {
			packageIDs[n.ID] = struct{}{}
		}
	}
	nodes := make([]Node, 0, len(packageIDs))
	for _, n := range g.Nodes {
		if n.Type == NodePackage {
			nodes = append(nodes, n)
		}
	}
	edges := make([]Edge, 0)
	for _, e := range g.Edges {
		if e.Kind != EdgeImport {
			continue
		}
		if _, ok := packageIDs[e.From]; !ok {
			continue
		}
		if _, ok := packageIDs[NodeID(e.To)]; !ok {
			continue
		}
		edges = append(edges, e)
	}
	return &Graph{Nodes: nodes, Edges: edges}
}
