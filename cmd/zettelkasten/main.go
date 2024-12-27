/*
Zettelkasten - plain text notes with toml metadata
*/
package main

import "flag"
import "fmt"
import "os"
import "time"

import "github.com/radiand/zettelkasten/internal/application"
import "github.com/radiand/zettelkasten/internal/application/commands"
import "github.com/radiand/zettelkasten/internal/application/queries"
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
	verbose    bool
	subcommand string
	subArgs    []string
}

type cmdCommitArgs struct {
	cooldown time.Duration
}

type cmdGetArgs struct {
	providePath bool
	query       []string
}

type cmdInitArgs struct {
	workspaceName string
}

type cmdNewArgs struct {
	workspaceName string
}

func parseGlobalArgs() globalArgs {
	configPath := flag.String(
		"f",
		"~/.config/zettelkasten/config.toml",
		"Path to config.toml file",
	)
	verbose := flag.Bool(
		"v",
		false,
		"Turn on verbose messages, e.g. detailed error backtrace",
	)
	usage := common.BuildUsage("zettelkasten", "Note management").WithCommands(COMMANDS)
	flag.Usage = func() { common.Flagprint(usage.Render(flag.CommandLine)) }
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify a subcommand. Available: new, link, commit.")
		os.Exit(1)
	}
	cmd, args := args[0], args[1:]

	return globalArgs{configPath: *configPath, verbose: *verbose, subcommand: cmd, subArgs: args}
}

func parseCmdNew(args []string) cmdNewArgs {
	flagset := flag.NewFlagSet("new", flag.ExitOnError)
	usage := common.BuildUsage(
		"zettelkasten new", COMMANDS["new"],
	).WithArguments(
		map[string]string{"workspace": "(optional) Workspace in which note will be created. Default from config if not specified."},
	)
	flagset.Usage = func() { common.Flagprint(usage.Render(flagset)) }
	err := flagset.Parse(args)
	try(err, "Invalid arguments")
	hint := "Provide name of the workspace to create a new note in."
	if flagset.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments. %s\n", hint)
		os.Exit(1)
	}
	workspaceName := ""
	if flagset.NArg() == 1 {
		workspaceName = flagset.Arg(0)
	}
	return cmdNewArgs{workspaceName: workspaceName}
}

func parseCmdCommit(args []string) cmdCommitArgs {
	flagset := flag.NewFlagSet("commit", flag.ExitOnError)
	cooldown := flagset.Duration(
		"c",
		time.Duration(0),
		"Setup how much time has to pass to allow commiting a file.",
	)
	usage := common.BuildUsage("zettelkasten commit", COMMANDS["commit"])
	flagset.Usage = func() { common.Flagprint(usage.Render(flagset)) }
	err := flagset.Parse(args)
	try(err, "Invalid arguments")
	return cmdCommitArgs{cooldown: *cooldown}
}

func parseCmdGet(args []string) cmdGetArgs {
	flagset := flag.NewFlagSet("get", flag.ExitOnError)
	providePath := flagset.Bool("p", false, "Print path instead of the content.")
	usage := common.BuildUsage("zettelkasten get", COMMANDS["get"])
	flagset.Usage = func() { common.Flagprint(usage.Render(flagset)) }
	err := flagset.Parse(args)
	try(err, "Invalid arguments")

	return cmdGetArgs{providePath: *providePath, query: flagset.Args()}
}

func parseCmdInit(args []string) cmdInitArgs {
	flagset := flag.NewFlagSet("init", flag.ExitOnError)
	usage := common.BuildUsage(
		"zettelkasten init", COMMANDS["init"],
	).WithArguments(
		map[string]string{"workspace": "(optional) Workspace to be created."},
	)
	flagset.Usage = func() { common.Flagprint(usage.Render(flagset)) }
	err := flagset.Parse(args)
	try(err, "Invalid arguments")
	hint := "Provide name for new workspace."
	if flagset.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments. %s\n", hint)
		os.Exit(1)
	}
	workspaceName := ""
	if flagset.NArg() == 1 {
		workspaceName = flagset.Arg(0)
	}
	return cmdInitArgs{workspaceName: workspaceName}
}

func parseCmdLink(args []string) {
	flagset := flag.NewFlagSet("link", flag.ExitOnError)
	usage := common.BuildUsage("zettelkasten link", COMMANDS["link"])
	flagset.Usage = func() { common.Flagprint(usage.Render(flagset)) }
	flagset.Parse(args)
}

func main() {
	globalArgs := parseGlobalArgs()

	if globalArgs.subcommand == "init" {
		parsedArgs := parseCmdInit(globalArgs.subArgs)
		cmdInitRunner := commands.Init{
			ConfigPath:    globalArgs.configPath,
			WorkspaceName: parsedArgs.workspaceName,
		}
		run(cmdInitRunner, globalArgs.verbose)
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
		workspaceName := config.DefaultWorkspace
		if parsedArgs.workspaceName != "" {
			workspaceName = parsedArgs.workspaceName
		}
		cmdNewRunner := commands.New{
			ZettelkastenDir: zettelkastenDir,
			WorkspaceName:   workspaceName,
			Nowtime:         common.Now,
		}
		run(cmdNewRunner, globalArgs.verbose)
	case "link":
		parseCmdLink(globalArgs.subArgs)
		cmdLinkRunner := commands.Link{
			ZettelkastenDir: zettelkastenDir,
		}
		run(cmdLinkRunner, globalArgs.verbose)
	case "commit":
		trackedDirectories := []string{zettelkastenDir}
		parsedArgs := parseCmdCommit(globalArgs.subArgs)
		cmdCommitRunner := commands.Commit{
			Dirs:       trackedDirectories,
			GitFactory: gitFactory,
			Nowtime:    common.Now,
			Modtime:    common.ModificationTime,
			Cooldown:   parsedArgs.cooldown,
		}
		run(cmdCommitRunner, globalArgs.verbose)
	case "get":
		parsedArgs := parseCmdGet(globalArgs.subArgs)
		cmdGetRunner := queries.Get{
			ConfigPath:  globalArgs.configPath,
			ProvidePath: parsedArgs.providePath,
			Query:       parsedArgs.query,
		}
		run(cmdGetRunner, globalArgs.verbose)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported command: '%s'\n", globalArgs.subcommand)
		os.Exit(1)
	}
}

func run(runnable application.Runnable, verbose bool) {
	out, err := runnable.Run()
	if err != nil {
		if verbose {
			fmt.Fprintln(os.Stderr, common.FmtErrors(err))
		} else {
			fmt.Fprintln(os.Stderr, common.LastError(err))
		}
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, out)
}

func try(err error, message string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, message)
		fmt.Fprintln(os.Stderr, common.FmtErrors(err))
		os.Exit(1)
	}
}
