/*
Zettelkasten - plain text notes with toml metadata
*/
package main

import "flag"
import "log"
import "os"

import "github.com/radiand/zettelkasten/internal/osutils"

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
	cmdNew.BoolVar(&cmdNewFlagStdout, "s", true, "(shorthand for --stdout)")

	// Parse global flags.
	flag.Parse()

	// Parse subcommand flags.
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Please specify a subcommand.")
	}
	cmd, args := args[0], args[1:]

	config, err := GetConfigFromFile(osutils.ExpandHomeDir(*flagConfigPath))
	if err != nil {
		logger.Fatalf("Cannot load config: %s", err.Error())
	}

	switch cmd {
	case "new":
		cmdNew.Parse(args)
		options := CmdNewOptions{
			RootDir: osutils.ExpandHomeDir(config.Path),
			Stdout:  cmdNewFlagStdout,
		}
		err := RunCmdNew(options)
		if err != nil {
			logger.Fatalf("Program failed due to: %s", err.Error())
		}
	case "health":
		options := CmdHealthOptions{
			RootDir: osutils.ExpandHomeDir(config.Path),
		}
		err := RunCmdHealth(options)
		if err != nil {
			logger.Fatalf("Program failed due to: %s", err.Error())
		}
	case "link":
		options := CmdLinkOptions{
			RootDir: osutils.ExpandHomeDir(config.Path),
		}
		err := RunCmdLink(options)
		if err != nil {
			logger.Fatalf("Program failed due to: %s", err.Error())
		}
	default:
		logger.Fatalf("Unsupported command: '%s':", cmd)
	}
}
