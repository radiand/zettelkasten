package main

import "fmt"

import "github.com/radiand/zettelkasten/internal/note"

// CmdLinkOptions is used to carry arguments for RunCmdLink.
type CmdLinkOptions struct {
	RootDir string
}

// RunCmdLink seeks for references between notes and updates their headers if
// there are any.
func RunCmdLink(options CmdLinkOptions) error {
	repository := note.NewFilesystemNoteRepository(options.RootDir)
	err := note.LinkNotes(repository)
	if err != nil {
		return fmt.Errorf("RunCmdLink failed due to: %w", err)
	}
	return nil
}
