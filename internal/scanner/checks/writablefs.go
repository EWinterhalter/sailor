package checks

import (
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckWritableFS(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Writable Filesystem Check")

	check := models.CheckResult{
		Name:        "Writable FS",
		Description: "Checks if root filesystem is writable",
		Issues:      []string{},
	}

	output, _ := docker.Exec(containerID, []string{"mount"})
	check.Output = output
	check.Duration = time.Since(start)

	if strings.Contains(output, "rw") {
		check.Status = "warn"
		check.Severity = "medium"
		check.Issues = append(check.Issues, "Root filesystem is writable")
		printerln.PrintCheckResult("Writable FS", "WARN", check.Duration, "Filesystem is writable")
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Writable FS", "PASS", check.Duration, "Filesystem is read-only")
	}

	return check
}
