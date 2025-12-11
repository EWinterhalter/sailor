package checks

import (
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
		Description: "Checks if writing to root filesystem is possible",
		Issues:      []string{},
	}

	testCmd := []string{"sh", "-c", "touch /.__sailor_test__ 2>/dev/null && echo OK || echo FAIL"}
	output, _ := docker.Exec(containerID, testCmd)

	check.Output = output
	check.Duration = time.Since(start)

	if output == "OK\n" || output == "OK" {
		check.Status = "warn"
		check.Severity = "medium"
		check.Issues = append(check.Issues, "Root filesystem is writable")
		printerln.PrintCheckResult("Writable FS", "WARN", check.Duration, "Filesystem is writable")
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Writable FS", "PASS", check.Duration, "Filesystem is read-only")
	}

	_, _ = docker.Exec(containerID, []string{"sh", "-c", "rm -f /.__sailor_test__ 2>/dev/null"})

	return check
}
