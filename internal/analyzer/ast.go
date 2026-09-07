package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"strings"
)

func ParseFile(fset *token.FileSet, path string) (*Graph, error) {
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		slog.Error("analyzer parse file", "error", err)
		return nil, err
	}

	g := Graph{
		Nodes: make([]Node, 0),
		Edges: make([]Edge, 0),
	}

	pkg := f.Name.Name   // "testdata"
	file := path         // "testdata/sample.go"
	parent := pkgID(pkg) // "pkg:testdata"

	apecs := Node{
		ID:   parent,
		Name: pkg,
		Type: NodePackage,
		File: file,
		Line: fset.Position(f.Name.Pos()).Line,
	}

	g.Nodes = append(g.Nodes, apecs)

	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		switch x := n.(type) {
		case *ast.ImportSpec:
			edge := newEdge(newEdgeProps{
				kind: EdgeImport,
				from: parent,
				to:   strings.Trim(x.Path.Value, "\""),
			})
			g.Edges = append(g.Edges, edge)
		case *ast.FuncDecl:
			node := newNode(newNodeProps{
				nodeType: NodeFunction,
				pkg:      pkg,
				name:     x.Name.Name,
				file:     file,
				parent:   parent,
				line:     fset.Position(x.Name.Pos()).Line,
			})

			g.Nodes = append(g.Nodes, node)
		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, val := range x.Specs {
					ts, ok := val.(*ast.TypeSpec)
					if !ok {
						continue
					}
					line := fset.Position(ts.Pos()).Line
					node := newNode(newNodeProps{
						nodeType: NodeTypeDecl,
						pkg:      pkg,
						name:     ts.Name.Name,
						file:     file,
						parent:   parent,
						line:     line,
					})
					g.Nodes = append(g.Nodes, node)
				}
			}
		}

		return true
	})

	return &g, nil

}

type newEdgeProps struct {
	kind EdgeType
	from NodeID
	to   string
}

func newEdge(props newEdgeProps) Edge {
	return Edge{
		ID:   edgeID(props.kind, props.from, props.to),
		From: props.from,
		To:   props.to,
		Kind: props.kind,
	}
}

type newNodeProps struct {
	nodeType NodeType
	pkg      string
	name     string
	file     string
	parent   NodeID
	line     int
}

func newNode(props newNodeProps) Node {
	id := nodeID(props.nodeType, props.pkg, props.name)
	return Node{
		ID:       id,
		Name:     props.name,
		Type:     props.nodeType,
		File:     props.file,
		Line:     props.line,
		ParentID: props.parent,
	}
}

func pkgID(pkg string) NodeID {
	return NodeID(fmt.Sprintf("pkg:%s", pkg))
}

func nodeID(kind NodeType, pkg, name string) NodeID {
	return NodeID(fmt.Sprintf("%s:%s:%s", kind, pkg, name))
}

func edgeID(kind EdgeType, from NodeID, to string) EdgeID {
	return EdgeID(fmt.Sprintf("%s|%s|%s", kind, from, to))
}
