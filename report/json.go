package report

import (
	"encoding/json"
	"os"

	"github.com/gotech-hub/gocheck/analyzer"
)

func GenerateJSON(findings []analyzer.Finding) {
	f, _ := os.Create("report.json")
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.Encode(findings)
}
