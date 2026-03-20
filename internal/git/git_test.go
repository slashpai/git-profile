package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// initTempRepo creates a temporary git repo and changes into it.
// Returns a cleanup function that restores the original working directory.
func initTempRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_CONFIG_NOSYSTEM=1",
			"HOME="+dir,
			"XDG_CONFIG_HOME="+filepath.Join(dir, ".config"),
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %s: %v", args, out, err)
		}
	}

	run("init")
	run("config", "user.name", "Test")
	run("config", "user.email", "test@test.com")

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Chdir(origDir)
	})

	return dir
}

func TestSetAndGetConfig(t *testing.T) {
	initTempRepo(t)

	if err := SetConfig(ScopeLocal, "user.name", "Alice"); err != nil {
		t.Fatalf("SetConfig() error: %v", err)
	}

	got, err := GetConfig("user.name")
	if err != nil {
		t.Fatalf("GetConfig() error: %v", err)
	}
	if got != "Alice" {
		t.Errorf("GetConfig(user.name) = %q, want %q", got, "Alice")
	}
}


func TestUnsetConfig(t *testing.T) {
	initTempRepo(t)

	if err := SetConfig(ScopeLocal, "test.customkey", "TESTVAL"); err != nil {
		t.Fatal(err)
	}

	if err := UnsetConfig(ScopeLocal, "test.customkey"); err != nil {
		t.Fatalf("UnsetConfig() error: %v", err)
	}

	_, err := GetConfig("test.customkey")
	if err == nil {
		t.Error("expected error after unsetting key, got nil")
	}
}

func TestGetConfigNotSet(t *testing.T) {
	initTempRepo(t)

	_, err := GetConfig("test.nonexistent")
	if err == nil {
		t.Error("expected error for unset config key, got nil")
	}
}

func TestSetConfigOutsideRepo(t *testing.T) {
	origDir, _ := os.Getwd()
	dir := t.TempDir()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(origDir) })

	err := SetConfig(ScopeLocal, "user.name", "Test")
	if err == nil {
		t.Error("expected error when setting local config outside a git repo")
	}
}
