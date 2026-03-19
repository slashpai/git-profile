package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfigPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}
	want := filepath.Join(home, ".git-profiles.yaml")
	got, err := DefaultConfigPath()
	if err != nil {
		t.Fatalf("DefaultConfigPath() returned error: %v", err)
	}
	if got != want {
		t.Errorf("DefaultConfigPath() = %q, want %q", got, want)
	}
}

func TestValidateProfileName(t *testing.T) {
	valid := []string{"personal", "work", "my-org", "work_2", "github.com", "A1"}
	for _, name := range valid {
		if err := ValidateProfileName(name); err != nil {
			t.Errorf("ValidateProfileName(%q) should be valid, got: %v", name, err)
		}
	}

	invalid := []string{"", "-dash", ".dot", "has space", "new\nline", "tab\there", "special!char", "/slash"}
	for _, name := range invalid {
		if err := ValidateProfileName(name); err == nil {
			t.Errorf("ValidateProfileName(%q) should be invalid, got nil", name)
		}
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	cfg, err := Load("/tmp/does-not-exist-git-profile-test.yaml")
	if err != nil {
		t.Fatalf("Load() returned error for missing file: %v", err)
	}
	if cfg.Profiles == nil {
		t.Fatal("expected non-nil Profiles map")
	}
	if len(cfg.Profiles) != 0 {
		t.Errorf("expected empty Profiles, got %d", len(cfg.Profiles))
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.yaml")
	if err := os.WriteFile(path, []byte(":::invalid"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.yaml")

	original := &Config{
		Profiles: map[string]Profile{
			"work": {
				Name:       "Alice",
				Email:      "alice@work.com",
				SigningKey: "KEY123",
				GPGSign:    true,
				SSHKey:     "~/.ssh/id_work",
			},
			"personal": {
				Name:  "Alice",
				Email: "alice@home.com",
			},
		},
	}

	if err := Save(path, original); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(loaded.Profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(loaded.Profiles))
	}

	work := loaded.Profiles["work"]
	if work.Name != "Alice" || work.Email != "alice@work.com" {
		t.Errorf("work profile mismatch: got name=%q email=%q", work.Name, work.Email)
	}
	if work.SigningKey != "KEY123" {
		t.Errorf("work signing key = %q, want %q", work.SigningKey, "KEY123")
	}
	if !work.GPGSign {
		t.Error("work GPGSign should be true")
	}
	if work.SSHKey != "~/.ssh/id_work" {
		t.Errorf("work SSHKey = %q, want %q", work.SSHKey, "~/.ssh/id_work")
	}

	personal := loaded.Profiles["personal"]
	if personal.Name != "Alice" || personal.Email != "alice@home.com" {
		t.Errorf("personal profile mismatch: got name=%q email=%q", personal.Name, personal.Email)
	}
	if personal.SigningKey != "" {
		t.Errorf("personal signing key should be empty, got %q", personal.SigningKey)
	}
	if personal.GPGSign {
		t.Error("personal GPGSign should be false")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "dir", "profiles.yaml")

	cfg := &Config{Profiles: map[string]Profile{
		"test": {Name: "Test", Email: "test@test.com"},
	}}

	if err := Save(path, cfg); err != nil {
		t.Fatalf("Save() should create nested dirs: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}
}

func TestSaveFilePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.yaml")

	cfg := &Config{Profiles: map[string]Profile{
		"test": {Name: "Test", Email: "test@test.com"},
	}}

	if err := Save(path, cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}

	want := fs.FileMode(0o600)
	got := info.Mode().Perm()
	if got != want {
		t.Errorf("file permissions = %o, want %o", got, want)
	}
}

func TestProfileNotFoundErrorEmpty(t *testing.T) {
	cfg := &Config{Profiles: make(map[string]Profile)}
	err := cfg.ProfileNotFoundError("work")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, `"work"`) {
		t.Errorf("error should mention profile name, got: %s", msg)
	}
	if !strings.Contains(msg, "no profiles configured") {
		t.Errorf("error should mention no profiles configured, got: %s", msg)
	}
	if !strings.Contains(msg, "git-profile add") {
		t.Errorf("error should hint to use 'git-profile add', got: %s", msg)
	}
}

func TestProfileNotFoundErrorWithProfiles(t *testing.T) {
	cfg := &Config{Profiles: map[string]Profile{
		"personal": {Name: "Alice", Email: "alice@home.com"},
		"work":     {Name: "Alice", Email: "alice@work.com"},
	}}
	err := cfg.ProfileNotFoundError("typo")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, `"typo"`) {
		t.Errorf("error should mention profile name, got: %s", msg)
	}
	if !strings.Contains(msg, "personal") || !strings.Contains(msg, "work") {
		t.Errorf("error should list available profiles, got: %s", msg)
	}
	if !strings.Contains(msg, "git-profile list") {
		t.Errorf("error should hint to use 'git-profile list', got: %s", msg)
	}
}

func TestLoadEmptyProfilesMap(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.yaml")
	if err := os.WriteFile(path, []byte("profiles:\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Profiles == nil {
		t.Fatal("Profiles map should be initialized, not nil")
	}
}
