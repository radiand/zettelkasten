package workspaces

import "errors"
import "io/fs"
import "os"
import "path"
import "slices"

// ErrOsFailure represents any read error.
var ErrOsFailure = errors.New("Cannot read path")
// ErrNotExists signals that file or dir does not exist (but filesystem is readable).
var ErrNotExists = errors.New("File or dir does not exist")
// ErrMalformed signals that file or dir exists, but does not conform current template.
var ErrMalformed = errors.New("File or dir exists, but is malformed")

// GetWorkspaces returns names of valid workspaces found in rootDir.
func GetWorkspaces(rootDir string) ([]string, error) {
	listing, err := os.ReadDir(rootDir)
	if err != nil {
		return []string{}, err
	}

	correctWorkspaces := []string{}
	for _, entry := range listing {
		if !entry.IsDir() {
			continue
		}
		if ok, _ := IsOkay(rootDir, entry.Name()); ok {
			correctWorkspaces = append(correctWorkspaces, entry.Name())
		}
	}
	return correctWorkspaces, nil
}

// CreateWorkspace creates directory tree for a workspace with given name.
func CreateWorkspace(rootDir string, workspaceName string) error {
	currentWorkspaces, err := GetWorkspaces(rootDir)
	if err != nil {
		return errors.Join(err, errors.New("Could not scan existing workspaces"))
	}
	if slices.Contains(currentWorkspaces, workspaceName) {
		return nil
	}
	os.MkdirAll(path.Join(rootDir, workspaceName, NotesDirName), 0744)
	os.MkdirAll(path.Join(rootDir, workspaceName, IndexDirName), 0744)
	return nil
}

// IsOkay checks if workspace of given name exists and contains expected
// directories and files.
func IsOkay(rootDir string, workspaceName string) (bool, error) {
	workspacePath := path.Join(rootDir, workspaceName)
	err := exists(workspacePath)
	if err != nil {
		return false, err
	}

	notesPath := path.Join(workspacePath, NotesDirName)
	err = exists(notesPath)
	if err != nil {
		return false, err
	}

	return true, nil
}

func exists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errors.Join(err, ErrNotExists)
		}
		return errors.Join(err, ErrOsFailure)
	}
	return nil
}
