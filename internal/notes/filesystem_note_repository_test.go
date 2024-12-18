package notes

import "os"
import "testing"

import "github.com/stretchr/testify/assert"

func TestListing(t *testing.T) {
	// GIVEN
	tmpdir := t.TempDir()
	repo := NewFilesystemNoteRepository(tmpdir)
	repo.Put(NewNote())

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

	repo.Put(NewNote())
	os.WriteFile(tmpdir + "/yolo.md", []byte("Garbage."), 0644)

	// WHEN
	uids, err := repo.List()

	// THEN
	assert.Nil(t, err)
	assert.Len(t, uids, 1)
}
