package note

import "errors"
import "os"
import "path/filepath"
import "strings"

// FilesystemNoteRepository provides Notes saved on disk.
type FilesystemNoteRepository struct {
	RootDir string
}

// Get obtains Note from disk.
func (self *FilesystemNoteRepository) Get(uid string) (Note, error) {
	content, err := os.ReadFile(self.getNotePath(uid))
	if err != nil {
		return Note{}, err
	}
	return UnmarshallNote(string(content))
}

// Put saves Note to disk.
func (self *FilesystemNoteRepository) Put(note Note) (string, error) {
	marshalled, err := note.ToToml()
	if err != nil {
		return "", errors.Join(err, errors.New("Cannot marshall note"))
	}
	path := self.getNotePath(note.Header.Uid)
	err = os.WriteFile(path, []byte(marshalled), 0644)
	if err != nil {
		return "", errors.Join(err, errors.New("Cannot save note"))
	}
	return path, nil
}

// List obtains array of saved Notes' Uids.
func (self *FilesystemNoteRepository) List() ([]string, error) {
	notePaths, err := os.ReadDir(self.RootDir)
	if err != nil {
		return []string{}, errors.Join(err, errors.New("Cannot list notes"))
	}

	noteUids := []string{}
	for _, file := range notePaths {
		uidRe := GetUidRegexp()
		if matches := uidRe.MatchString(file.Name()); matches {
			noteUids = append(noteUids, strings.TrimSuffix(file.Name(), ".md"))
		}
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
