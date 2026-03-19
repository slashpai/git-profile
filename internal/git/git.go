package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Scope string

const (
	ScopeLocal  Scope = "--local"
	ScopeGlobal Scope = "--global"
)

func SetConfig(scope Scope, key, value string) error {
	cmd := exec.Command("git", "config", string(scope), key, value)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git config %s %s: %s: %w", scope, key, strings.TrimSpace(string(out)), err)
	}
	return nil
}

func GetConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", key)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func UnsetConfig(scope Scope, key string) error {
	cmd := exec.Command("git", "config", string(scope), "--unset", key)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git config %s --unset %s: %s: %w", scope, key, strings.TrimSpace(string(out)), err)
	}
	return nil
}
