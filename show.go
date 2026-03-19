package main

import (
	"fmt"

	"github.com/slashpai/git-profile/internal/git"
)

type ShowCmd struct{}

func (cmd *ShowCmd) Run(ctx *Context) error {
	fields := []struct {
		label string
		key   string
	}{
		{"user.name", "user.name"},
		{"user.email", "user.email"},
		{"user.signingkey", "user.signingkey"},
		{"commit.gpgsign", "commit.gpgsign"},
	}

	fmt.Println("Current git identity:")
	for _, f := range fields {
		val, err := git.GetConfig(f.key)
		if err != nil {
			val = "(not set)"
		}
		fmt.Printf("  %-16s = %s\n", f.label, val)
	}

	return nil
}
