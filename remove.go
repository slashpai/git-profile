package main

import (
	"fmt"

	"github.com/slashpai/git-profile/internal/config"
)

type RemoveCmd struct {
	Name string `arg:"" help:"Name of the profile to remove."`
}

func (cmd *RemoveCmd) Run(ctx *Context) error {
	cfg, err := config.Load(ctx.ConfigPath)
	if err != nil {
		return err
	}

	if _, exists := cfg.Profiles[cmd.Name]; !exists {
		return cfg.ProfileNotFoundError(cmd.Name)
	}

	delete(cfg.Profiles, cmd.Name)

	if err := config.Save(ctx.ConfigPath, cfg); err != nil {
		return err
	}

	fmt.Printf("Profile %q removed.\n", cmd.Name)
	return nil
}
