package scanner

import (
	"time"

	"github.com/EWinterhalter/sailor/internal/scanner/checks"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func RunChecks(containerID string, timeout time.Duration) (*models.ScanResult, error) {
	startTime := time.Now()
	result := &models.ScanResult{
		Checks: make([]models.CheckResult, 0),
	}

	printerln.PrintHeader(containerID)

	timeoutCh := time.After(timeout)
	checksDone := make(chan bool)

	go func() {
		result.Checks = append(result.Checks, checks.CheckOpenPorts(containerID))
		result.Checks = append(result.Checks, checks.CheckConnections(containerID))
		result.Checks = append(result.Checks, checks.CheckEnvironment(containerID))
		result.Checks = append(result.Checks, checks.CheckRootUser(containerID))
		result.Checks = append(result.Checks, checks.CheckImageHistory(containerID))
		result.Checks = append(result.Checks, checks.CheckWritableFS(containerID))
		checksDone <- true
	}()

	select {
	case <-timeoutCh:
		printerln.PrintTimeout()
		result.HasIssues = true
	case <-checksDone:
	}

	result.TotalTime = time.Since(startTime)
	result.CalculateSummary()
	printerln.PrintSummary(result)

	return result, nil
}
