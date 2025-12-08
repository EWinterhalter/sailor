package report

import (
	"time"

	"github.com/EWinterhalter/sailor/internal/scanner/models"
)

func BuildReport(image string, containerID string, scan *models.ScanResult) *models.FinalReport {
	checks := make([]models.CheckEntry, 0, len(scan.Checks))

	for _, c := range scan.Checks {
		entry := models.CheckEntry{
			Name:        c.Name,
			Description: c.Description,
			Severity:    c.Severity,
			Status:      c.Status,
			DurationMs:  float64(c.Duration.Milliseconds()),
			Issues:      c.Issues,
		}

		if len(c.Output) > 300 {
			entry.OutputPreview = c.Output[:300] + "..."
			entry.OutputFull = nil
		} else {
			entry.Output = c.Output
			entry.OutputFull = nil
		}

		checks = append(checks, entry)
	}

	summary := models.Summary{
		TotalChecks: scan.Summary.TotalChecks,
		Passed:      scan.Summary.Passed,
		Warnings:    scan.Summary.Warnings,
		Failed:      scan.Summary.Failed,
		Critical:    scan.Summary.Critical,
		High:        scan.Summary.High,
		Medium:      scan.Summary.Medium,
		Low:         scan.Summary.Low,
		HasFailures: scan.Summary.Failed > 0 || scan.Summary.Critical > 0,
		HasWarnings: scan.Summary.Warnings > 0,
	}

	return &models.FinalReport{
		Timestamp: time.Now(),
		Scan: models.ScanInfo{
			Image:       image,
			ContainerID: containerID[:12],
			TotalTimeMs: float64(scan.TotalTime.Milliseconds()),
		},
		Results: models.ResultsInfo{
			Summary: summary,
			Checks:  checks,
		},
	}
}
