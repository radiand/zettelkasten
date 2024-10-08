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
	indexDir   string
}

// Run performs initialization command.
func (self *CmdInit) Run() error {
	expandedConfigPath := common.ExpandHomeDir(self.configPath)
	isConfigFile, _ := common.Exists(expandedConfigPath)

	var expandedNotesDir string
	var expandedIndexDir string

	if !isConfigFile {
		expandedNotesDir = common.ExpandHomeDir(self.notesDir)
		expandedIndexDir = common.ExpandHomeDir(self.indexDir)
	} else {
		configObj, _ := config.GetConfigFromFile(expandedConfigPath)
		expandedNotesDir = common.ExpandHomeDir(configObj.ZettelkastenDir)
		expandedIndexDir = common.ExpandHomeDir(configObj.IndexDir)
	}

	isNotesDir, _ := common.Exists(expandedNotesDir)
	isIndexDir, _ := common.Exists(expandedIndexDir)

	if !isConfigFile {
		fmt.Printf("Creating config at %s. ", expandedConfigPath)
		accepted := common.AskBool("Proceed?")
		if accepted {
			config.PutConfigToFile(
				expandedConfigPath,
				config.Config{ZettelkastenDir: expandedNotesDir, IndexDir: expandedIndexDir},
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

	if !isIndexDir {
		fmt.Printf("Creating index dir at %s. ", expandedIndexDir)
		accepted := common.AskBool("Proceed?")
		if accepted {
			os.MkdirAll(expandedIndexDir, 0744)
		}
	}

	if !isConfigFile {
		fmt.Println("Config was created with default paths defined. You can edit them now if you wish.\nSee", expandedConfigPath)
	}

	return nil
}
