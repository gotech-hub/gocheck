package main

import (
	"fmt"
	"os"

	"github.com/gotech-hub/gocheck/analyzer"
	"github.com/gotech-hub/gocheck/report"
	"github.com/gotech-hub/gocheck/scanner"
)

// Scan quét mã nguồn Go trong path, sinh báo cáo HTML/JSON nếu được chọn.
func Scan(path string, htmlOutput, jsonOutput bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("Invalid path: %s", path)
	}

	files := scanner.ScanDir(path)
	results := analyzer.AnalyzeFiles(files)

	if htmlOutput {
		report.GenerateHTML(results)
		fmt.Println("GoCheck: HTML report generated → report.html")
	}

	if jsonOutput {
		report.GenerateJSON(results)
		fmt.Println("GoCheck: JSON report generated → report.json")
	}

	return nil
}
