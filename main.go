/*
Zettelkasten - plain text notes with toml metadata
*/
package main

import "flag"
import "fmt"
import "os"
import "strings"
import "time"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"
import "github.com/radiand/zettelkasten/internal/git"

// COMMANDS stores help string for all subcommands.
var COMMANDS = map[string]string{
	"init":   "Create config and required directories.",
	"new":    "Create new note.",
	"link":   "Find link between notes and update headers.",
	"get":    "Get key from config or note by UID.",
	"commit": "Generate commit message and execute git commit.",
}

type globalArgs struct {
	configPath string
	subcommand string
	subArgs    []string
}

type cmdCommitArgs struct {
	cooldown time.Duration
}

type cmdGetArgs struct {
	query string
}

type cmdNewArgs struct {
	stdout bool
}

func flagprint(lines []string) {
	fmt.Fprintln(
		flag.CommandLine.Output(),
		strings.Join(
			lines,
			"\n",
		),
	)
}

func subcommandUsage(name string, help string) {
	flagprint([]string{
		help,
		"",
		"Usage:",
		"  zettelkasten [global-options] " + name,
	})
}

func subcommandUsageWithOptions(flagset *flag.FlagSet, name string, help string) {
	flagprint([]string{
		help,
		"",
		"Usage:",
		"  zettelkasten [global-options] " + name + " [options]",
		"",
		"Options:",
	})
	flagset.PrintDefaults()
}

func parseGlobalArgs() globalArgs {
	configPath := flag.String(
		"config",
		"~/.config/zettelkasten/config.toml",
		"Path to config.toml file",
	)
	flag.Usage = func() {
		flagprint([]string{
			"Usage:",
			"  zettelkasten [options] <command>",
			"",
			"Options:",
		})
		flag.PrintDefaults()

		commandsLines := []string{
			"",
			"Available commands:",
		}
		for cmdName, cmdHelp := range COMMANDS {
			commandsLines = append(commandsLines, fmt.Sprintf("  %-10s %s", cmdName, cmdHelp))
		}
		flagprint(commandsLines)
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify a subcommand. Available: new, link, commit.")
		os.Exit(1)
	}
	cmd, args := args[0], args[1:]

	return globalArgs{configPath: *configPath, subcommand: cmd, subArgs: args}
}

func parseCmdNew(args []string) cmdNewArgs {
	flagset := flag.NewFlagSet("new", flag.ExitOnError)
	stdout := flagset.Bool(
		"stdout",
		false,
		"If true, print new note to stdout, otherwise save to file.",
	)
	flagset.Usage = func() { subcommandUsageWithOptions(flagset, "new", COMMANDS["new"]) }
	err := flagset.Parse(args)
	try(err, "Invalid arguments")
	return cmdNewArgs{stdout: *stdout}
}

func parseCmdCommit(args []string) cmdCommitArgs {
	flagset := flag.NewFlagSet("commit", flag.ExitOnError)
	cooldown := flagset.Duration(
		"cooldown",
		time.Duration(0),
		"Setup how much time has to pass to allow commiting a file.",
	)
	flagset.Usage = func() { subcommandUsageWithOptions(flagset, "commit", COMMANDS["commit"]) }
	err := flagset.Parse(args)
	try(err, "Invalid arguments")
	return cmdCommitArgs{cooldown: *cooldown}
}

func parseCmdGet(args []string) cmdGetArgs {
	flagset := flag.NewFlagSet("get", flag.ExitOnError)
	flagset.Usage = func() { subcommandUsage("get", COMMANDS["get"]) }
	err := flagset.Parse(args)
	try(err, "Invalid arguments")

	hint := "Provide key from config or note UID."
	if flagset.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Argument required. %s\n", hint)
		os.Exit(1)
	}
	if flagset.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments. %s\n", hint)
		os.Exit(1)
	}

	return cmdGetArgs{query: flagset.Arg(0)}
}

func parseCmdInit(args []string) {
	flagset := flag.NewFlagSet("init", flag.ExitOnError)
	flagset.Usage = func() { subcommandUsage("init", COMMANDS["init"]) }
	flagset.Parse(args)
}

func parseCmdLink(args []string) {
	flagset := flag.NewFlagSet("link", flag.ExitOnError)
	flagset.Usage = func() { subcommandUsage("link", COMMANDS["link"]) }
	flagset.Parse(args)
}

func main() {
	globalArgs := parseGlobalArgs()

	if globalArgs.subcommand == "init" {
		parseCmdInit(globalArgs.subArgs)
		cmdInitRunner := CmdInit{
			configPath: globalArgs.configPath,
			notesDir:   "~/vault/zettelkasten/notes",
			indexDir:   "~/vault/zettelkasten/index",
		}
		try(cmdInitRunner.Run(), "Command failed.")
		os.Exit(0)
	}

	config, err := config.GetConfigFromFile(common.ExpandHomeDir(globalArgs.configPath))
	try(err, "Cannot get config.")

	zettelkastenDir := common.ExpandHomeDir(config.ZettelkastenDir)
	gitFactory := func(workdir string) git.IGit {
		return &git.ShellGit{WorktreePath: workdir}
	}

	switch globalArgs.subcommand {
	case "new":
		parsedArgs := parseCmdNew(globalArgs.subArgs)
		cmdNewRunner := CmdNew{
			zettelkastenDir: zettelkastenDir,
			stdout:          parsedArgs.stdout,
		}
		try(cmdNewRunner.Run(), "Command failed.")
	case "link":
		parseCmdLink(globalArgs.subArgs)
		cmdLinkRunner := CmdLink{
			zettelkastenDir: zettelkastenDir,
		}
		try(cmdLinkRunner.Run(), "Command failed.")
	case "commit":
		trackedDirectories := []string{zettelkastenDir}
		isIndexDirSet := len(config.IndexDir) != 0
		if isIndexDirSet {
			trackedDirectories = append(
				trackedDirectories, common.ExpandHomeDir(config.IndexDir),
			)
		}
		parsedArgs := parseCmdCommit(globalArgs.subArgs)
		cmdCommitRunner := CmdCommit{
			dirs:       trackedDirectories,
			gitFactory: gitFactory,
			nowtime:    common.Now,
			modtime:    common.ModificationTime,
			cooldown:   parsedArgs.cooldown,
		}
		try(cmdCommitRunner.Run(), "Command failed.")
	case "get":
		parsedArgs := parseCmdGet(globalArgs.subArgs)
		cmdGetRunner := CmdGet{
			configPath: globalArgs.configPath,
			query:      parsedArgs.query,
		}
		try(cmdGetRunner.Run(), "Command failed.")
	default:
		fmt.Fprintf(os.Stderr, "Unsupported command: '%s'\n", globalArgs.subcommand)
		os.Exit(1)
	}
}

func try(err error, message string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, message+"\n", common.FmtErrors(err))
		os.Exit(1)
	}
}
