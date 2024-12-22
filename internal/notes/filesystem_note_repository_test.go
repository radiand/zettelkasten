package notes

import "os"
import "testing"
import "time"

import "github.com/stretchr/testify/assert"

func TestListing(t *testing.T) {
	// GIVEN
	tmpdir := t.TempDir()
	repo := NewFilesystemNoteRepository(tmpdir)
	repo.Put(NewNote(time.Now()))

	// WHEN
	uids, err := repo.List()

	// THEN
	assert.Nil(t, err)
	assert.Len(t, uids, 1)
}

func TestListingIgnoresInvalidFilenames(t *testing.T) {
	// GIVEN
	tmpdir := t.TempDir()
	repo := NewFilesystemNoteRepository(tmpdir)

	repo.Put(NewNote(time.Now()))
	os.WriteFile(tmpdir+"/yolo.md", []byte("Garbage."), 0644)

	// WHEN
	uids, err := repo.List()

	// THEN
	assert.Nil(t, err)
	assert.Len(t, uids, 1)
}
