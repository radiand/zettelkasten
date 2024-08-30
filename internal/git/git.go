/*
Package git provides just enough integration with git.
*/
package git

import "bytes"
import "errors"
import "fmt"
import "os/exec"
import "strings"

// FileChanges represents number of files added, deleted, etc. in current
// staging.
type FileChanges struct {
	Added    int
	Copied   int
	Deleted  int
	Modified int
	Renamed  int
}

// Any returns true if at least one file has been changed.
func (fc *FileChanges) Any() bool {
	a := fc.Added == 0
	c := fc.Copied == 0
	d := fc.Deleted == 0
	m := fc.Modified == 0
	r := fc.Renamed == 0
	if a || c || d || m || r {
		return true
	}
	return false
}

// IGit interface provides version control functionalities with git.
type IGit interface {
	Add() error
	Commit() error
	Stat(staged bool) (FileChanges, error)
}

// ShellGit is a Git interface implementation based on spawning shell process.
type ShellGit struct {
	WorktreePath string
}

// Add performs file staging.
func (instance *ShellGit) Add() error {
	cmd := exec.Command("git", "add", instance.WorktreePath)
	cmd.Dir = instance.WorktreePath
	_, err := cmd.Output()
	if err != nil {
		return errors.Join(err, errors.New("Cannot perform git add"))
	}
	return nil
}

// Commit performs git commit with custom message.
func (instance *ShellGit) Commit() error {
	changes, err := instance.Stat(true)
	if err != nil {
		return errors.Join(
			err,
			errors.New("Cannot commit changes because could not obtain change list"),
		)
	}
	if !changes.Any() {
		return nil
	}

	changesStringified := []string{}
	populateWith := func(value int, adjective string) {
		if value == 0 {
			return
		}
		changesStringified = append(
			changesStringified,
			fmt.Sprintf("%d %s", value, adjective),
		)
	}
	populateWith(changes.Added, "added")
	populateWith(changes.Copied, "copied")
	populateWith(changes.Deleted, "deleted")
	populateWith(changes.Modified, "modified")
	populateWith(changes.Renamed, "renamed")
	commitmsg := "auto: " + strings.Join(changesStringified, ", ")

	cmd := exec.Command("git", "commit", "-m", commitmsg)
	cmd.Dir = instance.WorktreePath
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return errors.Join(err, errors.New("git commit failed"))
	}
	if len(stderr.Bytes()) > 0 {
		return errors.Join(
			err,
			fmt.Errorf("git commit returned: %s", string(stderr.Bytes())),
		)
	}
	return nil
}

// Stat obtains number of changed files.
func (instance *ShellGit) Stat(staged bool) (FileChanges, error) {
	args := []string{
		"diff",
		"--name-status",
	}
	if staged {
		args = append(args, "--cached")
	}
	cmd := exec.Command("git", args...)

	// cd into directory where notes are stored and run git there. Using 'git
	// -C <path>' is not an otion here because git fatals when called from
	// another repository.
	cmd.Dir = instance.WorktreePath
	out, err := cmd.Output()
	if err != nil {
		return FileChanges{}, errors.Join(err, errors.New("Cannot perform git diff"))
	}
	changes := FileChanges{}
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) < 1 {
			continue
		}
		switch line[0] {
		case 'A':
			changes.Added++
		case 'C':
			changes.Copied++
		case 'D':
			changes.Deleted++
		case 'M':
			changes.Modified++
		case 'R':
			changes.Renamed++
		}
	}
	return changes, nil
}
