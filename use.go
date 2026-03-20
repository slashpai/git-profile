package main

import (
	"fmt"

	"github.com/slashpai/git-profile/internal/config"
	"github.com/slashpai/git-profile/internal/git"
)

type UseCmd struct {
	Name   string `arg:"" help:"Name of the profile to apply."`
	Global bool   `help:"Apply globally instead of to the current repo." default:"false"`
}

func (cmd *UseCmd) Run(ctx *Context) error {
	cfg, err := config.Load(ctx.ConfigPath)
	if err != nil {
		return err
	}

	profile, exists := cfg.Profiles[cmd.Name]
	if !exists {
		return cfg.ProfileNotFoundError(cmd.Name)
	}

	scope := git.ScopeLocal
	scopeLabel := "local"
	if cmd.Global {
		scope = git.ScopeGlobal
		scopeLabel = "global"
	}

	if err := git.SetConfig(scope, "user.name", profile.Name); err != nil {
		return err
	}
	if err := git.SetConfig(scope, "user.email", profile.Email); err != nil {
		return err
	}

	if profile.SigningKey != "" {
		if err := git.SetConfig(scope, "user.signingkey", profile.SigningKey); err != nil {
			return err
		}
	} else {
		_ = git.UnsetConfig(scope, "user.signingkey")
	}

	if profile.GPGSign {
		if err := git.SetConfig(scope, "commit.gpgsign", "true"); err != nil {
			return err
		}
	} else {
		_ = git.UnsetConfig(scope, "commit.gpgsign")
	}

	fmt.Printf("Switched to profile %q (%s).\n", cmd.Name, scopeLabel)
	fmt.Printf("  user.name       = %s\n", profile.Name)
	fmt.Printf("  user.email      = %s\n", profile.Email)
	if profile.SigningKey != "" {
		fmt.Printf("  user.signingkey = %s\n", profile.SigningKey)
	}
	fmt.Printf("  commit.gpgsign  = %v\n", profile.GPGSign)

	return nil
}
