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
		// Rule 1: Hardcoded credentials
		if assign, ok := n.(*ast.AssignStmt); ok {
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

		// Rule 2: Use of exec.Command (possible command injection)
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := fun.X.(*ast.Ident); ok && pkg.Name == "exec" && fun.Sel.Name == "Command" {
					pos := fset.Position(call.Pos())
					results = append(results, Finding{
						File:       file,
						Line:       pos.Line,
						Message:    "Use of exec.Command detected (possible command injection)",
						Severity:   High,
						Suggestion: "Tránh truyền tham số không kiểm soát vào exec.Command. Sử dụng các biện pháp kiểm tra đầu vào.",
					})
				}
			}
		}

		// Rule 3: Use of http.ListenAndServe with :80 or :8080 (insecure port)
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := fun.X.(*ast.Ident); ok && pkg.Name == "http" && fun.Sel.Name == "ListenAndServe" {
					if len(call.Args) > 0 {
						if bl, ok := call.Args[0].(*ast.BasicLit); ok && bl.Kind.String() == "STRING" {
							v := strings.Trim(bl.Value, "\"")
							if v == ":80" || v == ":8080" {
								pos := fset.Position(call.Pos())
								results = append(results, Finding{
									File:       file,
									Line:       pos.Line,
									Message:    "Use of http.ListenAndServe on insecure port (:80 or :8080)",
									Severity:   Medium,
									Suggestion: "Nên sử dụng HTTPS (443) thay vì HTTP (80/8080) cho các dịch vụ production.",
								})
							}
						}
					}
				}
			}
		}

		// Rule 4: Use of md5.New or sha1.New (insecure hash)
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := fun.X.(*ast.Ident); ok {
					if (pkg.Name == "md5" && fun.Sel.Name == "New") || (pkg.Name == "sha1" && fun.Sel.Name == "New") {
						pos := fset.Position(call.Pos())
						results = append(results, Finding{
							File:       file,
							Line:       pos.Line,
							Message:    fmt.Sprintf("Use of insecure hash function: %s.New", pkg.Name),
							Severity:   High,
							Suggestion: "Không sử dụng md5 hoặc sha1 cho mục đích bảo mật. Hãy dùng sha256 hoặc các thuật toán mạnh hơn.",
						})
					}
				}
			}
		}

		return true
	})

	return results
}
