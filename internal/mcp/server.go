package mcp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	s := server.NewMCPServer(
		"Demo Go Arch Visualizer ",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogger(logger),
	)

	analyzeProjectTool := mcp.NewTool("analyze_project",
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The path to project"),
		),
		mcp.WithNumber("level",
			mcp.Description("The level analyze"),
			mcp.DefaultNumber(2),
		))

	s.AddTool(analyzeProjectTool, analyzeHandler)

	if err := server.ServeStdio(s); err != nil {
		slog.Error("MCP server error: %v", "error", err)
		return err
	}

	return nil

}

func analyzeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	p, err := request.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	l, err := request.RequireInt("level")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("analyzeHandler: path = %s, level = %d", p, l)), nil
}
