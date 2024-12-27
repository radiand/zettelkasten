package commands

import "errors"
import "fmt"
import "os"
import "path"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"
import "github.com/radiand/zettelkasten/internal/workspaces"

// Init creates configuration file and required directories if they do not
// exist.
type Init struct {
	ConfigPath    string
	WorkspaceName string
}

// Run performs initialization command.
func (self Init) Run() (string, error) {
	// Create config if not found.
	expandedConfigPath := common.ExpandHomeDir(self.ConfigPath)

	isConfigFile, _ := common.Exists(expandedConfigPath)
	if !isConfigFile {
		os.MkdirAll(path.Dir(expandedConfigPath), 0766)
		err := config.PutConfigToFile(expandedConfigPath, config.NewConfig())
		if err != nil {
			return "", err
		}
		out := fmt.Sprintf("Created configuration file in: %s.\nOpen it now, review and modify default values. When done, run init once again to finalize.", expandedConfigPath)
		return out, nil
	}

	configObj, _ := config.GetConfigFromFile(expandedConfigPath)
	expandedRootPath := common.ExpandHomeDir(configObj.ZettelkastenDir)
	os.MkdirAll(expandedRootPath, 0766)

	// Choose workspace name. If user did not provide any, use the default,
	// hence allowing for simple `$ zettelkasten init` to do everything that is
	// required to work.
	workspaceName := chooseFirstNonEmpty(self.WorkspaceName, configObj.DefaultWorkspace)

	_, err := workspaces.IsOkay(expandedRootPath, workspaceName)
	if errors.Is(err, workspaces.ErrOsFailure) {
		return "", err
	}
	if errors.Is(err, workspaces.ErrMalformed) {
		return "", errors.Join(
			err, fmt.Errorf(
				"Workspace %s exists, but does not conform template",
				path.Join(expandedRootPath, workspaceName),
			),
		)
	}

	if errors.Is(err, workspaces.ErrNotExists) {
		err := workspaces.CreateWorkspace(expandedRootPath, workspaceName)
		if err != nil {
			return "", err
		}
		out := fmt.Sprintf("Created workspace %s/%s.", expandedRootPath, workspaceName)
		return out, nil
	}

	return "Nothing to do.", nil
}

func chooseFirstNonEmpty(choice ...string) string {
	for _, item := range choice {
		if item != "" {
			return item
		}
	}
	panic("No valid workspace name found")
}
