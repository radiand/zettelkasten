package application

import "bytes"
import "os"
import "path"
import "strings"
import "testing"

import "github.com/stretchr/testify/assert"
import "github.com/radiand/zettelkasten/internal/notes"

// captureStdout calls given function and returns what it printed to stdout and
// errors that it returned.
func captureStdout(fn func() error) (string, error) {
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	capturedErr := fn()
	os.Stdout = originalStdout
	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String(), capturedErr
}

func TestCreateNote(t *testing.T) {
	zkdir := t.TempDir()
	wsname := "workspace_name"
	notesDir := path.Join(zkdir, wsname, "notes")

	// Prepare directories the way init command would do it. We don't call init
	// here because it is interactive.
	os.MkdirAll(notesDir, 0777)

	cmdNew := &CmdNew{
		ZettelkastenDir: zkdir,
		WorkspaceName:   wsname,
	}

	// CmdNew.Run() should print path of newly created note to stdout.
	notePath, err := captureStdout(cmdNew.Run)
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
	noteRepo := note.NewFilesystemNoteRepository(notesDir)
	openNote, err := noteRepo.Get(noteUid)
	assert.Nil(t, err)

	// Header should match filename UID.
	assert.Equal(t, openNote.Header.Uid, noteUid)
}
