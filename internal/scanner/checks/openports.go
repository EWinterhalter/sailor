package checks

import (
	"fmt"
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckOpenPorts(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Open Ports Analysis")

	check := models.CheckResult{
		Name:        "Open Ports",
		Description: "Checking for exposed network ports",
		Issues:      []string{},
	}

	output, err := docker.Exec(containerID, []string{
		"sh", "-c",
		"netstat -tulpn 2>/dev/null || ss -tulpn 2>/dev/null || echo 'no network tools'",
	})

	cleaned := cleanPortScanHeader(output)
	check.Output = cleaned

	check.Duration = time.Since(start)

	if err != nil || strings.Contains(output, "no network tools") {
		check.Status = "warn"
		check.Severity = "low"
		check.Issues = append(check.Issues, "Unable to check ports - network tools not available")
		printerln.PrintCheckResult("Open Ports", "WARN", check.Duration, "Network tools unavailable")
		return check
	}

	openPorts, suspiciousPorts := parseOpenPorts(cleaned)

	if len(suspiciousPorts) > 0 {
		check.Status = "fail"
		check.Severity = "high"
		check.Issues = append(check.Issues, fmt.Sprintf("Suspicious ports detected: %v", suspiciousPorts))
		printerln.PrintCheckResult(
			"Open Ports",
			"FAIL",
			check.Duration,
			fmt.Sprintf("%d open ports, %d suspicious", openPorts, len(suspiciousPorts)),
		)
	} else if openPorts > 5 {
		check.Status = "warn"
		check.Severity = "medium"
		check.Issues = append(check.Issues, fmt.Sprintf("Many open ports detected: %d", openPorts))
		printerln.PrintCheckResult("Open Ports", "WARN", check.Duration, fmt.Sprintf("%d open ports", openPorts))
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Open Ports", "PASS", check.Duration, fmt.Sprintf("%d open ports", openPorts))
	}

	return check
}

func parseOpenPorts(output string) (int, []string) {
	lines := strings.Split(output, "\n")
	openPorts := 0
	suspiciousPorts := []string{}
	dangerPorts := []string{"4444", "31337", "12345", "6667", "6666"}

	for _, line := range lines {
		if strings.Contains(line, "LISTEN") || strings.Contains(line, "UNCONN") {
			openPorts++
			for _, port := range dangerPorts {
				if strings.Contains(line, ":"+port) {
					suspiciousPorts = append(suspiciousPorts, port)
				}
			}
		}
	}

	return openPorts, suspiciousPorts
}

func cleanPortScanHeader(s string) string {
	lines := strings.Split(s, "\n")

	startIdx := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Proto") || strings.HasPrefix(trimmed, "Netid") {
			startIdx = i + 1
			break
		}
	}

	if startIdx < len(lines) {
		return strings.Join(lines[startIdx:], "\n")
	}

	return s
}
