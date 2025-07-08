package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/gotech-hub/gocheck/utils"
)

const (
	maxFuncLines      = 100
	maxFuncParams     = 4
	maxNestingDepth   = 3 // 1 is function body, >3 means nested more than 2 levels
	maxReturnStmts    = 2
	maxIfElseBranches = 3
	maxLocalVars      = 8
	minFuncNameLength = 3
)

func AnalyzeCleanCode(file string) []Finding {
	var results []Finding
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, parser.ParseComments, 0)
	if err != nil {
		return results
	}

	// Rule 9: Avoid global variables in the file
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						pos := fset.Position(name.Pos())
						results = append(results, Finding{
							File:       file,
							Line:       pos.Line,
							Message:    fmt.Sprintf("Global variable '%s' should be avoided", name.Name),
							Severity:   Medium,
							Suggestion: "Avoid using global variables. Use function parameters or struct fields instead.",
							Category:   "Clean",
						})
					}
				}
			}
		}
	}

	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Body != nil {
			pos := fset.Position(fn.Pos())

			// Rule 1: Function is too long
			if len(fn.Body.List) > maxFuncLines {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s is too long (%d lines)", fn.Name.Name, len(fn.Body.List)),
					Severity:   Medium,
					Suggestion: "Split the function into smaller functions for better readability and testability.",
					Category:   "Clean",
				})
			}

			// Rule 2: Function has too many parameters
			if fn.Type.Params != nil && len(fn.Type.Params.List) > maxFuncParams {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s has too many parameters (%d)", fn.Name.Name, len(fn.Type.Params.List)),
					Severity:   Medium,
					Suggestion: "Consider grouping parameters or using a struct.",
					Category:   "Clean",
				})
			}

			// Rule 3: Function is nested too deeply (>2 levels)
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
			if maxDepth > maxNestingDepth {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s is nested too deeply (%d levels)", fn.Name.Name, maxDepth),
					Severity:   Medium,
					Suggestion: "Reduce nesting, split logic into smaller functions.",
					Category:   "Clean",
				})
			}

			// Rule 4: Function has too many return statements (>2)
			returnCount := 0
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if _, ok := n.(*ast.ReturnStmt); ok {
					returnCount++
				}
				return true
			})
			if returnCount > maxReturnStmts {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s has too many return statements (%d)", fn.Name.Name, returnCount),
					Severity:   Low,
					Suggestion: "Consider simplifying the return flow.",
					Category:   "Clean",
				})
			}

			// Rule 5: Function has too many if/else branches (>3)
			ifCount := 0
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				switch n.(type) {
				case *ast.IfStmt:
					ifCount++
				}
				return true
			})
			if ifCount > maxIfElseBranches {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s has too many if/else branches (%d)", fn.Name.Name, ifCount),
					Severity:   Low,
					Suggestion: "Consider refactoring the conditional logic.",
					Category:   "Clean",
				})
			}

			// Rule 6: Function has too many local variables
			localVarCount := 0
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if decl, ok := n.(*ast.AssignStmt); ok {
					for _, lhs := range decl.Lhs {
						if ident, ok := lhs.(*ast.Ident); ok && ident.Obj != nil && ident.Obj.Kind == ast.Var {
							localVarCount++
						}
					}
				}
				if decl, ok := n.(*ast.DeclStmt); ok {
					if genDecl, ok := decl.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
						for _, spec := range genDecl.Specs {
							if valueSpec, ok := spec.(*ast.ValueSpec); ok {
								localVarCount += len(valueSpec.Names)
							}
						}
					}
				}
				return true
			})
			if localVarCount > maxLocalVars {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s has too many local variables (%d)", fn.Name.Name, localVarCount),
					Severity:   Low,
					Suggestion: "Reduce the number of local variables or split logic into smaller functions.",
					Category:   "Clean",
				})
			}

			// Rule 7: Function name is too short
			if len(fn.Name.Name) < minFuncNameLength {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function name '%s' is too short", fn.Name.Name),
					Severity:   Low,
					Suggestion: "Use a more descriptive function name.",
					Category:   "Clean",
				})
			}

			// Rule 8: Unused local variables
			usedVars := make(map[string]bool)
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if ident, ok := n.(*ast.Ident); ok && ident.Obj != nil && ident.Obj.Kind == ast.Var {
					usedVars[ident.Name] = true
				}
				return true
			})
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if decl, ok := n.(*ast.AssignStmt); ok {
					for _, lhs := range decl.Lhs {
						if ident, ok := lhs.(*ast.Ident); ok && ident.Obj != nil && ident.Obj.Kind == ast.Var {
							if !usedVars[ident.Name] {
								results = append(results, Finding{
									File:       file,
									Line:       pos.Line,
									Message:    fmt.Sprintf("Local variable '%s' declared but not used", ident.Name),
									Severity:   Low,
									Suggestion: "Remove unused local variable.",
									Category:   "Clean",
								})
							}
						}
					}
				}
				return true
			})

			// Rule 10: Avoid nested function declarations
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if innerFn, ok := n.(*ast.FuncDecl); ok && innerFn != fn {
					innerPos := fset.Position(innerFn.Pos())
					results = append(results, Finding{
						File:       file,
						Line:       innerPos.Line,
						Message:    fmt.Sprintf("Nested function '%s' should be avoided", innerFn.Name.Name),
						Severity:   Medium,
						Suggestion: "Declare functions at the top level, not inside other functions.",
						Category:   "Clean",
					})
				}
				return true
			})

			// Rule 11: Avoid magic numbers (except 0, 1, -1)
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				switch lit := n.(type) {
				case *ast.BasicLit:
					if lit.Kind == token.INT {
						if lit.Value != "0" && lit.Value != "1" && lit.Value != "-1" {
							magicPos := fset.Position(lit.Pos())
							results = append(results, Finding{
								File:       file,
								Line:       magicPos.Line,
								Message:    fmt.Sprintf("Magic number %s detected", lit.Value),
								Severity:   Low,
								Suggestion: "Replace magic numbers with named constants.",
								Category:   "Clean",
							})
						}
					}
				}
				return true
			})

			// Rule 12: Avoid too many comments in a function (>5)
			commentCount := 0
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if _, ok := n.(*ast.Comment); ok {
					commentCount++
				}
				return true
			})
			if commentCount > 5 {
				results = append(results, Finding{
					File:       file,
					Line:       pos.Line,
					Message:    fmt.Sprintf("Function %s has too many comments (%d)", fn.Name.Name, commentCount),
					Severity:   Low,
					Suggestion: "Refactor code to be self-explanatory and reduce excessive comments.",
					Category:   "Clean",
				})
			}

			// Rule 13: Avoid commented-out code (code that is commented out)
			if fn.Body != nil && node.Comments != nil {
				for _, cg := range node.Comments {
					for _, c := range cg.List {
						if utils.IsCommentedOutCode(c.Text) {
							commentPos := fset.Position(c.Pos())
							results = append(results, Finding{
								File:       file,
								Line:       commentPos.Line,
								Message:    "Commented-out code detected",
								Severity:   Low,
								Suggestion: "Remove commented-out code for better readability.",
								Category:   "Clean",
							})
						}
					}
				}
			}
		}
		return true
	})

	return results
}
