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
func (repo *FilesystemNoteRepository) Put(note Note) error {
	marshalled, err := note.ToToml()
	if err != nil {
		return fmt.Errorf("Cannot marshall note due to: %w", err)
	}
	err = os.WriteFile(repo.getNotePath(note.Header.Uid), []byte(marshalled), 0644)
	if err != nil {
		return fmt.Errorf("Cannot save note due to: %w", err)
	}
	return nil
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

// NewFilesystemNoteRepository creates new instance of the repository.
func NewFilesystemNoteRepository(rootDir string) *FilesystemNoteRepository {
	return &FilesystemNoteRepository{RootDir: rootDir}
}

func (repo *FilesystemNoteRepository) getNotePath(uid string) string {
	return filepath.Join(repo.RootDir, uid+".md")
}
