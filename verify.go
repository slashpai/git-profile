package main

import (
	"fmt"

	"github.com/slashpai/git-profile/internal/config"
	"github.com/slashpai/git-profile/internal/git"
)

type VerifyCmd struct {
	ShowIdentity bool `help:"Also display the current git identity." short:"s"`
}

func (cmd *VerifyCmd) Run(ctx *Context) error {
	currentName, err := git.GetConfig("user.name")
	if err != nil {
		return fmt.Errorf("no git user.name configured in this repo")
	}
	currentEmail, err := git.GetConfig("user.email")
	if err != nil {
		return fmt.Errorf("no git user.email configured in this repo")
	}

	cfg, err := config.Load(ctx.ConfigPath)
	if err != nil {
		return err
	}

	if len(cfg.Profiles) == 0 {
		fmt.Printf("Current identity: %s <%s>\n", currentName, currentEmail)
		fmt.Println("No saved profiles to compare against (use 'git-profile add' to create one).")
		return nil
	}

	matched := cfg.MatchingProfiles(currentName, currentEmail)

	if len(matched) == 0 {
		fmt.Printf("Warning: current identity does not match any saved profile.\n")
		fmt.Printf("  user.name  = %s\n", currentName)
		fmt.Printf("  user.email = %s\n", currentEmail)
		fmt.Println("\nUse 'git-profile list' to see saved profiles or 'git-profile use <name>' to switch.")
		return nil
	}

	for _, name := range matched {
		fmt.Printf("Current identity matches profile %q.\n", name)
	}

	if cmd.ShowIdentity {
		fmt.Println("\nCurrent git identity:")
		fields := []struct{ label, key string }{
			{"user.name", "user.name"},
			{"user.email", "user.email"},
			{"user.signingkey", "user.signingkey"},
			{"commit.gpgsign", "commit.gpgsign"},
		}
		for _, f := range fields {
			val, err := git.GetConfig(f.key)
			if err != nil {
				val = "(not set)"
			}
			fmt.Printf("  %-16s = %s\n", f.label, val)
		}
	}

	return nil
}
