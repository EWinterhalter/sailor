package docker

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func StartContainer(image string) (string, error) {
	cmd := exec.Command("docker", "run", "-d", "--rm", image, "sleep", "infinity")

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("docker run failed: %w", err)
	}

	containerID := strings.TrimSpace(string(out))
	return containerID, nil
}

func Exec(containerID string, command []string) (string, error) {
	args := append([]string{"exec", containerID}, command...)
	cmd := exec.Command("docker", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func StopContainer(containerID string) error {
	cmd := exec.Command("docker", "stop", containerID)
	_ = cmd.Run()
	return nil
}
