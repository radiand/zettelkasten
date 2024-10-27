package git

import "errors"
import "fmt"
import "os/exec"
import "strings"

// ShellGit is a Git interface implementation based on spawning shell process.
type ShellGit struct {
	WorktreePath string
}

// Add performs file staging.
func (instance *ShellGit) Add(paths ...string) error {
	args := []string{"-C", instance.WorktreePath, "add"}
	args = append(args, paths...)
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Join(err, fmt.Errorf("git add failed due to: %s", string(out)))
	}
	return nil
}

// Commit performs git commit with custom message.
func (instance *ShellGit) Commit(message string) error {
	cmd := exec.Command(
		"git", "-C", instance.WorktreePath, "commit", "-m", message,
	)
	err := cmd.Run()
	if err != nil {
		return errors.Join(
			err,
			fmt.Errorf("git commit failed due to: %s", fmtExitError(err)),
		)
	}
	return nil
}

// Status obtains git statuses of all paths in working directory. Note: this is
// "just enough" implementation and does not support many operations that could
// be performed in git repo.
func (instance *ShellGit) Status() ([]FileStatus, error) {
	cmd := exec.Command(
		"git",
		"-C",
		instance.WorktreePath,
		"status",
		"--porcelain=1",
	)
	out, err := cmd.Output()
	if err != nil {
		return []FileStatus{}, errors.Join(err, errors.New("git status failed"))
	}
	statuses, err := readGitStatusPorcelain(out)
	return statuses, err
}

// RootDir returns absolute path of repository root (typically where .git
// folder resides). Root directory can be used to join with paths from Status()
// to obtain absolute paths of files listed there.
func (instance *ShellGit) RootDir() (string, error) {
	cmd := exec.Command(
		"git",
		"-C",
		instance.WorktreePath,
		"rev-parse",
		"--show-toplevel",
	)
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Join(err, errors.New("git rev-parse failed"))
	}
	return strings.TrimSpace(string(out)), nil
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
