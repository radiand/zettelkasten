package main

import "errors"
import "fmt"
import "os"
import "path"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"
import "github.com/radiand/zettelkasten/internal/workspaces"

// DefaultRootPath is used when no config is available.
var DefaultRootPath = "~/vault/zettelkasten"

// DefaultWorkspaceName is used when workspace is neither defined in config nor
// set via CLI args.
var DefaultWorkspaceName = "main"

// CmdInit creates configuration file and required directories if they do not
// exist.
type CmdInit struct {
	configPath    string
	workspaceName string
}

// Run performs initialization command.
func (self *CmdInit) Run() error {
	expandedConfigPath := common.ExpandHomeDir(self.configPath)
	isConfigFile, _ := common.Exists(expandedConfigPath)

	var expandedRootPath string
	workspaceNameFromConfig := ""

	if !isConfigFile {
		expandedRootPath = common.ExpandHomeDir(DefaultRootPath)
	} else {
		configObj, _ := config.GetConfigFromFile(expandedConfigPath)
		expandedRootPath = common.ExpandHomeDir(configObj.ZettelkastenDir)
		workspaceNameFromConfig = configObj.DefaultWorkspace
	}

	if !isConfigFile {
		fmt.Printf("Creating config at %s. ", expandedConfigPath)
		accepted := common.AskBool("Proceed?")
		if accepted {
			config.PutConfigToFile(
				expandedConfigPath,
				config.Config{ZettelkastenDir: expandedRootPath},
			)
			fmt.Println("Config was created with defaults defined. You can edit them now if you wish.\nSee", expandedConfigPath)
		}
	}

	if ok, _ := common.Exists(expandedRootPath); !ok {
		os.MkdirAll(expandedRootPath, 0744)
	}

	// CLI arguments is most important, then config, then defaults.
	workspaceName := chooseFirstNonEmpty(self.workspaceName, workspaceNameFromConfig, DefaultWorkspaceName)
	_, err := workspaces.IsOkay(expandedRootPath, workspaceName)
	if errors.Is(err, workspaces.ErrOsFailure) {
		return err
	}
	if errors.Is(err, workspaces.ErrMalformed) {
		return errors.Join(
			err, fmt.Errorf(
				"Workspace %s exists, but does not conform template. Backup the content and run init once again",
				path.Join(expandedRootPath, workspaceName),
			),
		)
	}

	if errors.Is(err, workspaces.ErrNotExists) {
		fmt.Printf("Creating workspace '%s' at %s. ", workspaceName, expandedRootPath)
		accepted := common.AskBool("Proceed?")
		if accepted {
			err := workspaces.CreateWorkspace(expandedRootPath, workspaceName)
			if err != nil {
				return err
			}
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
