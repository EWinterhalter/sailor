package checks

import (
	"strings"
	"testing"
)

func TestCleanNetstatHeader_NoHeader(t *testing.T) {
	input := "tcp  ESTAB  0.0.0.0:80"
	result := cleanNetstatHeader(input)
	if result != input {
		t.Errorf("expected unchanged output, got %q", result)
	}
}

func TestCleanNetstatHeader_WithProto(t *testing.T) {
	input := "Proto Recv-Q Send-Q Local Address\ntcp  LISTEN 0.0.0.0:80\ntcp  ESTAB 0.0.0.0:443"
	result := cleanNetstatHeader(input)
	if strings.Contains(result, "Proto") {
		t.Error("Proto header line should be stripped")
	}
	if !strings.Contains(result, "tcp  LISTEN") {
		t.Error("data lines should be preserved")
	}
	if !strings.Contains(result, "ESTAB") {
		t.Error("all data lines should be preserved")
	}
}

func TestCleanNetstatHeader_EmptyInput(t *testing.T) {
	result := cleanNetstatHeader("")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestCleanNetstatHeader_OnlyHeader(t *testing.T) {
	input := "Proto Recv-Q Send-Q Local Address"
	result := cleanNetstatHeader(input)
	if strings.Contains(result, "Proto") {
		t.Error("Proto header should be stripped")
	}
}
