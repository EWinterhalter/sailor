package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/EWinterhalter/sailor/internal/colors"
	"github.com/EWinterhalter/sailor/internal/db"
	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/report"
	"github.com/EWinterhalter/sailor/internal/scanner"
	"github.com/EWinterhalter/sailor/internal/scanner/models"

	"github.com/spf13/cobra"
)

var (
	flagSavePath   string
	flagSaveDB     bool
	flagDBHost     string
	flagDBPort     int
	flagDBUser     string
	flagDBPassword string
	flagDBName     string
)

func init() {
	scanCmd.Flags().StringVar(&flagSavePath, "save-result", "", "path to save result JSON")
	scanCmd.Flags().BoolVar(&flagSaveDB, "save-db", false, "save results to PostgreSQL database")
	scanCmd.Flags().StringVar(&flagDBHost, "db-host", "localhost", "database host")
	scanCmd.Flags().IntVar(&flagDBPort, "db-port", 5432, "database port")
	scanCmd.Flags().StringVar(&flagDBUser, "db-user", "sailor", "database user")
	scanCmd.Flags().StringVar(&flagDBPassword, "db-password", "sailor_pass", "database password")
	scanCmd.Flags().StringVar(&flagDBName, "db-name", "sailor_db", "database name")

	rootCmd.AddCommand(scanCmd)
}

type scanTarget struct {
	Image       string
	ContainerID string
	StartErr    error
}

var scanCmd = &cobra.Command{
	Use:  "scan <image> [image...]",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		images := args
		timeout := 60 * time.Second

		targets := make([]scanTarget, len(images))
		var wg sync.WaitGroup

		fmt.Printf("%s[INFO]%s Starting %d container(s)...\n",
			colors.ColorBlue, colors.ColorReset, len(images))

		for i, image := range images {
			i := i
			image := image

			targets[i] = scanTarget{Image: image}

			wg.Add(1)
			go func() {
				defer wg.Done()

				containerID, err := docker.StartContainer(image)
				if err != nil {
					targets[i].StartErr = err
					return
				}

				targets[i].ContainerID = containerID
			}()
		}

		wg.Wait()

		startedTargets := make([]scanTarget, 0, len(targets))
		for _, t := range targets {
			if t.StartErr != nil {
				fmt.Printf("%s[ERROR]%s Failed to start container for image %s: %v\n",
					colors.ColorRed, colors.ColorReset, t.Image, t.StartErr)
				continue
			}

			fmt.Printf("%s[INFO]%s Container started for %s: %s%s%s\n",
				colors.ColorGreen, colors.ColorReset,
				t.Image,
				colors.ColorBold, t.ContainerID[:12], colors.ColorReset)

			startedTargets = append(startedTargets, t)
		}

		if len(startedTargets) == 0 {
			return fmt.Errorf("%s[ERROR]%s failed to start any containers", colors.ColorRed, colors.ColorReset)
		}

		time.Sleep(2 * time.Second)

		hasIssues := false

		for _, t := range startedTargets {
			fmt.Printf("\n%s[INFO]%s Scanning image: %s%s%s\n",
				colors.ColorBlue, colors.ColorReset,
				colors.ColorBold, t.Image, colors.ColorReset)

			results, err := scanner.RunChecks(t.ContainerID, timeout)
			if err != nil {
				fmt.Printf("%s[WARN]%s Scan error for %s: %v\n",
					colors.ColorYellow, colors.ColorReset, t.Image, err)
				hasIssues = true
				continue
			}

			reportData := report.BuildReport(t.Image, t.ContainerID, results)

			if flagSavePath != "" {
				filePath := flagSavePath
				if len(startedTargets) > 1 {
					filePath = fmt.Sprintf("%s_%s.json", flagSavePath, t.ContainerID[:12])
				}

				jsonBytes, err := json.MarshalIndent(reportData, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal report for %s: %w", t.Image, err)
				}

				err = os.WriteFile(filePath, jsonBytes, 0644)
				if err != nil {
					return fmt.Errorf("failed to save report for %s: %w", t.Image, err)
				}

				fmt.Printf("%s[INFO]%s Report saved to: %s\n",
					colors.ColorGreen, colors.ColorReset, filePath)
			}

			if flagSaveDB {
				dbConn, err := db.Connect(db.Config{
					Host:     flagDBHost,
					Port:     flagDBPort,
					User:     flagDBUser,
					Password: flagDBPassword,
					DBName:   flagDBName,
					SSLMode:  "disable",
				})
				if err != nil {
					fmt.Printf("%s[WARN]%s Failed to connect to database: %v\n",
						colors.ColorYellow, colors.ColorReset, err)
				} else {
					repo := db.NewRepository(dbConn)
					scanID, err := repo.SaveScanReport(reportData)
					if err != nil {
						fmt.Printf("%s[WARN]%s Failed to save to database: %v\n",
							colors.ColorYellow, colors.ColorReset, err)
					} else {
						fmt.Printf("%s[INFO]%s Report saved to database with ID: %d\n",
							colors.ColorGreen, colors.ColorReset, scanID)
					}
					_ = dbConn.Close()
				}
			}

			if results.HasIssues {
				hasIssues = true
			}
		}

		for _, t := range startedTargets {
			fmt.Printf("%s[INFO]%s Stopping and removing container: %s%s%s\n",
				colors.ColorBlue, colors.ColorReset,
				colors.ColorBold, t.ContainerID[:12], colors.ColorReset)

			_ = docker.StopContainer(t.ContainerID)

			fmt.Printf("%s[INFO]%s Container cleaned up: %s%s%s\n",
				colors.ColorGreen, colors.ColorReset,
				colors.ColorBold, t.ContainerID[:12], colors.ColorReset)
		}

		if hasIssues {
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
