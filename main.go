/*
Zettelkasten - plain text notes with toml metadata
*/
package main

import "flag"
import "os"
import "time"
import "fmt"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/git"

func main() {
	// Global flags.
	flagConfigPath := flag.String(
		"config",
		"~/.config/zettelkasten/config.toml",
		"Path to config.toml file",
	)

	// 'new' cmd flags.
	cmdNew := flag.NewFlagSet("new", flag.ExitOnError)
	var cmdNewFlagStdout bool
	cmdNew.BoolVar(
		&cmdNewFlagStdout,
		"stdout",
		true,
		"If true, print new note to stdout, otherwise save to file.",
	)

	cmdCommit := flag.NewFlagSet("commit", flag.ExitOnError)
	var cmdCommitFlagCooldown string
	cmdCommit.StringVar(
		&cmdCommitFlagCooldown,
		"cooldown",
		"0s",
		"Setup how much time has to pass to allow commiting a file.",
	)

	// Parse global flags.
	flag.Parse()

	// Parse subcommand flags.
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify a subcommand. Available: new, health, link, commit.")
		os.Exit(1)
	}
	cmd, args := args[0], args[1:]

	config, err := GetConfigFromFile(common.ExpandHomeDir(*flagConfigPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot get config.\n", common.FmtErrors(err))
		os.Exit(1)
	}

	zettelkastenDir := common.ExpandHomeDir(config.ZettelkastenDir)
	gitFactory := func(workdir string) git.IGit {
		return &git.ShellGit{WorktreePath: workdir}
	}

	switch cmd {
	case "new":
		cmdNew.Parse(args)
		cmdNewRunner := CmdNew{
			zettelkastenDir: zettelkastenDir,
			stdout:          cmdNewFlagStdout,
		}
		err := cmdNewRunner.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Command failed.\n", common.FmtErrors(err))
			os.Exit(1)
		}
	case "health":
		cmdHealthRunner := CmdHealth{
			zettelkastenDir: zettelkastenDir,
		}
		err := cmdHealthRunner.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Command failed.\n", common.FmtErrors(err))
			os.Exit(1)
		}
	case "link":
		cmdLinkRunner := CmdLink{
			zettelkastenDir: zettelkastenDir,
		}
		err := cmdLinkRunner.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Command failed.\n", common.FmtErrors(err))
			os.Exit(1)
		}
	case "commit":
		cmdCommit.Parse(args)
		cooldown, err := time.ParseDuration(cmdCommitFlagCooldown)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid cooldown.\n", common.FmtErrors(err))
			os.Exit(1)
		}
		cmdCommitRunner := CmdCommit{
			zettelkastenDir: zettelkastenDir,
			gitFactory:      gitFactory,
			nowtime:         common.Now,
			modtime:         common.ModificationTime,
			cooldown:        cooldown,
		}
		err = cmdCommitRunner.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Command failed.\n", common.FmtErrors(err))
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported command: '%s'\n", cmd)
		os.Exit(1)
	}
}
