package checks

import (
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckImageHistory(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Image History Check")

	check := models.CheckResult{
		Name:        "Image History",
		Description: "Check layers for secrets or sensitive files",
		Issues:      []string{},
	}

	output, _ := docker.ImageHistory(containerID)
	check.Output = output
	check.Duration = time.Since(start)

	suspiciousKeywords := []string{"PASSWORD", "SECRET", "TOKEN", "API_KEY"}
	found := []string{}

	for _, key := range suspiciousKeywords {
		if strings.Contains(strings.ToUpper(output), key) {
			found = append(found, key)
		}
	}

	if len(found) > 0 {
		check.Status = "fail"
		check.Severity = "high"
		check.Issues = append(check.Issues, "Secrets found in image history: "+strings.Join(found, ", "))
		printerln.PrintCheckResult("Image History", "FAIL", check.Duration, "Secrets detected: "+strings.Join(found, ", "))
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Image History", "PASS", check.Duration, "No secrets found in layers")
	}

	return check
}
