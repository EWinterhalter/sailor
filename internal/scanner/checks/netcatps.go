package checks

import (
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckNetcatProcess(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Netcat Process Check")

	check := models.CheckResult{
		Name:        "Netcat Process",
		Description: "Checks if netcat is running inside the container",
		Issues:      []string{},
	}

	output, _ := docker.Exec(containerID, []string{
		"sh", "-c", "ps aux | grep -E 'nc|netcat|ncat' | grep -v grep",
	})

	check.Output = output
	check.Duration = time.Since(start)

	if strings.TrimSpace(output) != "" {
		check.Status = "fail"
		check.Severity = "critical"
		check.Issues = append(check.Issues, "Netcat process is running")
		printerln.PrintCheckResult("Netcat Process", "FAIL", check.Duration, "nc detected")
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Netcat Process", "PASS", check.Duration, "not running")
	}

	return check
}
