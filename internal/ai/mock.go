package ai

import (
	"context"

	"github.com/coffee22coder/goarch-visualizer/internal/analyzer"
)

type MockAnalyzer struct{}

func (a *MockAnalyzer) Analyze(ctx context.Context, g *analyzer.Graph, focus string) (*Analysis, error) {
	return &Analysis{
		Nodes: []NodeStatus{
			{ID: "1", Status: "ok", Reason: "Node 1 is ok"},
		},
	}, nil
}
