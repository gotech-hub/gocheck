package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/gotech-hub/gocheck/analyzer"
	"github.com/gotech-hub/gocheck/report"
	"github.com/gotech-hub/gocheck/scanner"
)

const version = "gocheck v1.0.1"

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

func showHelp() {
	fmt.Println("gocheck - A tool for scanning Go code for issues")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  gocheck [flags]")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --path string     Path to scan (default: .)")
	fmt.Println("  --html            Generate HTML report (default: true)")
	fmt.Println("  --json            Generate JSON report (default: true)")
	fmt.Println("  --version         Show version information")
	fmt.Println("  --help            Show this help message")
	fmt.Println("  --verbose         Enable verbose output")
	fmt.Println("  --stats           Show statistics after scanning")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  gocheck --path ./myproject --html --json")
	fmt.Println("  gocheck --path ./src --json=false --html=true --verbose")
}

func main() {
	// Check for required tools
	missingTools := []string{}
	if _, err := exec.LookPath("gosec"); err != nil {
		missingTools = append(missingTools, "gosec")
	}
	if _, err := exec.LookPath("staticcheck"); err != nil {
		missingTools = append(missingTools, "staticcheck")
	}
	if len(missingTools) > 0 {
		fmt.Printf("\u274c Error: The following required tools are missing: %v\n", missingTools)
		fmt.Println("Please install them as described in the README before running gocheck.")
		os.Exit(1)
	}

	var (
		path    = flag.String("path", ".", "Path to scan")
		html    = flag.Bool("html", true, "Generate HTML report")
		json    = flag.Bool("json", true, "Generate JSON report")
		showVer = flag.Bool("version", false, "Show version information")
		help    = flag.Bool("help", false, "Show help information")
		verbose = flag.Bool("verbose", false, "Enable verbose output")
		stats   = flag.Bool("stats", false, "Show statistics after scanning")
	)

	flag.Parse()

	if *showVer {
		fmt.Println(version)
		fmt.Println("A tool for scanning Go code for issues")
		return
	}

	if *help {
		showHelp()
		return
	}

	if *path == "" {
		log.Fatal("❌ Error: --path flag is required")
	}

	if !*html && !*json {
		log.Fatal("❌ Error: At least one of --html or --json must be true")
	}

	if *verbose {
		fmt.Printf("🔍 Scanning path: %s\n", *path)
		fmt.Printf("  HTML report: %v\n", *html)
		fmt.Printf("  JSON report: %v\n", *json)
	}

	err := Scan(*path, *html, *json)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *stats {
		// Giả sử có hàm GetStats trả về map[string]int
		if *verbose {
			fmt.Println("📊 Statistics:")
		}
		// statsMap := GetStats() // Cần cài đặt hàm này nếu có
		// for k, v := range statsMap {
		// 	fmt.Printf("  %s: %d\n", k, v)
		// }
	}
}
