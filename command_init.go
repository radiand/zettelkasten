package main

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
	configPath    string
	rootPath      string
	workspaceName string
}

// Run performs initialization command.
func (self *CmdInit) Run() error {
	expandedConfigPath := common.ExpandHomeDir(self.configPath)
	isConfigFile, _ := common.Exists(expandedConfigPath)

	var expandedRootPath string

	if !isConfigFile {
		expandedRootPath = common.ExpandHomeDir(self.rootPath)
	} else {
		configObj, _ := config.GetConfigFromFile(expandedConfigPath)
		expandedRootPath = common.ExpandHomeDir(configObj.ZettelkastenDir)
	}

	if !isConfigFile {
		fmt.Printf("Creating config at %s. ", expandedConfigPath)
		accepted := common.AskBool("Proceed?")
		if accepted {
			config.PutConfigToFile(
				expandedConfigPath,
				config.Config{ZettelkastenDir: expandedRootPath},
			)
			fmt.Println("Config was created with default paths defined. You can edit them now if you wish.\nSee", expandedConfigPath)
		}
	}

	if ok, _ := common.Exists(expandedRootPath); !ok {
		os.MkdirAll(expandedRootPath, 0744)
	}

	_, err := workspaces.IsOkay(expandedRootPath, self.workspaceName)
	if errors.Is(err, workspaces.ErrOsFailure) {
		return err
	}
	if errors.Is(err, workspaces.ErrMalformed) {
		return errors.Join(
			err, fmt.Errorf(
				"Workspace %s exists, but does not conform template. Backup the content and run init once again",
				path.Join(expandedRootPath, self.workspaceName),
			),
		)
	}

	if errors.Is(err, workspaces.ErrNotExists) {
		fmt.Printf("Creating workspace '%s' at %s. ", self.workspaceName, expandedRootPath)
		accepted := common.AskBool("Proceed?")
		if accepted {
			err := workspaces.CreateWorkspace(expandedRootPath, self.workspaceName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
