package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/slashpai/git-profile/internal/config"
)

type AddCmd struct {
	Name string `arg:"" help:"Name for the new profile."`
}

func (cmd *AddCmd) Run(ctx *Context) error {
	cfg, err := config.Load(ctx.ConfigPath)
	if err != nil {
		return err
	}

	if _, exists := cfg.Profiles[cmd.Name]; exists {
		return fmt.Errorf("profile %q already exists", cmd.Name)
	}

	scanner := bufio.NewScanner(os.Stdin)
	profile := config.Profile{}

	profile.Name = prompt(scanner, "user.name")
	profile.Email = prompt(scanner, "user.email")
	profile.SigningKey = prompt(scanner, "user.signingkey - GPG key ID, run 'gpg --list-secret-keys --keyid-format long' to find it (optional, Enter to skip)")

	gpg := prompt(scanner, "commit.gpgsign [y/N]")
	profile.GPGSign = strings.EqualFold(gpg, "y") || strings.EqualFold(gpg, "yes")

	profile.SSHKey = prompt(scanner, "SSH key path, e.g. ~/.ssh/id_ed25519 (optional, Enter to skip)")

	cfg.Profiles[cmd.Name] = profile

	if err := config.Save(ctx.ConfigPath, cfg); err != nil {
		return err
	}

	fmt.Printf("Profile %q saved.\n", cmd.Name)
	return nil
}

func prompt(scanner *bufio.Scanner, label string) string {
	fmt.Printf("  %s: ", label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
