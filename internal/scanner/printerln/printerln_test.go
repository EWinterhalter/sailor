package printerln

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/EWinterhalter/sailor/internal/colors"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
)

func captureOutput(t *testing.T, f func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()
	return buf.String()
}

func withNoColors(t *testing.T) {
	t.Helper()
	origReset := colors.ColorReset
	origRed := colors.ColorRed
	origGreen := colors.ColorGreen
	origYellow := colors.ColorYellow
	origBlue := colors.ColorBlue
	origMagenta := colors.ColorMagenta
	origCyan := colors.ColorCyan
	origWhite := colors.ColorWhite
	origBold := colors.ColorBold
	origDim := colors.ColorDim

	colors.DisableColors()

	t.Cleanup(func() {
		colors.ColorReset = origReset
		colors.ColorRed = origRed
		colors.ColorGreen = origGreen
		colors.ColorYellow = origYellow
		colors.ColorBlue = origBlue
		colors.ColorMagenta = origMagenta
		colors.ColorCyan = origCyan
		colors.ColorWhite = origWhite
		colors.ColorBold = origBold
		colors.ColorDim = origDim
	})
}

// ---- PrintHeader ----

func TestPrintHeader_ContainsContainerID(t *testing.T) {
	withNoColors(t)
	out := captureOutput(t, func() {
		PrintHeader("abcdef123456789")
	})
	if !strings.Contains(out, "abcdef123456") {
		t.Errorf("expected container ID in output, got: %q", out)
	}
}

func TestPrintHeader_ContainsDate(t *testing.T) {
	withNoColors(t)
	out := captureOutput(t, func() {
		PrintHeader("abcdef123456789")
	})
	year := time.Now().Format("2006")
	if !strings.Contains(out, year) {
		t.Errorf("expected current year in output, got: %q", out)
	}
}

// ---- PrintCheckResult ----

func TestPrintCheckResult_Pass(t *testing.T) {
	withNoColors(t)
	activeSpinner = nil
	out := captureOutput(t, func() {
		PrintCheckResult("Root User", "PASS", 10*time.Millisecond, "non-root")
	})
	if !strings.Contains(out, "✓") {
		t.Errorf("expected ✓ in PASS output, got: %q", out)
	}
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS in output, got: %q", out)
	}
	if !strings.Contains(out, "Root User") {
		t.Errorf("expected name in output, got: %q", out)
	}
	if !strings.Contains(out, "non-root") {
		t.Errorf("expected details in output, got: %q", out)
	}
}

func TestPrintCheckResult_Warn(t *testing.T) {
	withNoColors(t)
	activeSpinner = nil
	out := captureOutput(t, func() {
		PrintCheckResult("Open Ports", "WARN", 5*time.Millisecond, "many ports")
	})
	if !strings.Contains(out, "⚠") {
		t.Errorf("expected ⚠ in WARN output, got: %q", out)
	}
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %q", out)
	}
}

func TestPrintCheckResult_Fail(t *testing.T) {
	withNoColors(t)
	activeSpinner = nil
	out := captureOutput(t, func() {
		PrintCheckResult("Environment", "FAIL", 2*time.Millisecond, "found secret")
	})
	if !strings.Contains(out, "✗") {
		t.Errorf("expected ✗ in FAIL output, got: %q", out)
	}
	if !strings.Contains(out, "FAIL") {
		t.Errorf("expected FAIL in output, got: %q", out)
	}
	if !strings.Contains(out, "found secret") {
		t.Errorf("expected details in output, got: %q", out)
	}
}

func TestPrintCheckResult_UnknownStatus(t *testing.T) {
	withNoColors(t)
	activeSpinner = nil
	out := captureOutput(t, func() {
		PrintCheckResult("Check", "UNKNOWN", 0, "")
	})
	if !strings.Contains(out, "?") {
		t.Errorf("expected ? for unknown status, got: %q", out)
	}
}

func TestPrintCheckResult_NoDetails(t *testing.T) {
	withNoColors(t)
	activeSpinner = nil
	out := captureOutput(t, func() {
		PrintCheckResult("Check", "PASS", 1*time.Millisecond, "")
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS in output, got: %q", out)
	}
}

func TestPrintCheckResult_ContainsDuration(t *testing.T) {
	withNoColors(t)
	activeSpinner = nil
	out := captureOutput(t, func() {
		PrintCheckResult("Check", "PASS", 42*time.Millisecond, "")
	})
	if !strings.Contains(out, "42ms") {
		t.Errorf("expected duration 42ms in output, got: %q", out)
	}
}

// ---- PrintCheckStart + PrintCheckResult (spinner lifecycle) ----

func TestPrintCheckStart_ThenResult(t *testing.T) {
	withNoColors(t)
	out := captureOutput(t, func() {
		PrintCheckStart("Image Scan")
		time.Sleep(90 * time.Millisecond)
		PrintCheckResult("Image Scan", "PASS", 90*time.Millisecond, "ok")
	})
	if activeSpinner != nil {
		t.Error("activeSpinner should be nil after PrintCheckResult")
	}
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS in output, got: %q", out)
	}
}

