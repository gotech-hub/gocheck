package report

import (
	"html/template"
	"os"

	"github.com/gotech-hub/gocheck/analyzer"
)

func GenerateHTML(findings []analyzer.Finding) {
	tmpl := `
    <html>
    <head>
    <style>
        body { font-family: Arial; background: #f9f9f9; padding: 20px; }
        h1 { color: #333; }
        .finding { padding: 10px; border-radius: 5px; margin-bottom: 10px; }
        .Low { background-color: #f6ffed; border-left: 5px solid #95de64; }
        .Medium { background-color: #fffbe6; border-left: 5px solid #fadb14; }
        .High { background-color: #fff2e8; border-left: 5px solid #fa8c16; }
        .Critical { background-color: #fff1f0; border-left: 5px solid #ff4d4f; }
        .file { font-weight: bold; }
        .suggestion { font-style: italic; color: #555; }
    </style>
    </head>
    <body>
        <h1>GoCheck Report</h1>
        {{range .}}
            <div class="finding {{.Severity}}">
                <div class="file">{{.File}}:{{.Line}}</div>
                <div>{{.Message}}</div>
                <div class="suggestion">💡 {{.Suggestion}}</div>
            </div>
        {{end}}
    </body>
    </html>`

	t := template.Must(template.New("report").Parse(tmpl))
	f, _ := os.Create("report.html")
	defer f.Close()
	t.Execute(f, findings)
}
