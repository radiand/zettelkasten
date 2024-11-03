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

// Workspace represents single workspace directory.
type Workspace struct {
	rootPath      string
	workspaceName string
}

// GetNotesPath constructs absolute path to notes directory. Can be passed to
// INoteRepository then.
func (self Workspace) GetNotesPath() string {
	return path.Join(self.rootPath, self.workspaceName, NotesDirName)
}

// GetWorkspacePath constructs absolute path to workspace.
func (self Workspace) GetWorkspacePath() string {
	return path.Join(self.rootPath, self.workspaceName)
}

// GetName provides name of the Workspace.
func (self Workspace) GetName() string {
	return self.workspaceName
}

// GetWorkspaces returns correct workspaces.
func GetWorkspaces(rootPath string) ([]Workspace, error) {
	listing, err := os.ReadDir(rootPath)
	if err != nil {
		return []Workspace{}, err
	}

	correctWorkspaces := []Workspace{}
	for _, entry := range listing {
		if !entry.IsDir() {
			continue
		}
		if ok, _ := IsOkay(rootPath, entry.Name()); ok {
			correctWorkspaces = append(
				correctWorkspaces,
				Workspace{rootPath: rootPath, workspaceName: entry.Name()},
			)
		}
	}
	return correctWorkspaces, nil
}

// GetWorkspaceNames returns names of valid workspaces found in rootDir.
func GetWorkspaceNames(rootDir string) ([]string, error) {
	foundWorkspaces, err := GetWorkspaces(rootDir)
	if err != nil {
		return []string{}, err
	}
	foundWorkspaceNames := []string{}
	for _, ws := range foundWorkspaces {
		foundWorkspaceNames = append(foundWorkspaceNames, ws.GetName())
	}
	return foundWorkspaceNames, nil
}

// CreateWorkspace creates directory tree for a workspace with given name.
func CreateWorkspace(rootDir string, workspaceName string) error {
	currentWorkspaces, err := GetWorkspaceNames(rootDir)
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
