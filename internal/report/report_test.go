package report

import (
	"strings"
	"testing"
	"time"

	"github.com/EWinterhalter/sailor/internal/scanner/models"
)

func makeScan(checks []models.CheckResult) *models.ScanResult {
	sr := &models.ScanResult{
		Checks:    checks,
		TotalTime: 500 * time.Millisecond,
	}
	sr.CalculateSummary()
	return sr
}

func TestBuildReport_BasicFields(t *testing.T) {
	scan := makeScan([]models.CheckResult{
		{Name: "Root User", Status: "pass", Severity: "low", Description: "uid check", Issues: []string{}},
	})

	report := BuildReport("my-image:latest", "abcdef123456789", scan)

	if report.Scan.Image != "my-image:latest" {
		t.Errorf("unexpected image: %s", report.Scan.Image)
	}
	if report.Scan.ContainerID != "abcdef123456" {
		t.Errorf("container ID should be truncated to 12 chars, got %s", report.Scan.ContainerID)
	}
	if report.Scan.TotalTimeMs != float64(scan.TotalTime.Milliseconds()) {
		t.Errorf("unexpected total time: %f", report.Scan.TotalTimeMs)
	}
	if report.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestBuildReport_CheckEntry(t *testing.T) {
	scan := makeScan([]models.CheckResult{
		{
			Name:        "Open Ports",
			Description: "port scan",
			Severity:    "high",
			Status:      "fail",
			Duration:    100 * time.Millisecond,
			Output:      "port 4444 open",
			Issues:      []string{"Suspicious port 4444"},
		},
	})

	report := BuildReport("img", "abcdef123456789", scan)

	if len(report.Results.Checks) != 1 {
		t.Fatalf("expected 1 check entry, got %d", len(report.Results.Checks))
	}

	entry := report.Results.Checks[0]
	if entry.Name != "Open Ports" {
		t.Errorf("unexpected name: %s", entry.Name)
	}
	if entry.Status != "fail" {
		t.Errorf("unexpected status: %s", entry.Status)
	}
	if entry.DurationMs != 100 {
		t.Errorf("unexpected duration: %f", entry.DurationMs)
	}
	if len(entry.Issues) != 1 || entry.Issues[0] != "Suspicious port 4444" {
		t.Errorf("unexpected issues: %v", entry.Issues)
	}
}

func TestBuildReport_LongOutputTruncated(t *testing.T) {
	longOutput := strings.Repeat("x", 400)
	scan := makeScan([]models.CheckResult{
		{Name: "Check", Status: "pass", Severity: "low", Output: longOutput, Issues: []string{}},
	})

	report := BuildReport("img", "abcdef123456789", scan)
	entry := report.Results.Checks[0]

	if entry.Output != "" {
		t.Error("Output should be empty when output > 300 chars")
	}
	if !strings.HasSuffix(entry.OutputPreview, "...") {
		t.Error("OutputPreview should end with ...")
	}
	if len(entry.OutputPreview) != 303 {
		t.Errorf("OutputPreview should be 303 chars, got %d", len(entry.OutputPreview))
	}
}

func TestBuildReport_ShortOutputNotTruncated(t *testing.T) {
	scan := makeScan([]models.CheckResult{
		{Name: "Check", Status: "pass", Severity: "low", Output: "short output", Issues: []string{}},
	})

	report := BuildReport("img", "abcdef123456789", scan)
	entry := report.Results.Checks[0]

	if entry.Output != "short output" {
		t.Errorf("short output should pass through unchanged, got %q", entry.Output)
	}
	if entry.OutputPreview != "" {
		t.Error("OutputPreview should be empty for short output")
	}
}

func TestBuildReport_SummaryHasFailures(t *testing.T) {
	scan := makeScan([]models.CheckResult{
		{Status: "fail", Severity: "critical", Issues: []string{"issue"}},
		{Status: "pass", Severity: "low", Issues: []string{}},
	})

	report := BuildReport("img", "abcdef123456789", scan)

	if !report.Results.Summary.HasFailures {
		t.Error("expected HasFailures=true")
	}
	if report.Results.Summary.TotalChecks != 2 {
		t.Errorf("expected 2 total checks, got %d", report.Results.Summary.TotalChecks)
	}
}

func TestBuildReport_SummaryHasWarnings(t *testing.T) {
	scan := makeScan([]models.CheckResult{
		{Status: "warn", Severity: "medium", Issues: []string{"something"}},
	})

	report := BuildReport("img", "abcdef123456789", scan)

	if !report.Results.Summary.HasWarnings {
		t.Error("expected HasWarnings=true")
	}
	if report.Results.Summary.HasFailures {
		t.Error("expected HasFailures=false when only warnings")
	}
}

func TestBuildReport_EmptyScan(t *testing.T) {
	scan := makeScan([]models.CheckResult{})
	report := BuildReport("img", "abcdef123456789", scan)

	if len(report.Results.Checks) != 0 {
		t.Errorf("expected 0 check entries, got %d", len(report.Results.Checks))
	}
	if report.Results.Summary.TotalChecks != 0 {
		t.Errorf("expected 0 total checks, got %d", report.Results.Summary.TotalChecks)
	}
}
