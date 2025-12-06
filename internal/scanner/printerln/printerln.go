package printerln

import (
	"fmt"
	"time"

	"github.com/EWinterhalter/sailor/internal/colors"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
)

func PrintHeader(containerID string) {
	fmt.Printf("%s🐳 Container ID:%s %s\n", colors.ColorBold, colors.ColorReset, containerID[:12])
	fmt.Printf("%s⏰ Scan started:%s %s\n", colors.ColorBold, colors.ColorReset, time.Now().Format("2006-01-02 15:04:05"))
}

func PrintCheckStart(name string) {
	fmt.Printf("%s▶ Running check:%s %s%s%s\n", colors.ColorBlue, colors.ColorReset, colors.ColorBold, name, colors.ColorReset)
}

func PrintCheckResult(name, status string, duration time.Duration, details string) {
	var statusColor, statusSymbol string

	switch status {
	case "PASS":
		statusColor = colors.ColorGreen
		statusSymbol = "✓"
	case "WARN":
		statusColor = colors.ColorYellow
		statusSymbol = "⚠"
	case "FAIL":
		statusColor = colors.ColorRed
		statusSymbol = "✗"
	default:
		statusColor = colors.ColorWhite
		statusSymbol = "?"
	}

	fmt.Printf("  %s%s %s%s %s[%dms]%s %s%s\n",
		statusColor, statusSymbol, status, colors.ColorReset,
		colors.ColorDim, duration.Milliseconds(), colors.ColorReset,
		colors.ColorDim, details)
}

func PrintTimeout() {
	fmt.Printf("\n%s⚠ TIMEOUT:%s Scan exceeded time limit\n", colors.ColorYellow, colors.ColorReset)
}

func PrintSummary(result *models.ScanResult) {
	fmt.Printf("  %sTotal Checks:%s    %d\n", colors.ColorBold, colors.ColorReset, result.Summary.TotalChecks)
	fmt.Printf("  %s%s✓ Passed:%s       %d\n", colors.ColorGreen, colors.ColorBold, colors.ColorReset, result.Summary.Passed)
	fmt.Printf("  %s%s⚠ Warnings:%s     %d\n", colors.ColorYellow, colors.ColorBold, colors.ColorReset, result.Summary.Warnings)
	fmt.Printf("  %s%s✗ Failed:%s       %d\n", colors.ColorRed, colors.ColorBold, colors.ColorReset, result.Summary.Failed)

	fmt.Printf("\n  %sSeverity Breakdown:%s\n", colors.ColorBold, colors.ColorReset)
	if result.Summary.Critical > 0 {
		fmt.Printf("    %s● Critical:%s  %d\n", colors.ColorRed, colors.ColorReset, result.Summary.Critical)
	}
	if result.Summary.High > 0 {
		fmt.Printf("    %s● High:%s      %d\n", colors.ColorRed, colors.ColorReset, result.Summary.High)
	}
	if result.Summary.Medium > 0 {
		fmt.Printf("    %s● Medium:%s    %d\n", colors.ColorYellow, colors.ColorReset, result.Summary.Medium)
	}
	if result.Summary.Low > 0 {
		fmt.Printf("    %s● Low:%s       %d\n", colors.ColorGreen, colors.ColorReset, result.Summary.Low)
	}

	fmt.Printf("%s⏱ Total Time:%s   %s\n", colors.ColorBold, colors.ColorReset, result.TotalTime.Round(time.Millisecond))

	if result.HasIssues {
		fmt.Printf("\n  %s%s⚠ SECURITY ISSUES DETECTED ⚠%s\n", colors.ColorBold, colors.ColorRed, colors.ColorReset)
	} else {
		fmt.Printf("%s%s✓ NO CRITICAL ISSUES FOUND%s", colors.ColorBold, colors.ColorGreen, colors.ColorReset)
	}

	fmt.Println()
}
