package main

import "errors"
import "fmt"
import "os"
import "reflect"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"
import "github.com/radiand/zettelkasten/internal/note"
import "github.com/radiand/zettelkasten/internal/workspaces"

// CmdGet allows reading and printing config.
type CmdGet struct {
	configPath  string
	providePath bool
	query       []string
}

// Run executes the command.
func (self *CmdGet) Run() error {
	expandedConfigPath := common.ExpandHomeDir(self.configPath)
	configObj, err := config.GetConfigFromFile(expandedConfigPath)
	if err != nil {
		return errors.Join(err, errors.New("Could not open config"))
	}

	if len(self.query) == 0 {
		fmt.Fprintf(os.Stderr, "Query must contain resource and key.\n")
		return errors.New("Invalid query")
	}

	switch self.query[0] {
	case "config":
		return handleConfigQuery(configObj, self.query)
	case "note":
		return handleNoteQuery(configObj, self.query, self.providePath)
	case "workspace":
		return handleWorkspaceQuery(configObj, self.query, self.providePath)
	}

	return fmt.Errorf("Resource '%s' is not supported", self.query[0])
}

func handleConfigQuery(cfg config.Config, query []string) error {
	if len(query) < 2 {
		fmt.Fprintf(os.Stderr, "To seek configuration, desired key is required.\n")
		return errors.New("Invalid query")
	}
	value := reflect.ValueOf(cfg).FieldByName(query[1])
	if !value.IsValid() {
		fmt.Fprintf(os.Stderr, "No key with name '%s'\n", query)
		return errors.New("Invalid query")
	}
	fmt.Println(value)
	return nil
}

func handleNoteQuery(cfg config.Config, query []string, providePath bool) error {
	if len(query) < 2 {
		fmt.Fprint(os.Stderr, "Missing note UID.\n")
		return errors.New("Invalid query")
	}

	uid := query[1]

	if !note.GetUidRegexp().MatchString(uid) {
		fmt.Fprintf(os.Stderr, "%s is not a valid note UID.\n", uid)
		return errors.New("Invalid query")
	}

	expandedRootPath := common.ExpandHomeDir(cfg.ZettelkastenDir)
	foundWorkspaces, _ := workspaces.GetWorkspaces(expandedRootPath)
	for _, ws := range foundWorkspaces {
		noteRepo := note.NewFilesystemNoteRepository(ws.GetNotesPath())
		noteObj, err := noteRepo.Get(uid)
		if err != nil {
			continue
		}
		if providePath {
			fmt.Println(noteRepo.GetNotePath(noteObj.Header.Uid))
			return nil
		}
		marshalled, err := noteObj.ToToml()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s could not be marshalled.\n", noteObj.Header.Uid)
			return err
		}
		fmt.Println(marshalled)
		return nil
	}
	fmt.Fprintf(os.Stderr, "Could not find note with UID %s\n", uid)
	return errors.New("Invalid query")
}

func handleWorkspaceQuery(cfg config.Config, query []string, providePath bool) error {
	if len(query) > 1 {
		fmt.Fprint(os.Stderr, "Querying workspaces does not accept additional arguments.\n")
	}

	expandedRootPath := common.ExpandHomeDir(cfg.ZettelkastenDir)
	foundWorkspaces, err := workspaces.GetWorkspaces(expandedRootPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find any workspaces in %s", expandedRootPath)
		return err
	}

	for _, ws := range foundWorkspaces {
		if providePath {
			fmt.Println(ws.GetWorkspacePath())
		} else {
			fmt.Println(ws.GetName())
		}
	}
	return nil
}
