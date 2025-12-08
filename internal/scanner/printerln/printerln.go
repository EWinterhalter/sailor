package printerln

import (
	"fmt"
	"sync"
	"time"

	"github.com/EWinterhalter/sailor/internal/colors"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
)

type Spinner struct {
	frames   []string
	interval time.Duration
	stop     chan bool
	done     chan bool
	mu       sync.Mutex
	active   bool
}

var activeSpinner *Spinner

func newSpinner() *Spinner {
	return &Spinner{
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		interval: 80 * time.Millisecond,
		stop:     make(chan bool),
		done:     make(chan bool),
	}
}

func (s *Spinner) start(name string) {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				s.done <- true
				return
			default:
				fmt.Printf("\r  %s%s%s %s%s%s",
					colors.ColorBlue, s.frames[i], colors.ColorReset,
					colors.ColorBold, name, colors.ColorReset)
				i = (i + 1) % len(s.frames)
				time.Sleep(s.interval)
			}
		}
	}()
}

func (s *Spinner) stopSpinner() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	s.stop <- true
	<-s.done

	s.mu.Lock()
	s.active = false
	s.mu.Unlock()
}

func PrintHeader(containerID string) {
	fmt.Printf("%s[🐳] Container ID:%s %s\n", colors.ColorBold, colors.ColorReset, containerID[:12])
	fmt.Printf("%s[⏰] Scan started:%s %s\n\n", colors.ColorBold, colors.ColorReset, time.Now().Format("2006-01-02 15:04:05"))
}

func PrintCheckStart(name string) {
	activeSpinner = newSpinner()
	activeSpinner.start(name)
}

func PrintCheckResult(name, status string, duration time.Duration, details string) {
	if activeSpinner != nil {
		activeSpinner.stopSpinner()
	}

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

	fmt.Printf("\r%s%s %s%s %s%s%s %s[%dms]%s",
		statusColor, statusSymbol, status, colors.ColorReset,
		colors.ColorBold, name, colors.ColorReset,
		colors.ColorDim, duration.Milliseconds(), colors.ColorReset)

	if details != "" {
		fmt.Printf(" %s", details)
	}
	fmt.Println()

	activeSpinner = nil
}

func PrintTimeout() {
	if activeSpinner != nil {
		activeSpinner.stopSpinner()
		fmt.Print("\r")
		activeSpinner = nil
	}
	fmt.Printf("\n%s⚠ TIMEOUT:%s Scan exceeded time limit\n", colors.ColorYellow, colors.ColorReset)
}

func PrintSummary(result *models.ScanResult) {
	fmt.Printf("\n%sTotal Checks:%s    %d\n", colors.ColorBold, colors.ColorReset, result.Summary.TotalChecks)
	fmt.Printf("%s%s✓ Passed:%s       %d\n", colors.ColorGreen, colors.ColorBold, colors.ColorReset, result.Summary.Passed)
	fmt.Printf("%s%s⚠ Warnings:%s     %d\n", colors.ColorYellow, colors.ColorBold, colors.ColorReset, result.Summary.Warnings)
	fmt.Printf("%s%s✗ Failed:%s       %d\n", colors.ColorRed, colors.ColorBold, colors.ColorReset, result.Summary.Failed)

	fmt.Printf("\n%sSeverity Breakdown:%s\n", colors.ColorBold, colors.ColorReset)
	if result.Summary.Critical > 0 {
		fmt.Printf("%s● Critical:%s  %d\n", colors.ColorRed, colors.ColorReset, result.Summary.Critical)
	}
	if result.Summary.High > 0 {
		fmt.Printf("%s● High:%s      %d\n", colors.ColorRed, colors.ColorReset, result.Summary.High)
	}
	if result.Summary.Medium > 0 {
		fmt.Printf("%s● Medium:%s    %d\n", colors.ColorYellow, colors.ColorReset, result.Summary.Medium)
	}
	if result.Summary.Low > 0 {
		fmt.Printf("%s● Low:%s       %d\n", colors.ColorGreen, colors.ColorReset, result.Summary.Low)
	}

	fmt.Printf("%s⏱ Total Time:%s   %s\n", colors.ColorBold, colors.ColorReset, result.TotalTime.Round(time.Millisecond))

	if result.HasIssues {
		fmt.Printf("\n%s%s⚠ SECURITY ISSUES DETECTED ⚠%s\n", colors.ColorBold, colors.ColorRed, colors.ColorReset)
	} else {
		fmt.Printf("%s%s✓ NO CRITICAL ISSUES FOUND%s", colors.ColorBold, colors.ColorGreen, colors.ColorReset)
	}

	fmt.Println()
}
