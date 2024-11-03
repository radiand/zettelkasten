package main

import "errors"
import "fmt"
import "reflect"
import "strings"
import "path"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"
import "github.com/radiand/zettelkasten/internal/note"
import "github.com/radiand/zettelkasten/internal/workspaces"

// CmdGet allows reading and printing config.
type CmdGet struct {
	configPath  string
	providePath bool
	query       string
}

// Run executes the command.
func (self *CmdGet) Run() error {
	expandedConfigPath := common.ExpandHomeDir(self.configPath)
	configObj, err := config.GetConfigFromFile(expandedConfigPath)
	if err != nil {
		return errors.Join(err, errors.New("Could not open config"))
	}
	expandedRootPath := common.ExpandHomeDir(configObj.ZettelkastenDir)

	sanitizedQuery := strings.TrimSpace(self.query)
	if note.GetUidRegexp().MatchString(sanitizedQuery) {
		foundWorkspaces, _ := workspaces.GetWorkspaces(expandedRootPath)
		for _, ws := range foundWorkspaces {
			noteRepo := note.NewFilesystemNoteRepository(path.Join(expandedRootPath, ws, workspaces.NotesDirName))
			noteObj, err := noteRepo.Get(sanitizedQuery)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if self.providePath {
				fmt.Println(noteRepo.GetNotePath(noteObj.Header.Uid))
				return nil
			}
			marshalled, err := noteObj.ToToml()
			if err != nil {
				return err
			}
			fmt.Println(marshalled)
			return nil
		}
		return fmt.Errorf("Could not get note with UID: %s", sanitizedQuery)
	}

	value := reflect.ValueOf(configObj).FieldByName(sanitizedQuery)
	if !value.IsValid() {
		return fmt.Errorf("No key with name '%s'", self.query)
	}

	fmt.Println(value)
	return nil
}
