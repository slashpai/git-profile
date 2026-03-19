package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/slashpai/git-profile/internal/config"
)

type ListCmd struct{}

func (cmd *ListCmd) Run(ctx *Context) error {
	cfg, err := config.Load(ctx.ConfigPath)
	if err != nil {
		return err
	}

	if len(cfg.Profiles) == 0 {
		fmt.Println("No profiles configured. Use 'git-profile add <name>' to create one.")
		return nil
	}

	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PROFILE\tNAME\tEMAIL\tSIGNING KEY\tGPG SIGN\tSSH KEY")
	fmt.Fprintln(w, "-------\t----\t-----\t-----------\t--------\t-------")

	for _, name := range names {
		p := cfg.Profiles[name]
		gpg := "no"
		if p.GPGSign {
			gpg = "yes"
		}
		signingKey := valueOrDash(p.SigningKey)
		sshKey := valueOrDash(p.SSHKey)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", name, p.Name, p.Email, signingKey, gpg, sshKey)
	}

	return w.Flush()
}

func valueOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
