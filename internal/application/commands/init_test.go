package commands

import "testing"
import "path"
import "os"

import "github.com/stretchr/testify/assert"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"

func TestInit(t *testing.T) {
	tempDir := t.TempDir()
	configPath := path.Join(tempDir, "config", "zettelkasten.toml")

	// First run, it should create config file.
	cmdInit := Init{
		ConfigPath:    configPath,
		WorkspaceName: "",
	}
	_, err := cmdInit.Run()
	assert.Nil(t, err)

	// Check if config was created.
	configObj, err := config.GetConfigFromFile(configPath)
	assert.Nil(t, err)

	// Modify default values.
	os.Remove(configPath)
	zkdir := path.Join(tempDir, "zkdir")
	configObj.ZettelkastenDir = zkdir
	config.PutConfigToFile(configPath, configObj)

	// Second run, it should create workspace.
	cmdInit = Init{
		ConfigPath:    configPath,
		WorkspaceName: "",
	}
	_, err = cmdInit.Run()
	assert.Nil(t, err)

	// Check if workspace was created.
	isWorkspace, _ := common.Exists(path.Join(tempDir, "zkdir", "main", "notes"))
	assert.Equal(t, true, isWorkspace)

	// Third run, create new workspace.
	cmdInit = Init{
		ConfigPath:    configPath,
		WorkspaceName: "fantastic_ws",
	}
	_, err = cmdInit.Run()
	assert.Nil(t, err)

	// Check if workspace was created.
	isWorkspace, _ = common.Exists(path.Join(tempDir, "zkdir", "fantastic_ws", "notes"))
	assert.Equal(t, true, isWorkspace)
}
