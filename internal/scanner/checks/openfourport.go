package checks

import (
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckOpenPort4444(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Port 4444 Check")

	check := models.CheckResult{
		Name:        "Open Port 4444",
		Description: "Checks if port 4444 is open inside container",
		Issues:      []string{},
	}

	output, _ := docker.Exec(containerID, []string{
		"sh", "-c", "ss -lntp | grep ':4444'",
	})

	check.Output = output
	check.Duration = time.Since(start)

	if strings.TrimSpace(output) != "" {
		check.Status = "fail"
		check.Severity = "critical"
		check.Issues = append(check.Issues, "Port 4444 is open")
		printerln.PrintCheckResult("Port 4444", "FAIL", check.Duration, "listening")
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Port 4444", "PASS", check.Duration, "closed")
	}

	return check
}
