package renderer

import (
	"strings"

	"github.com/coffee22coder/goarch-visualizer/internal/ai"
	"github.com/coffee22coder/goarch-visualizer/internal/analyzer"
)

func ToMermaid(g analyzer.Graph, a ai.Analysis) string {
	b := strings.Builder{}
	b.WriteString("graph LR\n")

	b.WriteString("classDef ok fill:#d4edda,stroke:#28a745,color:#000\n")
	b.WriteString("classDef warning fill:#fff3cd,stroke:#ffc107,color:#000\n")
	b.WriteString("classDef error fill:#f8d7da,stroke:#dc3545,color:#000\n")

	packageMap := make(map[analyzer.NodeID]analyzer.Node, len(g.Nodes))
	statusByID := make(map[analyzer.NodeID]ai.Status, len(a.Nodes))

	for _, n := range a.Nodes {
		statusByID[n.ID] = n.Status
	}

	for _, n := range g.Nodes {
		if n.Type == analyzer.NodePackage {
			packageMap[n.ID] = n

			b.WriteString(safeId(string(n.ID)) + "[" + n.Name + "]\n")
		}
	}

	for id, n := range packageMap {
		st := ai.StatusOk
		if s, ok := statusByID[id]; ok {
			st = s
		}
		b.WriteString("class " + safeId(string(n.ID)) + " " + string(st) + "\n")
	}

	for _, e := range g.Edges {
		if e.Kind != analyzer.EdgeImport {
			continue
		}
		if _, ok := packageMap[e.From]; !ok {
			continue
		}
		if _, ok := packageMap[analyzer.NodeID(e.To)]; !ok {
			continue
		}
		b.WriteString(safeId(string(e.From)) + "-->" + safeId(e.To) + "\n")
	}
	return b.String()
}

// замена "/" и ":" на "_"
func safeId(id string) string {
	id = strings.ReplaceAll(id, "/", "_")
	id = strings.ReplaceAll(id, ":", "_")
	id = strings.ReplaceAll(id, ".", "_")
	return id
}
