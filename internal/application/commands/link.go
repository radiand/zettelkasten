package commands

import "errors"

import "github.com/radiand/zettelkasten/internal/notes"
import "github.com/radiand/zettelkasten/internal/workspaces"

// Link carries required params to run command.
type Link struct {
	ZettelkastenDir string
}

// Run seeks for references between notes and updates their headers if there
// are any.
func (self Link) Run() error {
	foundWorkspaces, err := workspaces.GetWorkspaces(self.ZettelkastenDir)
	if err != nil {
		return errors.Join(err, errors.New("Could not link because no workspaces were found"))
	}

	for _, ws := range foundWorkspaces {
		repository := notes.NewFilesystemNoteRepository(ws.GetNotesPath())
		err := notes.LinkNotes(repository)
		if err != nil {
			return errors.Join(err, errors.New("CmdLink failed"))
		}
	}

	return nil
}
