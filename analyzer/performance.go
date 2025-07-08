package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func analyzePerformance(file string) []Finding {
	var results []Finding
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return results
	}

	ast.Inspect(node, func(n ast.Node) bool {
		loop, ok := n.(*ast.ForStmt)
		if ok {
			pos := fset.Position(loop.Pos())
			if !strings.Contains(file, "_test.go") {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    "For-loop detected — review for potential performance impact",
					Severity:   Low,
					Suggestion: "Kiểm tra lại điều kiện dừng hoặc xử lý nhiều vòng lặp lồng nhau.",
				})
			}
		}
		return true
	})

	return results
}
