package application

import "fmt"
import "os"
import "path"
import "strings"
import "testing"
import "time"

import "github.com/stretchr/testify/assert"

import "github.com/radiand/zettelkasten/internal/application/commands"
import "github.com/radiand/zettelkasten/internal/notes"

func TestCreateNote(t *testing.T) {
	zkdir := t.TempDir()
	wsname := "workspace_name"
	notesDir := path.Join(zkdir, wsname, "notes")

	// Prepare directories the way init command would do it. We don't call init
	// here because it is interactive.
	os.MkdirAll(notesDir, 0777)

	cmdNew := commands.New{
		ZettelkastenDir: zkdir,
		WorkspaceName:   wsname,
		Nowtime:         time.Now,
	}

	// CmdNew.Run() should return path of newly created note to stdout.
	notePath, err := cmdNew.Run()
	assert.Nil(t, err)

	// Printed path should be absolute.
	assert.True(t, strings.HasPrefix(notePath, "/"))

	// There should be file in printed path.
	noteBytes, err := os.ReadFile(strings.TrimSpace(notePath))
	assert.Nil(t, err)
	assert.NotEmpty(t, noteBytes)

	// Extract note UID.
	dir, noteFilename := path.Split(strings.TrimSpace(notePath))
	assert.Equal(t, notesDir, strings.TrimRight(dir, "/"))
	noteUid := strings.TrimRight(noteFilename, ".md") // revive:disable-line

	// Open the note once again, but this time using proper Repository.
	noteRepo := notes.NewFilesystemNoteRepository(notesDir)
	note, err := noteRepo.Get(noteUid)
	assert.Nil(t, err)

	// Header should match filename UID.
	assert.Equal(t, note.Header.Uid, noteUid)
}

func TestLinkTwoNotes(t *testing.T) {
	zkdir := t.TempDir()
	wsname := "workspace_name"
	notesDir := path.Join(zkdir, wsname, "notes")

	// Prepare directories the way init command would do it. We don't call init
	// here because it is interactive.
	os.MkdirAll(notesDir, 0777)

	// Create two notes.
	cmdNew := commands.New{
		ZettelkastenDir: zkdir,
		WorkspaceName:   wsname,
		Nowtime:         func() time.Time { return time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC) },
	}
	_, err := cmdNew.Run()
	assert.Nil(t, err)

	cmdNew = commands.New{
		ZettelkastenDir: zkdir,
		WorkspaceName:   wsname,
		Nowtime:         func() time.Time { return time.Date(2024, 2, 2, 2, 2, 2, 2, time.UTC) },
	}
	_, err = cmdNew.Run()
	assert.Nil(t, err)

	// Obtain list of created note UIDs, there should be two.
	noteRepo := notes.NewFilesystemNoteRepository(notesDir)
	noteUids, err := noteRepo.List()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(noteUids))

	// Refer to first note from second.
	note1, _ := noteRepo.Get(noteUids[0])
	note2, _ := noteRepo.Get(noteUids[1])
	note2.Body = fmt.Sprintf("Refers to [[%s]]", note1.Header.Uid)
	_, err = noteRepo.Put(note2)
	assert.Nil(t, err)

	cmdLink := commands.Link{
		ZettelkastenDir: zkdir,
	}
	_, err = cmdLink.Run()
	assert.Nil(t, err)

	// Check if notes are really linked now and saved.
	note1, _ = noteRepo.Get(noteUids[0])
	note2, _ = noteRepo.Get(noteUids[1])

	assert.Equal(t, []string{note2.Header.Uid}, note1.Header.ReferredFrom)
	assert.Equal(t, []string{note1.Header.Uid}, note2.Header.RefersTo)
}
