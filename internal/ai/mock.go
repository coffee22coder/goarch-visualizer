package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/coffee22coder/goarch-visualizer/internal/analyzer"
)

type MockAnalyzer struct{}

func (a *MockAnalyzer) Analyze(ctx context.Context, g *analyzer.Graph, focus string) (*Analysis, error) {
	nodes := make([]NodeStatus, 0, len(g.Nodes))

	for _, n := range g.Nodes {
		var status Status
		var reason string
		if n.Type == analyzer.NodePackage {
			status = StatusOk
		} else {
			status = StatusWarn
		}

		reason = getReason(n)

		nodes = append(nodes, NodeStatus{
			ID:     n.ID,
			Status: status,
			Reason: reason,
		})
	}

	return &Analysis{
		Nodes: nodes,
		Recommendations: []string{
			"Keep cmd packages thin; move logic to internal/",
		},
	}, nil
}

func (a *MockAnalyzer) ParseQuery(ctx context.Context, g *analyzer.Graph, question string) (analyzer.Query, error) {
	return analyzer.Query{}, fmt.Errorf("mock: ParseQuery not implemented")
}

func getReason(n analyzer.Node) string {
	switch n.Type {
	case analyzer.NodePackage:
		id := string(n.ID)
		if strings.Contains(id, "/cmd/") {
			return "entry point package"
		}
		if strings.Contains(id, "/internal/") {
			return "internal layer package"
		}
		return "package outside cmd/internal layout"
	case analyzer.NodeFunction:
		return "function not reviewed yet"
	case analyzer.NodeTypeDecl:
		return "type declaration"
	}
	return ""
}
