package application

import "errors"
import "fmt"
import "path"
import "time"

import "github.com/radiand/zettelkasten/internal/notes"
import "github.com/radiand/zettelkasten/internal/workspaces"

// CmdNew carries required params to run command.
type CmdNew struct {
	ZettelkastenDir string
	WorkspaceName   string
	Nowtime         func() time.Time
}

// Run creates new note file and prints its path to stdout.
func (self *CmdNew) Run() error {
	newNote := notes.NewNote(self.Nowtime())
	if ok, err := workspaces.IsOkay(self.ZettelkastenDir, self.WorkspaceName); !ok {
		return errors.Join(
			err, errors.New("Cannot create note in invalid workspace. Consider initializing workspace before"),
		)
	}
	destinationDirPath := path.Join(self.ZettelkastenDir, self.WorkspaceName, workspaces.NotesDirName)
	repo := notes.NewFilesystemNoteRepository(destinationDirPath)
	notePath, err := repo.Put(newNote)
	if err != nil {
		return errors.Join(err, errors.New("Cannot save note"))
	}
	fmt.Println(notePath)
	return nil
}
