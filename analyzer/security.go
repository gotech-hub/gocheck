package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func analyzeSecurity(file string) []Finding {
	var results []Finding
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		return results
	}

	ast.Inspect(node, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if ok {
			for _, expr := range assign.Rhs {
				bl, ok := expr.(*ast.BasicLit)
				if ok && bl.Kind.String() == "STRING" {
					v := strings.ToLower(bl.Value)
					if strings.Contains(v, "key") || strings.Contains(v, "password") {
						pos := fset.Position(bl.Pos())
						results = append(results, Finding{
							File:       file,
							Line:       pos.Line,
							Message:    fmt.Sprintf("Hardcoded credential: %s", bl.Value),
							Severity:   High,
							Suggestion: "Không hardcode mật khẩu/API key. Dùng biến môi trường hoặc config file.",
						})
					}
				}
			}
		}
		return true
	})

	return results
}
