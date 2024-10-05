/*
Zettelkasten - plain text notes with toml metadata
*/
package main

import "flag"
import "os"
import "time"
import "fmt"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"
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
		false,
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

	cmdGet := flag.NewFlagSet("get", flag.ExitOnError)

	// Parse global flags.
	flag.Parse()

	// Parse subcommand flags.
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify a subcommand. Available: new, health, link, commit.")
		os.Exit(1)
	}
	cmd, args := args[0], args[1:]

	if cmd == "init" {
		cmdInitRunner := CmdInit{
			configPath: *flagConfigPath,
			notesDir:   "~/vault/zettelkasten/notes",
		}
		try(cmdInitRunner.Run(), "Command failed.")
		os.Exit(0)
	}

	config, err := config.GetConfigFromFile(common.ExpandHomeDir(*flagConfigPath))
	try(err, "Cannot get config.")

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
		try(cmdNewRunner.Run(), "Command failed.")
	case "health":
		cmdHealthRunner := CmdHealth{
			zettelkastenDir: zettelkastenDir,
		}
		try(cmdHealthRunner.Run(), "Command failed.")
	case "link":
		cmdLinkRunner := CmdLink{
			zettelkastenDir: zettelkastenDir,
		}
		try(cmdLinkRunner.Run(), "Command failed.")
	case "commit":
		cmdCommit.Parse(args)
		cooldown, err := time.ParseDuration(cmdCommitFlagCooldown)
		try(err, "Invalid cooldown.")
		cmdCommitRunner := CmdCommit{
			dirs:       []string{zettelkastenDir},
			gitFactory: gitFactory,
			nowtime:    common.Now,
			modtime:    common.ModificationTime,
			cooldown:   cooldown,
		}
		try(cmdCommitRunner.Run(), "Command failed.")
	case "get":
		cmdGet.Parse(args)

		hint := "Provide key from config or note UID."
		if cmdGet.NArg() < 1 {
			fmt.Fprintf(os.Stderr, "Argument required. %s\n", hint)
			os.Exit(1)
		}
		if cmdGet.NArg() > 1 {
			fmt.Fprintf(os.Stderr, "Too many arguments. %s\n", hint)
			os.Exit(1)
		}

		cmdGetRunner := CmdGet{
			configPath: *flagConfigPath,
			query:      cmdGet.Arg(0),
		}
		try(cmdGetRunner.Run(), "Command failed.")
	default:
		fmt.Fprintf(os.Stderr, "Unsupported command: '%s'\n", cmd)
		os.Exit(1)
	}
}

func try(err error, message string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, message+"\n", common.FmtErrors(err))
		os.Exit(1)
	}
}
