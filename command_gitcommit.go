package main

import "github.com/radiand/zettelkasten/internal/git"

// CmdGitCommitOptions is used to carry arguments for RunCmdGitCommit.
type CmdGitCommitOptions struct {
	RootDir string
}

// RunCmdGitCommit performs git commit with all changes that happened in
// RootDir directory.
func RunCmdGitCommit(options CmdGitCommitOptions) error {
	g := git.ShellGit{WorktreePath: options.RootDir}
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
