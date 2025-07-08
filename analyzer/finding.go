package analyzer

type Severity string

const (
	Low      Severity = "Low"
	Medium   Severity = "Medium"
	High     Severity = "High"
	Critical Severity = "Critical"
)

type Finding struct {
	File       string   `json:"file"`
	Line       int      `json:"line"`
	Message    string   `json:"message"`
	Severity   Severity `json:"severity"`
	Suggestion string   `json:"suggestion"`
	Category   string   `json:"category"` // e.g., "Clean", "Performance", "Security"
}
