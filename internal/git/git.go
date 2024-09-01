/*
Package git provides just enough integration with git.
*/
package git

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
	Add(paths ...string) error
	Commit() error
	Stat(staged bool) (FileChanges, error)
}

// ShellGit is a Git interface implementation based on spawning shell process.
type ShellGit struct {
	WorktreePath string
}

// Add performs file staging.
func (instance *ShellGit) Add(paths ...string) error {
	cmd := exec.Command(
		"git", "-C", instance.WorktreePath, "add", strings.Join(paths, " "),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Join(err, fmt.Errorf("git add failed due to: %s", string(out)))
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

	cmd := exec.Command("git", "-C", instance.WorktreePath, "commit", "-m", commitmsg)

	err = cmd.Run()
	if err != nil {
		return errors.Join(
			err,
			fmt.Errorf("git commit failed due to: %s", fmtExitError(err)),
		)
	}
	return nil
}

// Stat obtains number of changed files.
func (instance *ShellGit) Stat(staged bool) (FileChanges, error) {
	args := []string{
		"-C",
		instance.WorktreePath,
		"diff",
		"--name-status",
	}
	if staged {
		args = append(args, "--cached")
	}
	cmd := exec.Command("git", args...)

	out, err := cmd.Output()
	if err != nil {
		return FileChanges{}, errors.Join(err, errors.New("git diff failed"))
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

func fmtExitError(err error) string {
	if exitErr, ok := err.(*exec.ExitError); ok {
		if len(exitErr.Stderr) > 0 {
			return fmt.Sprintf(
				"%s (stderr: %s)",
				exitErr.Error(),
				string(exitErr.Stderr),
			)
		}
	}
	return err.Error()
}
