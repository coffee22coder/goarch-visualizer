package ai

import "github.com/coffee22coder/goarch-visualizer/internal/analyzer"

func Validate(g *analyzer.Graph, a *Analysis) *Analysis {
	allowed := map[analyzer.NodeID]struct{}{}

	for _, n := range g.Nodes {
		allowed[n.ID] = struct{}{}
	}

	checked := make([]NodeStatus, 0, len(a.Nodes))

	for _, ns := range a.Nodes {
		if _, ok := allowed[ns.ID]; !ok {
			continue
		}
		checked = append(checked, ns)
	}

	return &Analysis{
		Nodes:           checked,
		Recommendations: a.Recommendations,
	}
}
