package application

import "errors"
import "fmt"
import "os"
import "path"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"
import "github.com/radiand/zettelkasten/internal/workspaces"

// CmdInit creates configuration file and required directories if they do not
// exist.
type CmdInit struct {
	ConfigPath    string
	WorkspaceName string
}

// Run performs initialization command.
func (self *CmdInit) Run() error {
	// Create config if not found.
	expandedConfigPath := common.ExpandHomeDir(self.ConfigPath)

	isConfigFile, _ := common.Exists(expandedConfigPath)
	if !isConfigFile {
		fmt.Printf("Creating configuration file in %s.\n", expandedConfigPath)
		os.MkdirAll(path.Dir(expandedConfigPath), 0766)
		err := config.PutConfigToFile(expandedConfigPath, config.NewConfig())
		if err != nil {
			return err
		}
		fmt.Println("Please open configuration file, review default values and modify them as you wish. When done, run init command once again to finalize.")
		return nil
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
		return err
	}
	if errors.Is(err, workspaces.ErrMalformed) {
		return errors.Join(
			err, fmt.Errorf(
				"Workspace %s exists, but does not conform template",
				path.Join(expandedRootPath, workspaceName),
			),
		)
	}

	if errors.Is(err, workspaces.ErrNotExists) {
		fmt.Printf("Creating workspace %s/%s.\n", expandedRootPath, workspaceName)
		err := workspaces.CreateWorkspace(expandedRootPath, workspaceName)
		if err != nil {
			return err
		}
	}

	return nil
}

func chooseFirstNonEmpty(choice ...string) string {
	for _, item := range choice {
		if item != "" {
			return item
		}
	}
	panic("No valid workspace name found")
}
