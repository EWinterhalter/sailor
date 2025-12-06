package checks

import (
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

var sensitiveKeys = []string{"PASSWORD", "SECRET", "TOKEN", "API_KEY", "PRIVATE_KEY", "AWS_SECRET", "DB_PASS"}

func CheckEnvironment(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Environment Variables")

	check := models.CheckResult{
		Name:        "Environment",
		Description: "Checking for sensitive data in environment",
		Issues:      []string{},
	}

	output, err := docker.Exec(containerID, []string{"env"})
	check.Output = output
	check.Duration = time.Since(start)

	if err != nil {
		check.Status = "warn"
		check.Severity = "low"
		printerln.PrintCheckResult("Environment", "WARN", check.Duration, "Unable to read environment")
		return check
	}

	foundVars := findSensitiveEnv(output)

	if len(foundVars) > 0 {
		check.Status = "fail"
		check.Severity = "high"
		check.Issues = append(check.Issues, "Sensitive data exposed: "+strings.Join(foundVars, ", "))
		printerln.PrintCheckResult("Environment", "FAIL", check.Duration, "Found: "+strings.Join(foundVars, ", "))
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Environment", "PASS", check.Duration, "No sensitive data exposed")
	}

	return check
}

func findSensitiveEnv(envOutput string) []string {
	lines := strings.Split(envOutput, "\n")
	found := []string{}

	for _, line := range lines {
		upperLine := strings.ToUpper(line)
		for _, key := range sensitiveKeys {
			if strings.Contains(upperLine, key) && !contains(found, key) {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) > 0 {
					found = append(found, parts[0])
				}
			}
		}
	}
	return found
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
