package checks

import (
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckNetcatListener(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Netcat Listener Check")

	check := models.CheckResult{
		Name:        "Netcat Listener",
		Description: "Checks if netcat is running in listen mode",
		Issues:      []string{},
	}

	output, _ := docker.Exec(containerID, []string{
		"sh", "-c", "ps aux | grep -E 'nc .* -l' | grep -v grep",
	})

	check.Output = output
	check.Duration = time.Since(start)

	if strings.TrimSpace(output) != "" {
		check.Status = "fail"
		check.Severity = "critical"
		check.Issues = append(check.Issues, "Netcat running in listen mode")
		printerln.PrintCheckResult("Netcat Listener", "FAIL", check.Duration, "bind shell detected")
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Netcat Listener", "PASS", check.Duration, "not detected")
	}

	return check
}
