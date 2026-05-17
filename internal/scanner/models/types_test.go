package models

import (
	"testing"
	"time"
)

func TestCalculateSummary_Empty(t *testing.T) {
	sr := &ScanResult{}
	sr.CalculateSummary()

	if sr.Summary.TotalChecks != 0 {
		t.Errorf("expected 0 total checks, got %d", sr.Summary.TotalChecks)
	}
	if sr.HasIssues {
		t.Error("expected HasIssues=false for empty scan")
	}
}

func TestCalculateSummary_Counts(t *testing.T) {
	sr := &ScanResult{
		Checks: []CheckResult{
			{Status: "pass", Severity: "low", Duration: time.Millisecond},
			{Status: "pass", Severity: "medium", Duration: time.Millisecond},
			{Status: "warn", Severity: "high", Duration: time.Millisecond},
			{Status: "fail", Severity: "critical", Duration: time.Millisecond},
			{Status: "fail", Severity: "high", Duration: time.Millisecond},
		},
	}

	sr.CalculateSummary()

	if sr.Summary.TotalChecks != 5 {
		t.Errorf("expected 5, got %d", sr.Summary.TotalChecks)
	}
	if sr.Summary.Passed != 2 {
		t.Errorf("expected 2 passed, got %d", sr.Summary.Passed)
	}
	if sr.Summary.Warnings != 1 {
		t.Errorf("expected 1 warning, got %d", sr.Summary.Warnings)
	}
	if sr.Summary.Failed != 2 {
		t.Errorf("expected 2 failed, got %d", sr.Summary.Failed)
	}
	if sr.Summary.Critical != 1 {
		t.Errorf("expected 1 critical, got %d", sr.Summary.Critical)
	}
	if sr.Summary.High != 2 {
		t.Errorf("expected 2 high, got %d", sr.Summary.High)
	}
	if sr.Summary.Medium != 1 {
		t.Errorf("expected 1 medium, got %d", sr.Summary.Medium)
	}
	if sr.Summary.Low != 1 {
		t.Errorf("expected 1 low, got %d", sr.Summary.Low)
	}
	if !sr.HasIssues {
		t.Error("expected HasIssues=true when there are failures")
	}
}

func TestCalculateSummary_NoFailures(t *testing.T) {
	sr := &ScanResult{
		Checks: []CheckResult{
			{Status: "pass", Severity: "low"},
			{Status: "warn", Severity: "medium"},
		},
	}

	sr.CalculateSummary()

	if sr.HasIssues {
		t.Error("expected HasIssues=false when no failures")
	}
	if sr.Summary.Failed != 0 {
		t.Errorf("expected 0 failed, got %d", sr.Summary.Failed)
	}
}

func TestCalculateSummary_AllSeverities(t *testing.T) {
	sr := &ScanResult{
		Checks: []CheckResult{
			{Status: "pass", Severity: "critical"},
			{Status: "pass", Severity: "high"},
			{Status: "pass", Severity: "medium"},
			{Status: "pass", Severity: "low"},
		},
	}

	sr.CalculateSummary()

	if sr.Summary.Critical != 1 {
		t.Errorf("expected 1 critical, got %d", sr.Summary.Critical)
	}
	if sr.Summary.High != 1 {
		t.Errorf("expected 1 high, got %d", sr.Summary.High)
	}
	if sr.Summary.Medium != 1 {
		t.Errorf("expected 1 medium, got %d", sr.Summary.Medium)
	}
	if sr.Summary.Low != 1 {
		t.Errorf("expected 1 low, got %d", sr.Summary.Low)
	}
}
