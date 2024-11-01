package main

import "errors"
import "fmt"
import "path"

import "github.com/radiand/zettelkasten/internal/note"
import "github.com/radiand/zettelkasten/internal/workspaces"

// CmdNew carries required params to run command.
type CmdNew struct {
	zettelkastenDir string
	workspaceName   string
}

// Run creates new note file and prints its path to stdout.
func (self *CmdNew) Run() error {
	newNote := note.NewNote()
	if ok, err := workspaces.IsOkay(self.zettelkastenDir, self.workspaceName); !ok {
		return errors.Join(
			err, errors.New("Cannot create note in invalid workspace. Consider initializing workspace before"),
		)
	}
	destinationDirPath := path.Join(self.zettelkastenDir, self.workspaceName, workspaces.NotesDirName)
	repo := note.NewFilesystemNoteRepository(destinationDirPath)
	notePath, err := repo.Put(newNote)
	if err != nil {
		return errors.Join(err, errors.New("Cannot save note"))
	}
	fmt.Println(notePath)
	return nil
}
