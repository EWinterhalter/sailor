package checks

import (
	"strconv"
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckConnections(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Network Connections")

	check := models.CheckResult{
		Name:        "Network Connections",
		Description: "Checking active network connections",
		Issues:      []string{},
	}

	output, err := docker.Exec(containerID, []string{"sh", "-c", "ss -tunap 2>/dev/null || netstat -tunap 2>/dev/null || echo 'no tools'"})
	check.Output = output
	check.Duration = time.Since(start)

	if err != nil || strings.Contains(output, "no tools") {
		check.Status = "warn"
		check.Severity = "low"
		printerln.PrintCheckResult("Connections", "WARN", check.Duration, "Network tools unavailable")
		return check
	}

	established := strings.Count(output, "ESTAB")
	listening := strings.Count(output, "LISTEN")

	if established > 10 {
		check.Status = "warn"
		check.Severity = "medium"
		check.Issues = append(check.Issues, "High number of established connections: "+strconv.Itoa(established))
		printerln.PrintCheckResult("Connections", "WARN", check.Duration,
			strconv.Itoa(established)+" established, "+strconv.Itoa(listening)+" listening")
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Connections", "PASS", check.Duration,
			strconv.Itoa(established)+" established, "+strconv.Itoa(listening)+" listening")
	}

	return check
}
