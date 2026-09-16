package mcp

import (
	"context"
	"log/slog"

	"github.com/coffee22coder/goarch-visualizer/internal/ai"
	"github.com/coffee22coder/goarch-visualizer/internal/analyzer"
	"github.com/coffee22coder/goarch-visualizer/internal/renderer"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Run(ctx context.Context, a ai.Analyzer, logger *slog.Logger) error {
	s := server.NewMCPServer(
		"Demo Go Arch Visualizer ",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogger(logger),
	)

	analyzeProjectTool := mcp.NewTool("analyze_project",
		mcp.WithDescription(
			"Analyze Go project architecture: builds package dependency graph from AST "+
				"and returns a Mermaid diagram. Use this instead of manually reading imports "+
				"when the user asks about Go project structure, layers, or dependencies.",
		),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The path to project"),
		),
		mcp.WithNumber("level",
			mcp.Description("The level analyze"),
			mcp.DefaultNumber(2),
		))

	s.AddTool(analyzeProjectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return analyzeHandler(ctx, request, a)
	})

	if err := server.ServeStdio(s); err != nil {
		slog.Error("MCP server error: %v", "error", err)
		return err
	}

	return nil

}

func analyzeHandler(ctx context.Context, request mcp.CallToolRequest, a ai.Analyzer) (*mcp.CallToolResult, error) {
	p, err := request.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	// level := request.GetInt("level", 2)

	graph, err := analyzer.AnalyzeDir(p)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	analysis, err := a.Analyze(context.Background(), graph, "all")
	if err != nil {
		slog.Error("analyze failed", "error", err)
		return nil, err
	}

	valid := ai.Validate(graph, analysis)

	mermaid := renderer.ToMermaid(*graph, *valid)

	return mcp.NewToolResultText(mermaid), nil
}
