package main

import "fmt"
import "os"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"

// CmdInit creates configuration file and required directories if they do not
// exist.
type CmdInit struct {
	configPath string
	notesDir   string
}

// Run performs initialization command.
func (self *CmdInit) Run() error {
	expandedConfigPath := common.ExpandHomeDir(self.configPath)
	isConfigFile, _ := common.Exists(expandedConfigPath)

	var expandedNotesDir string

	if !isConfigFile {
		expandedNotesDir = common.ExpandHomeDir(self.notesDir)
	} else {
		configObj, _ := config.GetConfigFromFile(expandedConfigPath)
		expandedNotesDir = common.ExpandHomeDir(configObj.ZettelkastenDir)
	}

	isNotesDir, _ := common.Exists(expandedNotesDir)

	if !isConfigFile {
		fmt.Printf("Creating config at %s. ", expandedConfigPath)
		accepted := common.AskBool("Proceed?")
		if accepted {
			config.PutConfigToFile(
				expandedConfigPath,
				config.Config{ZettelkastenDir: expandedNotesDir},
			)
		}
	}

	if !isNotesDir {
		fmt.Printf("Creating notes dir at %s. ", expandedNotesDir)
		accepted := common.AskBool("Proceed?")
		if accepted {
			os.MkdirAll(expandedNotesDir, 0744)
		}
	}

	if !isConfigFile {
		fmt.Println("Config was created with default paths defined. You can edit them now if you wish.\nSee", expandedConfigPath)
	}

	return nil
}
