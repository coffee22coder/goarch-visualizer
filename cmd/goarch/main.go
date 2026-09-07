package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/coffee22coder/goarch-visualizer/internal/ai"
	"github.com/coffee22coder/goarch-visualizer/internal/analyzer"
	"github.com/coffee22coder/goarch-visualizer/internal/mcp"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	if len(os.Args) >= 2 && os.Args[1] == "analyze" {
		if len(os.Args) < 3 || os.Args[2] == "" {
			slog.Error("CLI mode: path required")
			os.Exit(1)
		}
		runCLI(os.Args[2])
		return
	}

	err := runMCP(context.Background(), logger)
	if err != nil {
		slog.Error("MCP mode", "error", err)
		os.Exit(1)
	}
}

func runCLI(path string) {
	g, err := analyzer.AnalyzeDir(path)
	if err != nil {
		os.Exit(1)
	}
	slog.Info("analyze done", "nodes", len(g.Nodes), "edges", len(g.Edges))

	ollama := &ai.OllamaAnalyzer{}

	analysis, err := ollama.Analyze(context.Background(), g, "all")
	if err != nil {
		slog.Error("analyze failed", "error", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		slog.Error("marshal analysis", "error", err)
		os.Exit(1)
	}
	fmt.Println(string(out)) // stdout — удобнее смотреть, чем slog
}

func runMCP(ctx context.Context, logger *slog.Logger) error {

	err := mcp.Run(ctx, logger)
	if err != nil {
		return err
	}
	return nil
}
