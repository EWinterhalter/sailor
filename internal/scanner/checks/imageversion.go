package checks

import (
	"fmt"
	"strings"
	"time"

	"github.com/EWinterhalter/sailor/internal/docker"
	"github.com/EWinterhalter/sailor/internal/scanner/models"
	"github.com/EWinterhalter/sailor/internal/scanner/printerln"
)

func CheckImageVersion(containerID string) models.CheckResult {
	start := time.Now()
	printerln.PrintCheckStart("Image Version Check")

	check := models.CheckResult{
		Name:        "Image Version",
		Description: "Checking if container image is up-to-date",
		Issues:      []string{},
	}

	imageInfo, err := docker.GetImageInfo(containerID)
	check.Output = imageInfo
	check.Duration = time.Since(start)

	if err != nil {
		check.Status = "warn"
		check.Severity = "medium"
		check.Issues = append(check.Issues, "Unable to get image info")
		printerln.PrintCheckResult("Image Version", "WARN", check.Duration, "Cannot retrieve image info")
		return check
	}

	if strings.Contains(imageInfo, "latest") || strings.Contains(imageInfo, "deprecated") {
		check.Status = "warn"
		check.Severity = "medium"
		check.Issues = append(check.Issues, fmt.Sprintf("Using outdated or latest image tag: %s", imageInfo))
		printerln.PrintCheckResult("Image Version", "WARN", check.Duration, fmt.Sprintf("Image tag: %s", imageInfo))
	} else {
		check.Status = "pass"
		check.Severity = "low"
		printerln.PrintCheckResult("Image Version", "PASS", check.Duration, fmt.Sprintf("Image tag: %s", imageInfo))
	}

	return check
}
