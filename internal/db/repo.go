package db

import (
	"database/sql"
	"fmt"

	"github.com/EWinterhalter/sailor/internal/scanner/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveScanReport(report *models.FinalReport) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var scanID int64
	err = tx.QueryRow(`
		INSERT INTO scans (
			timestamp, image, container_id, total_time_ms,
			total_checks, passed, warnings, failed,
			critical, high, medium, low,
			has_failures, has_warnings
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`,
		report.Timestamp,
		report.Scan.Image,
		report.Scan.ContainerID,
		report.Scan.TotalTimeMs,
		report.Results.Summary.TotalChecks,
		report.Results.Summary.Passed,
		report.Results.Summary.Warnings,
		report.Results.Summary.Failed,
		report.Results.Summary.Critical,
		report.Results.Summary.High,
		report.Results.Summary.Medium,
		report.Results.Summary.Low,
		report.Results.Summary.HasFailures,
		report.Results.Summary.HasWarnings,
	).Scan(&scanID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert scan: %w", err)
	}

	for _, check := range report.Results.Checks {
		var checkResultID int64

		output := check.Output
		if output == "" && check.OutputPreview != "" {
			output = check.OutputPreview
		}

		err = tx.QueryRow(`
			INSERT INTO check_results (
				scan_id, name, description, severity, status, duration_ms, output, output_preview
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`,
			scanID,
			check.Name,
			check.Description,
			check.Severity,
			check.Status,
			check.DurationMs,
			output,
			check.OutputPreview,
		).Scan(&checkResultID)

		if err != nil {
			return 0, fmt.Errorf("failed to insert check result: %w", err)
		}

		for _, issue := range check.Issues {
			_, err = tx.Exec(`
				INSERT INTO check_issues (check_result_id, issue)
				VALUES ($1, $2)
			`, checkResultID, issue)

			if err != nil {
				return 0, fmt.Errorf("failed to insert issue: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return scanID, nil
}
