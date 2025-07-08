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
		if ok && fn.Body != nil {
			pos := fset.Position(fn.Pos())

			// Rule 1: Hàm quá dài
			if len(fn.Body.List) > 20 {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s is too long (%d lines)", fn.Name.Name, len(fn.Body.List)),
					Severity:   Medium,
					Suggestion: "Tách hàm ra thành nhiều hàm nhỏ để dễ đọc và test.",
				})
			}

			// Rule 2: Hàm có quá nhiều tham số
			if fn.Type.Params != nil && len(fn.Type.Params.List) > 4 {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s has too many parameters (%d)", fn.Name.Name, len(fn.Type.Params.List)),
					Severity:   Medium,
					Suggestion: "Cân nhắc gom nhóm các tham số hoặc sử dụng struct.",
				})
			}

			// Rule 3: Hàm lồng nhau quá sâu (>2 cấp)
			var maxDepth, currDepth int
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				switch n.(type) {
				case *ast.BlockStmt:
					currDepth++
					if currDepth > maxDepth {
						maxDepth = currDepth
					}
					return true
				default:
					return true
				}
			})
			if maxDepth > 3 { // 1 là thân hàm, >3 là lồng nhau >2 cấp
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s is nested too deeply (%d levels)", fn.Name.Name, maxDepth),
					Severity:   Medium,
					Suggestion: "Giảm độ lồng nhau, tách hàm/phân rã logic.",
				})
			}

			// Rule 4: Hàm có nhiều return statement (>2)
			returnCount := 0
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if _, ok := n.(*ast.ReturnStmt); ok {
					returnCount++
				}
				return true
			})
			if returnCount > 2 {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s has too many return statements (%d)", fn.Name.Name, returnCount),
					Severity:   Low,
					Suggestion: "Cân nhắc đơn giản hóa luồng trả về.",
				})
			}

			// Rule 5: Hàm có nhiều nhánh if/else (>3)
			ifCount := 0
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				switch n.(type) {
				case *ast.IfStmt:
					ifCount++
				}
				return true
			})
			if ifCount > 3 {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s has too many if/else branches (%d)", fn.Name.Name, ifCount),
					Severity:   Low,
					Suggestion: "Cân nhắc refactor lại logic điều kiện.",
				})
			}
		}
		return true
	})

	return results
}
