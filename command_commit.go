package main

import "github.com/radiand/zettelkasten/internal/git"

// CmdCommit carries required params to run command.
type CmdCommit struct {
	RootDir string
}

// Run performs git commit with all changes that happened in RootDir directory.
func (cmd CmdCommit) Run() error {
	g := git.ShellGit{WorktreePath: cmd.RootDir}
	err := g.Add()
	if err != nil {
		return err
	}
	err = g.Commit()
	if err != nil {
		return err
	}
	return nil
}
