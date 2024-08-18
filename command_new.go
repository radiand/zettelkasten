package main

import "fmt"
import "os"

import "github.com/radiand/zettelkasten/internal/note"

// CmdNewOptions is used to carry arguments for RunCmdNew.
type CmdNewOptions struct {
	RootDir string
	Stdout  bool
}

// RunCmdNew creates new note. It can be instructed to print new note to stdout
// (default) or to file, printing only the path to created note. Printing just
// paths can be useful to integrate this application with external text
// editors.
func RunCmdNew(options CmdNewOptions) error {
	note := note.NewNote()
	marshaled, _ := note.ToToml()
	if options.Stdout {
		fmt.Print(marshaled)
	} else {
		notePath := fmt.Sprintf("%s/%s.md", options.RootDir, note.Header.Uid)
		err := os.WriteFile(notePath, []byte(marshaled), 0644)
		if err != nil {
			return fmt.Errorf("Cannot save note due to: %w", err)
		}
		fmt.Println(notePath)
	}
	return nil
}
