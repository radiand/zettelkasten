package main

import "errors"
import "fmt"

import "github.com/radiand/zettelkasten/internal/note"

// CmdHealthOptions is used to carry arguments for CmdNewOptions.
type CmdHealthOptions struct {
	RootDir string
}

// RunCmdHealth tries to read all notes and checks if they can be correctly
// parsed.
func RunCmdHealth(options CmdHealthOptions) error {
	repository := note.NewFilesystemNoteRepository(options.RootDir)
	uids, err := repository.List()
	if err != nil {
		return err
	}

	var health = make(map[string]error)
	for _, uid := range uids {
		_, err := repository.Get(uid)
		if err != nil {
			health[uid] = err
		}
	}
	if len(health) == 0 {
		return nil
	}
	for uid, err := range health {
		fmt.Printf("%s: %s\n", uid, err.Error())
	}
	return errors.New("Found unhealthy notes")
}
