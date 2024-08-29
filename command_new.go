package main

import "fmt"
import "os"
import "errors"

import "github.com/radiand/zettelkasten/internal/note"

// CmdNew carries required params to run command.
type CmdNew struct {
	RootDir string
	Stdout  bool
}

// Run creates new note. It can be instructed to print new note to stdout
// (default) or to file, printing only the path to created note. Printing just
// paths can be useful to integrate this application with external text
// editors.
func (cmd *CmdNew) Run() error {
	note := note.NewNote()
	marshaled, _ := note.ToToml()
	if cmd.Stdout {
		fmt.Print(marshaled)
	} else {
		notePath := fmt.Sprintf("%s/%s.md", cmd.RootDir, note.Header.Uid)
		err := os.WriteFile(notePath, []byte(marshaled), 0644)
		if err != nil {
			return errors.Join(err, errors.New("Cannot save note"))
		}
		fmt.Println(notePath)
	}
	return nil
}
