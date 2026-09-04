package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "testdata/sample.go", nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	var packageName string
	if f.Name.Name != "" {
		packageName = f.Name.Name
	} else {
		packageName = "package not found"
	}

	fmt.Printf("package: %s\n", packageName)

	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch x := n.(type) {
		case *ast.ImportSpec:
			fmt.Printf("import: %s\n", x.Path.Value)
		case *ast.FuncDecl:
			fmt.Printf("func: %s\n", x.Name.Name)
		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, val := range x.Specs {
					ts, ok := val.(*ast.TypeSpec)
					if !ok {
						continue
					}
					line := fset.Position(ts.Pos()).Line
					l := strconv.Itoa(line)
					fmt.Printf("type: %s; line: %s\n", ts.Name.Name, l)
				}
			}
		}
		return true
	})

}
