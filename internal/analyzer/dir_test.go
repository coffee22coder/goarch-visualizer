package analyzer_test

import (
	"testing"

	"github.com/coffee22coder/goarch-visualizer/internal/analyzer"
)

func TestAnalyzeDir(t *testing.T) {
	g, err := analyzer.AnalyzeDir("../../testdata/sample-project")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nodes:", len(g.Nodes), "edges:", len(g.Edges))

	if len(g.Nodes) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(g.Nodes))
	}
	if len(g.Edges) != 2 {
		t.Fatalf("expected 2 edge, got %d", len(g.Edges))
	}
}