func TestPrintCheckStart_AlreadyActive(t *testing.T) {
	withNoColors(t)
	captureOutput(t, func() {
		PrintCheckStart("First")
		activeSpinner.start("Second")
		PrintCheckResult("First", "PASS", 0, "")
	})
}

// ---- PrintTimeout ----

func TestPrintTimeout_NoSpinner(t *testing.T) {
	withNoColors(t)
	activeSpinner = nil
	out := captureOutput(t, func() {
		PrintTimeout()
	})
	if !strings.Contains(out, "TIMEOUT") {
		t.Errorf("expected TIMEOUT in output, got: %q", out)
	}
}

func TestPrintTimeout_WithActiveSpinner(t *testing.T) {
	withNoColors(t)
	out := captureOutput(t, func() {
		PrintCheckStart("Slow Check")
		PrintTimeout()
	})
	if activeSpinner != nil {
		t.Error("activeSpinner should be nil after PrintTimeout")
	}
	if !strings.Contains(out, "TIMEOUT") {
		t.Errorf("expected TIMEOUT in output, got: %q", out)
	}
}

// ---- PrintSummary ----

func TestPrintSummary_Counts(t *testing.T) {
	withNoColors(t)
	result := &models.ScanResult{
		Summary: models.Summary{
			TotalChecks: 5,
			Passed:      3,
			Warnings:    1,
			Failed:      1,
		},
		TotalTime: 200 * time.Millisecond,
	}
	out := captureOutput(t, func() {
		PrintSummary(result)
	})
	if !strings.Contains(out, "5") {
		t.Errorf("expected total checks 5, got: %q", out)
	}
	if !strings.Contains(out, "Passed") {
		t.Errorf("expected Passed in output, got: %q", out)
	}
	if !strings.Contains(out, "Failed") {
		t.Errorf("expected Failed in output, got: %q", out)
	}
}

func TestPrintSummary_SeverityBreakdown(t *testing.T) {
	withNoColors(t)
	result := &models.ScanResult{
		Summary: models.Summary{
			TotalChecks: 4,
			Critical:    1,
			High:        1,
			Medium:      1,
			Low:         1,
		},
	}
	out := captureOutput(t, func() {
		PrintSummary(result)
	})
	for _, s := range []string{"Critical", "High", "Medium", "Low"} {
		if !strings.Contains(out, s) {
			t.Errorf("expected %s in severity breakdown, got: %q", s, out)
		}
	}
}

func TestPrintSummary_SeverityBreakdown_ZeroValuesHidden(t *testing.T) {
	withNoColors(t)
	result := &models.ScanResult{
		Summary: models.Summary{TotalChecks: 1, Passed: 1},
	}
	out := captureOutput(t, func() {
		PrintSummary(result)
	})
	for _, s := range []string{"Critical", "High", "Medium"} {
		if strings.Contains(out, s) {
			t.Errorf("expected %s to be hidden when count is 0, got: %q", s, out)
		}
	}
}

func TestPrintSummary_HasIssues(t *testing.T) {
	withNoColors(t)
	result := &models.ScanResult{
		HasIssues: true,
		Summary:   models.Summary{TotalChecks: 1, Failed: 1},
	}
	out := captureOutput(t, func() {
		PrintSummary(result)
	})
	if !strings.Contains(out, "SECURITY ISSUES DETECTED") {
		t.Errorf("expected security warning, got: %q", out)
	}
}

func TestPrintSummary_NoIssues(t *testing.T) {
	withNoColors(t)
	result := &models.ScanResult{
		HasIssues: false,
		Summary:   models.Summary{TotalChecks: 1, Passed: 1},
	}
	out := captureOutput(t, func() {
		PrintSummary(result)
	})
	if !strings.Contains(out, "NO CRITICAL ISSUES FOUND") {
		t.Errorf("expected no-issues message, got: %q", out)
	}
}

func TestPrintSummary_TotalTime(t *testing.T) {
	withNoColors(t)
	result := &models.ScanResult{
		TotalTime: 123 * time.Millisecond,
	}
	out := captureOutput(t, func() {
		PrintSummary(result)
	})
	if !strings.Contains(out, "123ms") {
		t.Errorf("expected 123ms in output, got: %q", out)
	}
}
