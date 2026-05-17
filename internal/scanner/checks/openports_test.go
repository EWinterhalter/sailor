package checks

import (
	"strings"
	"testing"
)

func TestParseOpenPorts_Empty(t *testing.T) {
	count, suspicious := parseOpenPorts("")
	if count != 0 {
		t.Errorf("expected 0 ports, got %d", count)
	}
	if len(suspicious) != 0 {
		t.Errorf("expected no suspicious ports, got %v", suspicious)
	}
}

func TestParseOpenPorts_NoListen(t *testing.T) {
	output := "tcp  ESTAB  0.0.0.0:80  0.0.0.0:*\ntcp  ESTAB  0.0.0.0:443  0.0.0.0:*"
	count, suspicious := parseOpenPorts(output)
	if count != 0 {
		t.Errorf("expected 0 listening ports, got %d", count)
	}
	if len(suspicious) != 0 {
		t.Errorf("expected no suspicious ports, got %v", suspicious)
	}
}

func TestParseOpenPorts_NormalPorts(t *testing.T) {
	output := "tcp  LISTEN  0.0.0.0:80  0.0.0.0:*\ntcp  LISTEN  0.0.0.0:443  0.0.0.0:*"
	count, suspicious := parseOpenPorts(output)
	if count != 2 {
		t.Errorf("expected 2 listening ports, got %d", count)
	}
	if len(suspicious) != 0 {
		t.Errorf("expected no suspicious ports, got %v", suspicious)
	}
}

func TestParseOpenPorts_SuspiciousPort4444(t *testing.T) {
	output := "tcp  LISTEN  0.0.0.0:4444  0.0.0.0:*"
	count, suspicious := parseOpenPorts(output)
	if count != 1 {
		t.Errorf("expected 1 port, got %d", count)
	}
	if len(suspicious) != 1 || suspicious[0] != "4444" {
		t.Errorf("expected [4444] suspicious, got %v", suspicious)
	}
}

func TestParseOpenPorts_MultipleSniffPorts(t *testing.T) {
	output := strings.Join([]string{
		"tcp  LISTEN  0.0.0.0:31337  0.0.0.0:*",
		"tcp  LISTEN  0.0.0.0:6667   0.0.0.0:*",
		"tcp  LISTEN  0.0.0.0:80     0.0.0.0:*",
	}, "\n")
	count, suspicious := parseOpenPorts(output)
	if count != 3 {
		t.Errorf("expected 3 ports, got %d", count)
	}
	if len(suspicious) != 2 {
		t.Errorf("expected 2 suspicious ports, got %v", suspicious)
	}
}

func TestParseOpenPorts_UNCONNCounted(t *testing.T) {
	output := "udp  UNCONN  0.0.0.0:53  0.0.0.0:*"
	count, _ := parseOpenPorts(output)
	if count != 1 {
		t.Errorf("expected 1 UNCONN port, got %d", count)
	}
}

func TestCleanPortScanHeader_NoHeader(t *testing.T) {
	input := "data line 1\ndata line 2"
	result := cleanPortScanHeader(input)
	if result != input {
		t.Errorf("expected unchanged output, got %q", result)
	}
}

func TestCleanPortScanHeader_ProtoHeader(t *testing.T) {
	input := "Proto Recv-Q Send-Q Local\ntcp  LISTEN 0.0.0.0:80"
	result := cleanPortScanHeader(input)
	if strings.Contains(result, "Proto") {
		t.Error("header line should be stripped")
	}
	if !strings.Contains(result, "tcp  LISTEN") {
		t.Error("data lines should be preserved")
	}
}

func TestCleanPortScanHeader_NetidHeader(t *testing.T) {
	input := "Netid State Local\ntcp   LISTEN 0.0.0.0:443"
	result := cleanPortScanHeader(input)
	if strings.Contains(result, "Netid") {
		t.Error("Netid header should be stripped")
	}
	if !strings.Contains(result, "tcp") {
		t.Error("data lines should be preserved")
	}
}
