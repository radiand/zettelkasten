package main

import "fmt"
import "errors"
import "strings"
import "time"
import "path"

import "github.com/radiand/zettelkasten/internal/git"

// CmdCommit carries required params to run command.
type CmdCommit struct {
	dirs       []string
	gitFactory func(workdir string) git.IGit
	nowtime    func() time.Time
	modtime    func(path string) (time.Time, error)
	cooldown   time.Duration
}

// Run performs git commit with all changes that happened in RootDir directory.
func (self CmdCommit) Run() error {
	for _, path := range self.dirs {
		err := self.run(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self CmdCommit) run(workdir string) error {
	gitHandler := self.gitFactory(workdir)
	var addErr error
	if self.cooldown > 0 {
		pathsToIgnore, err := self.filterPathsStillInCooldown(workdir)
		if err != nil {
			return err
		}
		pathsToIgnore = wrapWithIgnore(pathsToIgnore)
		pathsPassedToAdd := append([]string{workdir}, pathsToIgnore...)
		addErr = gitHandler.Add(pathsPassedToAdd...)
	} else {
		addErr = gitHandler.Add(workdir)
	}

	if addErr != nil {
		return addErr
	}

	statuses, err := gitHandler.Status()
	if err != nil {
		return errors.Join(err, errors.New("Could not obtain git status"))
	}

	aggregated := countStaged(statuses)
	if !aggregated.any() {
		return nil
	}

	commitMsg := composeCommitMessage(aggregated)

	err = gitHandler.Commit(commitMsg)
	if err != nil {
		return err
	}
	return nil
}

func (self CmdCommit) filterPathsStillInCooldown(workdir string) ([]string, error) {
	gitHandler := self.gitFactory(workdir)
	statuses, err := gitHandler.Status()
	if err != nil {
		return []string{}, errors.Join(err, errors.New("Could not obtain git status"))
	}

	gitRootDir, err := gitHandler.RootDir()
	if err != nil {
		return []string{}, errors.Join(err, errors.New("Could not obtain root dir"))
	}

	paths := []string{}
	now := self.nowtime()
	for _, status := range statuses {
		path := path.Join(gitRootDir, status.Path)
		modtime, err := self.modtime(path)
		if err != nil {
			return []string{}, errors.Join(err, fmt.Errorf("Could not get mod time of path %s", path))
		}
		if now.Sub(modtime) <= self.cooldown {
			paths = append(paths, status.Path)
		}
	}
	return paths, nil
}

func wrapWithIgnore(paths []string) []string {
	wrapped := []string{}
	for _, path := range paths {
		wrapped = append(wrapped, ":!"+path)
	}
	return wrapped
}

func composeCommitMessage(changes aggregation) string {
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
	populateWith(changes.added, "added")
	populateWith(changes.copied, "copied")
	populateWith(changes.deleted, "deleted")
	populateWith(changes.modified, "modified")
	populateWith(changes.renamed, "renamed")
	return "auto: " + strings.Join(changesStringified, ", ")
}

type aggregation struct {
	added    int
	copied   int
	deleted  int
	modified int
	renamed  int
}

func (self *aggregation) any() bool {
	a := self.added != 0
	c := self.copied != 0
	d := self.deleted != 0
	m := self.modified != 0
	r := self.renamed != 0
	if a || c || d || m || r {
		return true
	}
	return false
}

func countStaged(statuses []git.FileStatus) aggregation {
	aggr := aggregation{0, 0, 0, 0, 0}
	for _, st := range statuses {
		switch st.Staged {
		case git.Added:
			aggr.added++
		case git.Copied:
			aggr.copied++
		case git.Deleted:
			aggr.deleted++
		case git.Modified:
			aggr.modified++
		case git.Renamed:
			aggr.renamed++
		}
	}
	return aggr
}
