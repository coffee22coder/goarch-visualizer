package analyzer

import (
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func AnalyzeDir(root string) (*Graph, error) {
	fset := token.NewFileSet()
	result := &Graph{Nodes: make([]Node, 0), Edges: make([]Edge, 0)}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if d.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		g, err := ParseFile(fset, path)
		if err != nil {
			return err
		}

		mergeGraph(result, g)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func mergeGraph(dst *Graph, src *Graph) {
	nodeSet := make(map[NodeID]struct{}, len(dst.Nodes))

	for _, n := range dst.Nodes {
		nodeSet[n.ID] = struct{}{}
	}

	for _, n := range src.Nodes {
		if _, ok := nodeSet[n.ID]; !ok {
			dst.Nodes = append(dst.Nodes, n)
			nodeSet[n.ID] = struct{}{}
		}
	}

	edgeSet := make(map[EdgeID]struct{}, len(dst.Edges))

	for _, e := range dst.Edges {
		edgeSet[e.ID] = struct{}{}
	}

	for _, e := range src.Edges {
		if _, ok := edgeSet[e.ID]; !ok {
			dst.Edges = append(dst.Edges, e)
			edgeSet[e.ID] = struct{}{}
		}
	}

}
