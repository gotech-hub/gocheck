package analyzer

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
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
							Suggestion: "Do not hardcode passwords/API keys. Use environment variables or configuration files instead.",
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
						Suggestion: "Avoid passing unchecked input to exec.Command. Use input validation and sanitization.",
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
									Suggestion: "Use HTTPS (443) instead of HTTP (80/8080) for production services.",
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
							Suggestion: "Do not use md5 or sha1 for security purposes. Use sha256 or stronger algorithms instead.",
						})
					}
				}
			}
		}

		// Rule 5: tls.Config{InsecureSkipVerify: true}
		if cl, ok := n.(*ast.CompositeLit); ok {
			if se, ok := cl.Type.(*ast.SelectorExpr); ok {
				if pkg, ok := se.X.(*ast.Ident); ok && pkg.Name == "tls" && se.Sel.Name == "Config" {
					for _, elt := range cl.Elts {
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.Ident); ok && key.Name == "InsecureSkipVerify" {
								if val, ok := kv.Value.(*ast.Ident); ok && val.Name == "true" {
									pos := fset.Position(kv.Pos())
									results = append(results, Finding{
										File:       file,
										Line:       pos.Line,
										Message:    "tls.Config with InsecureSkipVerify: true detected (insecure TLS)",
										Severity:   Critical,
										Suggestion: "Never set InsecureSkipVerify to true in production. This disables certificate validation and is highly insecure.",
									})
								}
							}
						}
					}
				}
			}
		}

		return true
	})

	// --- Gosec integration ---
	gosecFindings := runGosec(file)
	results = append(results, gosecFindings...)

	return results
}

// runGosec executes gosec on the given file and parses the JSON output into []Finding
func runGosec(file string) []Finding {
	cmd := exec.Command("gosec", "-fmt=json", file)
	output, err := cmd.Output()
	if err != nil {
		return nil // Nếu không có gosec hoặc lỗi, bỏ qua
	}
	type gosecResult struct {
		Issues []struct {
			Severity   string `json:"severity"`
			Confidence string `json:"confidence"`
			Cwe        struct {
				ID  string `json:"ID"`
				URL string `json:"URL"`
			} `json:"cwe"`
			RuleID  string `json:"rule_id"`
			Details string `json:"details"`
			File    string `json:"file"`
			Code    string `json:"code"`
			Line    int    `json:"line"`
		} `json:"issues"`
	}
	var res gosecResult
	err = json.Unmarshal(output, &res)
	if err != nil {
		return nil
	}
	var findings []Finding
	for _, issue := range res.Issues {
		sev := Medium
		switch strings.ToLower(issue.Severity) {
		case "low":
			sev = Low
		case "medium":
			sev = Medium
		case "high":
			sev = High
		case "critical":
			sev = Critical
		}
		findings = append(findings, Finding{
			File:       issue.File,
			Line:       issue.Line,
			Message:    fmt.Sprintf("[gosec][%s] %s", issue.RuleID, issue.Details),
			Severity:   sev,
			Suggestion: issue.Cwe.URL,
		})
	}
	return findings
}
