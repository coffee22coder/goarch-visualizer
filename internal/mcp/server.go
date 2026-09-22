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
		mcp.WithString("question",
			mcp.Description("Optional. Who imports X, path from A to B, ..."),
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
	question := request.GetString("question", "")

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	graph, err := analyzer.AnalyzeDir(p)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if question == "" {
		mermaid := renderer.ToMermaid(*graph, ai.Analysis{})
		return mcp.NewToolResultText(mermaid), nil
	} else {
		query, err := a.ParseQuery(context.Background(), graph, question)
		if err != nil {
			slog.Error("analyze failed", "error", err)
			return nil, err
		}

		subgraph, err := analyzer.Execute(graph, query)
		if err != nil {
			slog.Error("analyze failed", "error", err)
			return nil, err
		}

		mermaid := renderer.ToMermaid(*subgraph, ai.Analysis{})

		return mcp.NewToolResultText(mermaid), nil
	}

}
