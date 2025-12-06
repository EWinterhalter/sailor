package checks

import (
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckRootUser(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Root User Check")

	check := models.CheckResult{
		Name:        "Root User",
		Description: "Checks if container runs as root",
		Issues:      []string{},
	}

	output, _ := docker.Exec(containerID, []string{"id", "-u"})
	check.Output = output
	check.Duration = time.Since(start)

	if strings.TrimSpace(output) == "0" {
		check.Status = "fail"
		check.Severity = "high"
		check.Issues = append(check.Issues, "Container runs as root (UID=0)")
		printerln.PrintCheckResult("Root User", "FAIL", check.Duration, "UID=0")
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Root User", "PASS", check.Duration, "Non-root user")
	}

	return check
}
