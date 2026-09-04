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

	ParentID    string
	ChildrenIDs []string
}

type Edge struct {
	ID   EdgeID
	From string
	To   string
	Kind EdgeType
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}
