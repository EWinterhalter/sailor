package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/EWinterhalter/sailor/internal/colors"
	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/report"
	"github.com/EWinterhalter/sailor/internal/scanner"
	"github.com/EWinterhalter/sailor/internal/scanner/models"

	"github.com/spf13/cobra"
)

var (
	flagSavePath string
)

func init() {
	scanCmd.Flags().StringVar(&flagSavePath, "save-result", "", "path to save result JSON")

	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:  "scan <image>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		image := args[0]

		fmt.Printf("%s[INFO]%s Starting security scan for image: %s%s%s\n",
			colors.ColorBlue, colors.ColorReset,
			colors.ColorBold, image, colors.ColorReset)

		timeout := time.Duration(60) * time.Second

		fmt.Printf("%s[INFO]%s Starting container...\n", colors.ColorBlue, colors.ColorReset)
		containerID, err := docker.StartContainer(image)
		if err != nil {
			return fmt.Errorf("%s[ERROR]%s Failed to start container: %w", colors.ColorRed, colors.ColorReset, err)
		}
		fmt.Printf("%s[INFO]%s Container started: %s%s%s\n",
			colors.ColorGreen, colors.ColorReset,
			colors.ColorBold, containerID[:12], colors.ColorReset)

		time.Sleep(2 * time.Second)

		results, err := scanner.RunChecks(containerID, timeout)
		if err != nil {
			return fmt.Errorf("%s[ERROR]%s Scan error: %w", colors.ColorRed, colors.ColorReset, err)
		}

		reportData := report.BuildReport(image, containerID, results)

		if flagSavePath != "" {
			jsonBytes, err := json.MarshalIndent(reportData, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal report: %w", err)
			}

			err = os.WriteFile(flagSavePath, jsonBytes, 0644)
			if err != nil {
				return fmt.Errorf("failed to save report: %w", err)
			}
			fmt.Printf("%s[INFO]%s Report saved to: %s\n",
				colors.ColorGreen, colors.ColorReset, flagSavePath)
		}

		fmt.Printf("%s[INFO]%s Stopping and removing container...\n",
			colors.ColorBlue, colors.ColorReset)
		_ = docker.StopContainer(containerID)
		fmt.Printf("%s[INFO]%s Container cleaned up\n",
			colors.ColorGreen, colors.ColorReset)

		if results.HasIssues {
			fmt.Printf("\n%s[ALERT]%s Security issues detected - exiting with error code 1\n",
				colors.ColorRed, colors.ColorReset)
			os.Exit(1)
		}

		return nil
	},
}

func printDetailedResults(results *models.ScanResult) {
	for i, check := range results.Checks {
		fmt.Printf("%s[%d] %s%s\n", colors.ColorBold, i+1, check.Name, colors.ColorReset)
		fmt.Printf("    Description: %s\n", check.Description)
		fmt.Printf("    Status: %s\n", getStatusString(check.Status))
		fmt.Printf("    Severity: %s\n", getSeverityString(check.Severity))
		fmt.Printf("    Duration: %dms\n", check.Duration.Milliseconds())

		if len(check.Issues) > 0 {
			fmt.Printf("    Issues:\n")
			for _, issue := range check.Issues {
				fmt.Printf("      - %s%s%s\n", colors.ColorRed, issue, colors.ColorReset)
			}
		}

		if check.Output != "" && len(check.Output) < 500 {
			fmt.Printf("    Output preview:\n")
			fmt.Printf("%s%s%s\n", colors.ColorDim, truncateString(check.Output, 300), colors.ColorReset)
		}

		fmt.Println()
	}
}

func getStatusString(status string) string {
	switch status {
	case "pass":
		return fmt.Sprintf("%s✓ PASS%s", colors.ColorGreen, colors.ColorReset)
	case "warn":
		return fmt.Sprintf("%s⚠ WARN%s", colors.ColorYellow, colors.ColorReset)
	case "fail":
		return fmt.Sprintf("%s✗ FAIL%s", colors.ColorRed, colors.ColorReset)
	default:
		return status
	}
}

func getSeverityString(severity string) string {
	switch severity {
	case "critical":
		return fmt.Sprintf("%s● CRITICAL%s", colors.ColorRed, colors.ColorReset)
	case "high":
		return fmt.Sprintf("%s● HIGH%s", colors.ColorRed, colors.ColorReset)
	case "medium":
		return fmt.Sprintf("%s● MEDIUM%s", colors.ColorYellow, colors.ColorReset)
	case "low":
		return fmt.Sprintf("%s● LOW%s", colors.ColorGreen, colors.ColorReset)
	default:
		return severity
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
