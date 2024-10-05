package main

import "errors"
import "fmt"
import "reflect"
import "strings"

import "github.com/radiand/zettelkasten/internal/common"
import "github.com/radiand/zettelkasten/internal/config"
import "github.com/radiand/zettelkasten/internal/note"

// CmdGet allows reading and printing config.
type CmdGet struct {
	configPath string
	query      string
}

// Run executes the command.
func (self *CmdGet) Run() error {
	expandedConfigPath := common.ExpandHomeDir(self.configPath)
	configObj, err := config.GetConfigFromFile(expandedConfigPath)
	if err != nil {
		return errors.Join(err, errors.New("Could not open config"))
	}

	sanitizedQuery := strings.TrimSpace(self.query)
	if note.GetUidRegexp().MatchString(sanitizedQuery) {
		return errors.New("Getting note by UID is not implemented yet")
	}

	value := reflect.ValueOf(configObj).FieldByName(sanitizedQuery)
	if !value.IsValid() {
		return fmt.Errorf("No key with name '%s'", self.query)
	}

	fmt.Println(value)
	return nil
}
