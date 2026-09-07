package ai

type Status string // "ok" | "warning" | "error"

type NodeStatus struct {
	ID     string `json:"id"`
	Status Status `json:"status"`
	Reason string `json:"reason"`
}
type Analysis struct {
	Nodes           []NodeStatus `json:"nodes"`
	Recommendations []string     `json:"recommendations"`
}
