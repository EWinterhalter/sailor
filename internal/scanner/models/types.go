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
	TotalChecks int `json:"total_checks"`
	Passed      int `json:"passed"`
	Warnings    int `json:"warnings"`
	Failed      int `json:"failed"`
	Critical    int `json:"critical"`
	High        int `json:"high"`
	Medium      int `json:"medium"`
	Low         int `json:"low"`
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
