package ai

import (
	"context"

	"github.com/coffee22coder/goarch-visualizer/internal/analyzer"
)

type Analyzer interface {
	Analyze(ctx context.Context, g *analyzer.Graph, focus string) (*Analysis, error)
	ParseQuery(ctx context.Context, g *analyzer.Graph, question string) (analyzer.Query, error)
}
