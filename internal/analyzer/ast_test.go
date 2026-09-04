package analyzer_test

import (
	"go/token"
	"testing"

	"github.com/coffee22coder/goarch-visualizer/internal/analyzer"
)

func TestParseFile(t *testing.T) {
	fset := token.NewFileSet()
	g, err := analyzer.ParseFile(fset, "../../testdata/sample.go")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(g.Nodes) != 6 {
		t.Fatalf("expected 6 nodes, got %d", len(g.Nodes))
	}
	if len(g.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(g.Edges))
	}
}
