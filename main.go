/*
Zettelkasten - plain text notes with toml metadata
*/
package main

import "flag"
import "log"
import "os"
import "time"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/git"

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "", 0)
}

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
		log.Fatal("Please specify a subcommand.")
	}
	cmd, args := args[0], args[1:]

	config, err := GetConfigFromFile(common.ExpandHomeDir(*flagConfigPath))
	if err != nil {
		logger.Fatal("Cannot get config.\n", common.FmtErrors(err))
	}

	rootDir := common.ExpandHomeDir(config.RootDir)
	zettelkastenDir := common.ExpandHomeDir(config.ZettelkastenDir)

	switch cmd {
	case "new":
		cmdNew.Parse(args)
		cmdNewRunner := CmdNew{
			zettelkastenDir: zettelkastenDir,
			stdout:          cmdNewFlagStdout,
		}
		err := cmdNewRunner.Run()
		if err != nil {
			logger.Fatal("Command failed.\n", common.FmtErrors(err))
		}
	case "health":
		cmdHealthRunner := CmdHealth{
			zettelkastenDir: zettelkastenDir,
		}
		err := cmdHealthRunner.Run()
		if err != nil {
			logger.Fatal("Command failed.\n", common.FmtErrors(err))
		}
	case "link":
		cmdLinkRunner := CmdLink{
			zettelkastenDir: zettelkastenDir,
		}
		err := cmdLinkRunner.Run()
		if err != nil {
			logger.Fatal("Command failed.\n", common.FmtErrors(err))
		}
	case "commit":
		cmdCommit.Parse(args)
		cooldown, err := time.ParseDuration(cmdCommitFlagCooldown)
		if err != nil {
			logger.Fatal("Invalid cooldown.\n", common.FmtErrors(err))
		}
		cmdCommitRunner := CmdCommit{
			rootDir:         rootDir,
			zettelkastenDir: zettelkastenDir,
			git:             &git.ShellGit{WorktreePath: rootDir},
			nowtime:         common.Now,
			modtime:         common.ModificationTime,
			cooldown:        cooldown,
		}
		err = cmdCommitRunner.Run()
		if err != nil {
			logger.Fatal("Command failed.\n", common.FmtErrors(err))
		}
	default:
		logger.Fatalf("Unsupported command: '%s':", cmd)
	}
}
