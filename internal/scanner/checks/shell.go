package checks

import (
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckRemoteShell(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Remote Shell Check")

	check := models.CheckResult{
		Name:        "Remote Shell",
		Description: "Checks for netcat spawning a shell",
		Issues:      []string{},
	}

	output, _ := docker.Exec(containerID, []string{
		"sh", "-c", "ps aux | grep -E 'nc.*(/bin/sh|/bin/bash)' | grep -v grep",
	})

	check.Output = output
	check.Duration = time.Since(start)

	if strings.TrimSpace(output) != "" {
		check.Status = "fail"
		check.Severity = "critical"
		check.Issues = append(check.Issues, "Remote shell detected via netcat")
		printerln.PrintCheckResult("Remote Shell", "FAIL", check.Duration, "shell exposed")
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Remote Shell", "PASS", check.Duration, "no shell")
	}

	return check
}
