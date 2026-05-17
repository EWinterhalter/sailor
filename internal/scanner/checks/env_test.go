package checks

import "testing"

func TestFindSensitiveEnv_Empty(t *testing.T) {
	found := findSensitiveEnv("")
	if len(found) != 0 {
		t.Errorf("expected empty, got %v", found)
	}
}

func TestFindSensitiveEnv_NoSensitive(t *testing.T) {
	env := "PATH=/usr/bin:/bin\nHOME=/root\nTERM=xterm"
	found := findSensitiveEnv(env)
	if len(found) != 0 {
		t.Errorf("expected no sensitive vars, got %v", found)
	}
}

func TestFindSensitiveEnv_Password(t *testing.T) {
	env := "PATH=/usr/bin\nDB_PASSWORD=supersecret\nHOME=/root"
	found := findSensitiveEnv(env)
	if len(found) != 1 {
		t.Errorf("expected 1 sensitive var, got %v", found)
	}
	if found[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %s", found[0])
	}
}

func TestFindSensitiveEnv_MultipleKeys(t *testing.T) {
	env := "API_KEY=abc123\nSECRET=xyz\nPATH=/bin"
	found := findSensitiveEnv(env)
	if len(found) != 2 {
		t.Errorf("expected 2 sensitive vars, got %v", found)
	}
}

func TestFindSensitiveEnv_CaseInsensitive(t *testing.T) {
	env := "aws_secret_access_key=AKIAIOSFODNN7EXAMPLE"
	found := findSensitiveEnv(env)
	if len(found) != 1 {
		t.Errorf("expected 1 sensitive var (case-insensitive match), got %v", found)
	}
}

func TestFindSensitiveEnv_NoDuplicates(t *testing.T) {
	env := "PASSWORD_SECRET=value"
	found := findSensitiveEnv(env)
	seen := map[string]int{}
	for _, v := range found {
		seen[v]++
		if seen[v] > 1 {
			t.Errorf("duplicate entry %q in results", v)
		}
	}
}

func TestContains_Found(t *testing.T) {
	if !contains([]string{"a", "b", "c"}, "b") {
		t.Error("expected true")
	}
}

func TestContains_NotFound(t *testing.T) {
	if contains([]string{"a", "b", "c"}, "z") {
		t.Error("expected false")
	}
}

func TestContains_Empty(t *testing.T) {
	if contains([]string{}, "a") {
		t.Error("expected false for empty slice")
	}
}
