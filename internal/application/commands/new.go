package commands

import "errors"
import "path"
import "time"

import "github.com/radiand/zettelkasten/internal/notes"
import "github.com/radiand/zettelkasten/internal/workspaces"

// New carries required params to run command.
type New struct {
	ZettelkastenDir string
	WorkspaceName   string
	Nowtime         func() time.Time
}

// Run creates new note file and prints its path to stdout.
func (self New) Run() (string, error) {
	newNote := notes.NewNote(self.Nowtime())
	if ok, err := workspaces.IsOkay(self.ZettelkastenDir, self.WorkspaceName); !ok {
		return "", errors.Join(
			err, errors.New("Cannot create note in invalid workspace. Consider initializing workspace before"),
		)
	}
	destinationDirPath := path.Join(self.ZettelkastenDir, self.WorkspaceName, workspaces.NotesDirName)
	repo := notes.NewFilesystemNoteRepository(destinationDirPath)
	notePath, err := repo.Put(newNote)
	if err != nil {
		return "", errors.Join(err, errors.New("Cannot save note"))
	}
	return notePath, nil
}
