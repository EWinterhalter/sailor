package models

import "time"

type CheckResult struct {
	Name        string        `json:"name"`
	Status      string        `json:"status"`
	Output      string        `json:"output"`
	Issues      []string      `json:"issues,omitempty"`
	Duration    time.Duration `json:"duration"`
	Severity    string        `json:"severity"`
	Description string        `json:"description"`
}

type ScanResult struct {
	Checks    []CheckResult `json:"checks"`
	Summary   Summary       `json:"summary"`
	HasIssues bool          `json:"has_issues"`
	TotalTime time.Duration `json:"total_time"`
}

type Summary struct {
	TotalChecks int  `json:"total_checks"`
	Passed      int  `json:"passed"`
	Warnings    int  `json:"warnings"`
	Failed      int  `json:"failed"`
	Critical    int  `json:"critical"`
	High        int  `json:"high"`
	Medium      int  `json:"medium"`
	Low         int  `json:"low"`
	HasFailures bool `json:"has_failures"`
	HasWarnings bool `json:"has_warnings"`
}

type FinalReport struct {
	Timestamp time.Time   `json:"timestamp"`
	Scan      ScanInfo    `json:"scan"`
	Results   ResultsInfo `json:"results"`
}

type ScanInfo struct {
	Image       string  `json:"image"`
	ContainerID string  `json:"container_id"`
	TotalTimeMs float64 `json:"total_time_ms"`
}

type ResultsInfo struct {
	Summary Summary      `json:"summary"`
	Checks  []CheckEntry `json:"checks"`
}

type CheckEntry struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Severity      string   `json:"severity"`
	Status        string   `json:"status"`
	DurationMs    float64  `json:"duration_ms"`
	Output        string   `json:"output,omitempty"`
	OutputPreview string   `json:"output_preview,omitempty"`
	OutputFull    *string  `json:"output_full"`
	Issues        []string `json:"issues"`
}

func (sr *ScanResult) CalculateSummary() {
	sr.Summary.TotalChecks = len(sr.Checks)
	for _, check := range sr.Checks {
		switch check.Status {
		case "pass":
			sr.Summary.Passed++
		case "warn":
			sr.Summary.Warnings++
		case "fail":
			sr.Summary.Failed++
			sr.HasIssues = true
		}

		switch check.Severity {
		case "critical":
			sr.Summary.Critical++
		case "high":
			sr.Summary.High++
		case "medium":
			sr.Summary.Medium++
		case "low":
			sr.Summary.Low++
		}
	}
}
