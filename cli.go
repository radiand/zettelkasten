package main

import "flag"
import "fmt"
import "strings"

func flagprint(lines []string) {
	fmt.Fprintln(
		flag.CommandLine.Output(),
		strings.Join(
			lines,
			"\n",
		),
	)
}

func renderOptions(flagset *flag.FlagSet) []string {
	lines := []string{}
	flagset.VisitAll(
		func(f *flag.Flag) {
			lines = append(lines, fmt.Sprintf("  -%-9s %s (default: %s)", f.Name, f.Usage, f.DefValue))
		},
	)
	if len(lines) > 0 {
		lines = append([]string{"Options:"}, lines...)
		lines = append(lines, "")
	}
	return lines
}

func renderPositionals(header string, positionals map[string]string) []string {
	if len(positionals) == 0 {
		return []string{}
	}

	lines := []string{header}
	for argName, argHelp := range positionals {
		lines = append(lines, fmt.Sprintf("  %-10s %s", argName, argHelp))
	}
	return append(lines, "")
}

func renderCommands(commands map[string]string) []string {
	return renderPositionals("Commands:", commands)
}

func renderArguments(arguments map[string]string) []string {
	return renderPositionals("Positional arguments:", arguments)
}

// Usage can be rendered in CLI as command (or subcommand) help.
type Usage struct {
	name      string
	help      string
	commands  map[string]string
	arguments map[string]string
}

// BuildUsage creates new instance of Usage.
func BuildUsage(name string, help string) Usage {
	return Usage{name: name, help: help, commands: map[string]string{}, arguments: map[string]string{}}
}

// WithCommands attaches given command to the final Usage. Subcommands should
// not be mixed witn Arguments.
func (self Usage) WithCommands(commands map[string]string) Usage {
	self.commands = commands
	return self
}

// WithArguments attaches given positional arguments to the final Usage. Arguments should
// not be mixed witn Commands.
func (self Usage) WithArguments(arguments map[string]string) Usage {
	self.arguments = arguments
	return self
}

// Render creates slice of strings that can be joined to create final usage
// string ready to be printed.
func (self Usage) Render(flagset *flag.FlagSet) []string {
	optionsRendered := renderOptions(flagset)
	commandsRendered := renderCommands(self.commands)
	argumentsRendered := renderArguments(self.arguments)

	invocation := "  " + self.name
	if len(optionsRendered) > 0 {
		invocation += " [<options>]"
	}
	if len(commandsRendered) > 0 {
		invocation += " <command>"
	}
	for argName := range self.arguments {
		invocation += " <" + argName + ">"
	}
	lines := []string{
		fmt.Sprintf("%s: %s", self.name, self.help),
		"",
		"Usage:",
		invocation,
		"",
	}
	lines = append(lines, optionsRendered...)
	lines = append(lines, commandsRendered...)
	lines = append(lines, argumentsRendered...)
	return lines
}
