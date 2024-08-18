package note

import "fmt"
import "os"
import "path/filepath"
import "strings"

// FilesystemNoteRepository provides Notes saved on disk.
type FilesystemNoteRepository struct {
	RootDir string
}

// Get obtains Note from disk.
func (repo *FilesystemNoteRepository) Get(uid string) (Note, error) {
	content, err := os.ReadFile(repo.getNotePath(uid))
	if err != nil {
		return Note{}, err
	}
	return LoadNote(string(content))
}

// Put saves Note to disk.
func (repo *FilesystemNoteRepository) Put(_ Note) error {
	return fmt.Errorf("Not implemented")
}

// List obtains array of saved Notes' Uids.
func (repo *FilesystemNoteRepository) List() ([]string, error) {
	notePaths, err := os.ReadDir(repo.RootDir)
	if err != nil {
		return []string{}, fmt.Errorf("Cannot list notes due to: %w", err)
	}

	noteUids := []string{}
	for _, file := range notePaths {
		noteUids = append(noteUids, strings.TrimSuffix(file.Name(), ".md"))
	}
	return noteUids, nil
}

func (repo *FilesystemNoteRepository) getNotePath(uid string) string {
	return filepath.Join(repo.RootDir, uid+".md")
}
