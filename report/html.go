package report

import (
	"html/template"
	"os"
	"strings"

	"github.com/gotech-hub/gocheck/analyzer"
)

func GenerateHTML(findings []analyzer.Finding) {
	// Tính toán số lượng từng loại severity
	stats := map[analyzer.Severity]int{
		analyzer.Low:      0,
		analyzer.Medium:   0,
		analyzer.High:     0,
		analyzer.Critical: 0,
	}
	for _, f := range findings {
		stats[f.Severity]++
	}
	type ReportData struct {
		Findings            []analyzer.Finding
		Stats               map[analyzer.Severity]int
		Total               int
		CleanCodeFindings   []analyzer.Finding
		PerformanceFindings []analyzer.Finding
		SecurityFindings    []analyzer.Finding
	}
	// Phân loại findings
	var cleanCodeFindings, performanceFindings, securityFindings []analyzer.Finding
	for _, f := range findings {
		// Dựa vào Message để phân loại (theo cách gọi trong analyzer.go)
		if isCleanCodeFinding(f) {
			cleanCodeFindings = append(cleanCodeFindings, f)
		} else if isPerformanceFinding(f) {
			performanceFindings = append(performanceFindings, f)
		} else if isSecurityFinding(f) {
			securityFindings = append(securityFindings, f)
		}
	}
	data := ReportData{
		Findings:            findings,
		Stats:               stats,
		Total:               len(findings),
		CleanCodeFindings:   cleanCodeFindings,
		PerformanceFindings: performanceFindings,
		SecurityFindings:    securityFindings,
	}

	tmpl := `
    <html>
    <head>
    <style>
        body { font-family: Arial; background: #f9f9f9; padding: 20px; }
        h1 { color: #333; }
        .stats {
            background: #fff; border-radius: 8px; padding: 16px; margin-bottom: 24px; box-shadow: 0 2px 8px #eee;
            display: flex; gap: 24px; align-items: center;
        }
        .stat {
            display: flex; flex-direction: column; align-items: center; min-width: 80px;
        }
        .stat-label { font-size: 14px; color: #888; }
        .stat-value { font-size: 24px; font-weight: bold; }
        .Low { background-color: #f6ffed; border-left: 5px solid #95de64; }
        .Medium { background-color: #fffbe6; border-left: 5px solid #fadb14; }
        .High { background-color: #fff2e8; border-left: 5px solid #fa8c16; }
        .Critical { background-color: #fff1f0; border-left: 5px solid #ff4d4f; }
        .finding { padding: 10px; border-radius: 5px; margin-bottom: 10px; }
        .file { font-weight: bold; }
        .suggestion { font-style: italic; color: #555; }
        .tab { display: inline-block; padding: 10px 24px; margin-right: 8px; background: #eee; border-radius: 8px 8px 0 0; cursor: pointer; }
        .tab.active { background: #fff; border-bottom: 2px solid #fff; font-weight: bold; }
        .tab-content { display: none; background: #fff; border-radius: 0 0 8px 8px; padding: 16px; box-shadow: 0 2px 8px #eee; }
        .tab-content.active { display: block; }
    </style>
    <script>
        function showTab(tab) {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
            document.getElementById(tab+'-tab').classList.add('active');
            document.getElementById(tab+'-content').classList.add('active');
        }
        window.onload = function() { showTab('cleancode'); };
    </script>
    </head>
    <body>
        <h1>GoCheck Report</h1>
        <div class="stats">
            <div class="stat">
                <span class="stat-label">Total</span>
                <span class="stat-value">{{.Total}}</span>
            </div>
            <div class="stat Low">
                <span class="stat-label">Low</span>
                <span class="stat-value">{{index .Stats "Low"}}</span>
            </div>
            <div class="stat Medium">
                <span class="stat-label">Medium</span>
                <span class="stat-value">{{index .Stats "Medium"}}</span>
            </div>
            <div class="stat High">
                <span class="stat-label">High</span>
                <span class="stat-value">{{index .Stats "High"}}</span>
            </div>
            <div class="stat Critical">
                <span class="stat-label">Critical</span>
                <span class="stat-value">{{index .Stats "Critical"}}</span>
            </div>
        </div>
        <div>
            <div id="cleancode-tab" class="tab" onclick="showTab('cleancode')">Clean Code</div>
            <div id="performance-tab" class="tab" onclick="showTab('performance')">Performance</div>
            <div id="security-tab" class="tab" onclick="showTab('security')">Security</div>
        </div>
        <div id="cleancode-content" class="tab-content">
            {{range .CleanCodeFindings}}
                <div class="finding {{.Severity}}">
                    <div class="file">{{.File}}:{{.Line}}</div>
                    <div>{{.Message}}</div>
                    <div class="suggestion">💡 {{.Suggestion}}</div>
                </div>
            {{else}}
                <div>No Clean Code findings.</div>
            {{end}}
        </div>
        <div id="performance-content" class="tab-content">
            {{range .PerformanceFindings}}
                <div class="finding {{.Severity}}">
                    <div class="file">{{.File}}:{{.Line}}</div>
                    <div>{{.Message}}</div>
                    <div class="suggestion">💡 {{.Suggestion}}</div>
                </div>
            {{else}}
                <div>No Performance findings.</div>
            {{end}}
        </div>
        <div id="security-content" class="tab-content">
            {{range .SecurityFindings}}
                <div class="finding {{.Severity}}">
                    <div class="file">{{.File}}:{{.Line}}</div>
                    <div>{{.Message}}</div>
                    <div class="suggestion">💡 {{.Suggestion}}</div>
                </div>
            {{else}}
                <div>No Security findings.</div>
            {{end}}
        </div>
    </body>
    </html>`

	t := template.Must(template.New("report").Parse(tmpl))
	f, _ := os.Create("report.html")
	defer f.Close()
	t.Execute(f, data)
}

// Thêm các hàm phân loại findings
func isCleanCodeFinding(f analyzer.Finding) bool {
	// Có thể dựa vào các từ khóa đặc trưng trong Message
	// (hoặc tốt hơn: thêm trường Type vào Finding, nhưng tạm thời dùng heuristic)
	msg := f.Message
	return containsAny(msg, []string{"function", "variable", "comment", "magic number", "global variable", "nested", "return statement", "if/else", "local variable", "function name", "commented-out code"})
}
func isPerformanceFinding(f analyzer.Finding) bool {
	msg := f.Message
	return containsAny(msg, []string{"For-loop", "performance", "defer", "goroutine", "string concatenation"})
}
func isSecurityFinding(f analyzer.Finding) bool {
	msg := f.Message
	return containsAny(msg, []string{"credential", "exec.Command", "ListenAndServe", "hash function", "md5", "sha1", "insecure"})
}
func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(strings.ToLower(s), strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
