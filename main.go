package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/slashpai/git-profile/internal/config"
)

var version = "dev"

type CLI struct {
	Config  string `help:"Path to config file." default:"${config_path}" type:"path"`
	Version kong.VersionFlag `help:"Show version." short:"v"`

	Add    AddCmd    `cmd:"" help:"Add a new git profile."`
	Update UpdateCmd `cmd:"" help:"Update an existing git profile."`
	Remove RemoveCmd `cmd:"" aliases:"rm" help:"Remove a saved git profile."`
	List   ListCmd   `cmd:"" aliases:"ls" help:"List all saved git profiles."`
	Use    UseCmd    `cmd:"" help:"Apply a git profile to the current repo or globally."`
	Show   ShowCmd   `cmd:"" help:"Show the active git identity."`
	Verify VerifyCmd `cmd:"" help:"Check if current git identity matches a saved profile."`
}

type Context struct {
	ConfigPath string
}

func main() {
	defaultCfgPath, err := config.DefaultConfigPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("git-profile"),
		kong.Description("Manage multiple git profiles easily."),
		kong.Vars{
			"config_path": defaultCfgPath,
			"version":     version,
		},
		kong.UsageOnError(),
	)
	cfgPath, err := config.ValidateConfigPath(cli.Config)
	if err != nil {
		ctx.FatalIfErrorf(err)
	}
	err = ctx.Run(&Context{ConfigPath: cfgPath})
	ctx.FatalIfErrorf(err)
}
