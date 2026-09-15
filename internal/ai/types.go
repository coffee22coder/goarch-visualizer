package ai

import "github.com/coffee22coder/goarch-visualizer/internal/analyzer"

type Status string

const (
	StatusOk   Status = "ok"
	StatusWarn Status = "warning"
	StatusErr  Status = "error"
)

type NodeStatus struct {
	ID     analyzer.NodeID `json:"id"`
	Status Status          `json:"status"`
	Reason string          `json:"reason"`
}
type Analysis struct {
	Nodes           []NodeStatus `json:"nodes"`
	Recommendations []string     `json:"recommendations"`
}
