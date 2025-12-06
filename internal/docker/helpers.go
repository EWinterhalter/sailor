package docker

import (
	"os/exec"
	"strings"
)

func GetImageInfo(containerID string) (string, error) {
	cmd := exec.Command("docker", "inspect", "--format", "{{.Image}}", containerID)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	imageID := strings.TrimSpace(string(out))

	cmd = exec.Command("docker", "images", "--no-trunc", "--format", "{{.Repository}}:{{.Tag}} {{.ID}}")
	out, err = cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[1] == imageID {
			return parts[0], nil
		}
	}

	return imageID, nil
}

func ImageHistory(containerID string) (string, error) {
	imageID, err := getImageID(containerID)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("docker", "history", "--no-trunc", "--format", "{{.CreatedBy}}", imageID)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func getImageID(containerID string) (string, error) {
	cmd := exec.Command("docker", "inspect", "--format", "{{.Image}}", containerID)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
