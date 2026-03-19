package main

import (
	"github.com/alecthomas/kong"
	"github.com/slashpai/git-profile/internal/config"
)

var version = "dev"

type CLI struct {
	Config  string `help:"Path to config file." default:"${config_path}" type:"path"`
	Version kong.VersionFlag `help:"Show version." short:"v"`

	Add    AddCmd    `cmd:"" help:"Add a new git profile."`
	Remove RemoveCmd `cmd:"" aliases:"rm" help:"Remove a saved git profile."`
	List   ListCmd   `cmd:"" aliases:"ls" help:"List all saved git profiles."`
	Use    UseCmd    `cmd:"" help:"Apply a git profile to the current repo or globally."`
	Show   ShowCmd   `cmd:"" help:"Show the active git identity."`
}

type Context struct {
	ConfigPath string
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("git-profile"),
		kong.Description("Manage multiple git profiles easily."),
		kong.Vars{
			"config_path": config.DefaultConfigPath(),
			"version":     version,
		},
		kong.UsageOnError(),
	)
	err := ctx.Run(&Context{ConfigPath: cli.Config})
	ctx.FatalIfErrorf(err)
}
