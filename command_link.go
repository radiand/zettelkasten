package main

import "errors"

import "github.com/radiand/zettelkasten/internal/note"

// CmdLink carries required params to run command.
type CmdLink struct {
	RootDir string
}

// Run seeks for references between notes and updates their headers if there
// are any.
func (cmd *CmdLink) Run() error {
	repository := note.NewFilesystemNoteRepository(cmd.RootDir)
	err := note.LinkNotes(repository)
	if err != nil {
		return errors.Join(err, errors.New("CmdLink failed"))
	}
	return nil
}
