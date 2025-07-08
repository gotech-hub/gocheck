package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func AnalyzeCleanCode(file string) []Finding {
	var results []Finding
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return results
	}

	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Body != nil && len(fn.Body.List) > 20 {
			pos := fset.Position(fn.Pos())
			results = append(results, Finding{
				File:       file,
				Line:       pos.Line,
				Message:    fmt.Sprintf("Function %s is too long (%d lines)", fn.Name.Name, len(fn.Body.List)),
				Severity:   Medium,
				Suggestion: "Tách hàm ra thành nhiều hàm nhỏ để dễ đọc và test.",
			})
		}
		return true
	})

	return results
}
