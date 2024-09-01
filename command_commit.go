package main

import "github.com/radiand/zettelkasten/internal/git"

// CmdCommit carries required params to run command.
type CmdCommit struct {
	rootDir string
	git   git.IGit
}

// Run performs git commit with all changes that happened in RootDir directory.
func (cmd CmdCommit) Run() error {
	err := cmd.git.Add(cmd.rootDir)
	if err != nil {
		return err
	}
	err = cmd.git.Commit()
	if err != nil {
		return err
	}
	return nil
}
