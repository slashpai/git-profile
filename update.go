package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/slashpai/git-profile/internal/config"
)

type UpdateCmd struct {
	Name string `arg:"" help:"Name of the profile to update."`
}

func (cmd *UpdateCmd) Run(ctx *Context) error {
	cfg, err := config.Load(ctx.ConfigPath)
	if err != nil {
		return err
	}

	profile, exists := cfg.Profiles[cmd.Name]
	if !exists {
		return cfg.ProfileNotFoundError(cmd.Name)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Updating profile %q (press Enter to keep current value):\n", cmd.Name)

	profile.Name = promptRequiredWithDefault(scanner, "user.name", profile.Name)
	profile.Email = promptRequiredWithDefault(scanner, "user.email", profile.Email)
	profile.SigningKey = promptWithDefault(scanner, "user.signingkey - GPG key ID, run 'gpg --list-secret-keys --keyid-format long' to find it", profile.SigningKey)

	gpgDefault := "n"
	if profile.GPGSign {
		gpgDefault = "y"
	}
	gpg := promptWithDefault(scanner, "commit.gpgsign [y/N]", gpgDefault)
	profile.GPGSign = strings.EqualFold(gpg, "y") || strings.EqualFold(gpg, "yes")

	profile.SSHKey = promptWithDefault(scanner, "SSH key path, e.g. ~/.ssh/id_ed25519", profile.SSHKey)

	cfg.Profiles[cmd.Name] = profile

	if err := config.Save(ctx.ConfigPath, cfg); err != nil {
		return err
	}

	fmt.Printf("Profile %q updated.\n", cmd.Name)
	return nil
}

