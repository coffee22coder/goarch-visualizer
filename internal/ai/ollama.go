package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coffee22coder/goarch-visualizer/internal/analyzer"
)

type ollamaResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

type OllamaAnalyzer struct{}

func (a *OllamaAnalyzer) Analyze(ctx context.Context, g *analyzer.Graph, focus string) (*Analysis, error) {

	pkgGraph := analyzer.PackageGraph(g)

	graphJSON, err := json.Marshal(pkgGraph)
	if err != nil {
		return nil, err
	}

	body := map[string]any{
		"model":   "qwen2.5-coder:7b",
		"stream":  false,
		"format":  "json",
		"options": map[string]int{"temperature": 0},
		"messages": []map[string]any{
			map[string]any{
				"role":    "system",
				"content": "You analyze a Go architecture graph (packages only). Reply with ONLY JSON matching this schema: {\"nodes\":[{\"id\":\"string\",\"status\":\"ok|warning|error\",\"reason\":\"string\"}],\"recommendations\":[\"string\"]}. Use only node ids from the input. Do not invent ids Return status for every package node id from the input..",
			},
			map[string]any{
				"role":    "user",
				"content": "Graph:\n" + string(graphJSON) + "\nFocus: (all)",
			},
		},
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:11434/api/chat", bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %d", response.StatusCode)
	}

	var resp ollamaResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	var analysis Analysis
	err = json.Unmarshal([]byte(resp.Message.Content), &analysis)
	if err != nil {
		return nil, err
	}

	return &analysis, nil
}
