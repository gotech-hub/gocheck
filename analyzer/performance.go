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
					Suggestion: "Check the loop's exit condition or review for nested loops that may impact performance.",
				})

				// Rule: defer in loop
				ast.Inspect(loop.Body, func(n ast.Node) bool {
					if deferStmt, ok := n.(*ast.DeferStmt); ok {
						deferPos := fset.Position(deferStmt.Pos())
						results = append(results, Finding{
							File:       file,
							Line:       deferPos.Line,
							Message:    "Use of 'defer' inside a loop can cause performance issues.",
							Severity:   Medium,
							Suggestion: "Move 'defer' outside the loop if possible, or consider alternative resource management.",
						})
					}
					return true
				})

				// Rule: go statement in loop
				ast.Inspect(loop.Body, func(n ast.Node) bool {
					if goStmt, ok := n.(*ast.GoStmt); ok {
						goPos := fset.Position(goStmt.Pos())
						results = append(results, Finding{
							File:       file,
							Line:       goPos.Line,
							Message:    "Launching goroutines inside a loop can cause race conditions or high overhead.",
							Severity:   Medium,
							Suggestion: "Consider batching data or using a worker pool instead of launching goroutines inside a loop.",
						})
					}
					return true
				})

				// Rule: string += in loop
				ast.Inspect(loop.Body, func(n ast.Node) bool {
					if assign, ok := n.(*ast.AssignStmt); ok {
						if assign.Tok == token.ADD_ASSIGN {
							if rhs, ok := assign.Rhs[0].(*ast.BasicLit); ok && rhs.Kind == token.STRING {
								assignPos := fset.Position(assign.Pos())
								results = append(results, Finding{
									File:       file,
									Line:       assignPos.Line,
									Message:    "String concatenation (+=) inside a loop can be slow.",
									Severity:   Low,
									Suggestion: "Use strings.Builder for string concatenation inside loops.",
								})
							}
						}
					}
					return true
				})
			}
		}
		return true
	})

	return results
}
