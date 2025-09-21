package queries

import "errors"
import "fmt"
import "reflect"
import "strings"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"
import "github.com/radiand/zettelkasten/internal/notes"
import "github.com/radiand/zettelkasten/internal/workspaces"

// Get allows reading and printing config.
type Get struct {
	ConfigPath  string
	ProvidePath bool
	Query       []string
}

// Run executes the command.
func (self Get) Run() (string, error) {
	expandedConfigPath := common.ExpandHomeDir(self.ConfigPath)
	configObj, err := config.GetConfigFromFile(expandedConfigPath)
	if err != nil {
		return "", errors.Join(err, errors.New("Could not open config"))
	}

	if len(self.Query) == 0 {
		return "", errors.New("Query must specify resource (available: config, note, workspace)")
	}

	switch self.Query[0] {
	case "config":
		return handleConfigQuery(configObj, self.Query)
	case "note":
		return handleNoteQuery(configObj, self.Query, self.ProvidePath)
	case "notes":
		return handleNotesQuery(configObj, self.Query, self.ProvidePath)
	case "workspace", "workspaces":
		return handleWorkspaceQuery(configObj, self.Query, self.ProvidePath)
	}

	return "", fmt.Errorf("Resource '%s' is not supported", self.Query[0])
}

func handleConfigQuery(cfg config.Config, query []string) (string, error) {
	if len(query) < 2 {
		return "", errors.New("Query must specify config field to be read")
	}
	value := reflect.ValueOf(cfg).FieldByName(query[1])
	if !value.IsValid() {
		return "", fmt.Errorf("No key with name '%s", query)
	}
	return fmt.Sprintf("%v", value), nil
}

func handleNoteQuery(cfg config.Config, query []string, providePath bool) (string, error) {
	if len(query) < 2 {
		return "", errors.New("Query must contain note UID")
	}

	uid := query[1]

	if !notes.GetUidRegexp().MatchString(uid) {
		return "", fmt.Errorf("%s is not valid note UID", uid)
	}

	expandedRootPath := common.ExpandHomeDir(cfg.ZettelkastenDir)
	foundWorkspaces, _ := workspaces.GetWorkspaces(expandedRootPath)
	for _, ws := range foundWorkspaces {
		noteRepo := notes.NewFilesystemNoteRepository(ws.GetNotesPath())
		noteObj, err := noteRepo.Get(uid)
		if err != nil {
			continue
		}
		if providePath {
			return noteRepo.GetNotePath(noteObj.Header.Uid), nil
		}
		marshalled, err := noteObj.ToToml()
		if err != nil {
			return "", fmt.Errorf("%s could not be marshalled", noteObj.Header.Uid)
		}
		return marshalled, nil
	}
	return "", fmt.Errorf("Could not find note with UID %s", uid)
}

func handleNotesQuery(cfg config.Config, query []string, providePath bool) (string, error) {
	var selectedWorkspace *string
	if len(query) > 1 {
		selectedWorkspace = &query[1]
	}

	expandedRootPath := common.ExpandHomeDir(cfg.ZettelkastenDir)
	foundWorkspaces, err := workspaces.GetWorkspaces(expandedRootPath)

	if err != nil {
		return "", fmt.Errorf("Could not find any workspaces in %s", expandedRootPath)
	}

	lines := []string{}
	for _, ws := range foundWorkspaces {
		if selectedWorkspace != nil && ws.GetName() != *selectedWorkspace {
			continue
		}
		noteRepo := notes.NewFilesystemNoteRepository(ws.GetNotesPath())
		uids, err := noteRepo.List()
		if err != nil {
			return "", fmt.Errorf("Cannot list notes in workspace %s", ws.GetName())
		}
		for _, uid := range uids {
			if providePath {
				lines = append(lines, noteRepo.GetNotePath(uid))
			} else {
				lines = append(lines, uid)
			}
		}
	}
	return strings.Join(lines, "\n"), nil
}

func handleWorkspaceQuery(cfg config.Config, query []string, providePath bool) (string, error) {
	if len(query) > 1 {
		return "", errors.New("Querying workspaces does not accept additional arguments")
	}

	expandedRootPath := common.ExpandHomeDir(cfg.ZettelkastenDir)
	foundWorkspaces, err := workspaces.GetWorkspaces(expandedRootPath)

	if err != nil {
		return "", fmt.Errorf("Could not find any workspaces in %s", expandedRootPath)
	}

	lines := []string{}
	for _, ws := range foundWorkspaces {
		if providePath {
			lines = append(lines, ws.GetWorkspacePath())
		}
		lines = append(lines, ws.GetName())
	}
	return strings.Join(lines, "\n"), nil
}
