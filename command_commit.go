package main

import "fmt"
import "errors"
import "strings"

import "github.com/radiand/zettelkasten/internal/git"

// CmdCommit carries required params to run command.
type CmdCommit struct {
	zettelkastenDir string
	git             git.IGit
}

// Run performs git commit with all changes that happened in RootDir directory.
func (cmd CmdCommit) Run() error {
	err := cmd.git.Add(cmd.zettelkastenDir)
	if err != nil {
		return err
	}

	statuses, err := cmd.git.Status()
	if err != nil {
		return errors.Join(err, errors.New("Could not obtain git status"))
	}

	aggregated := countStaged(statuses)
	if !aggregated.any() {
		return nil
	}

	commitMsg := composeCommitMessage(aggregated)

	err = cmd.git.Commit(commitMsg)
	if err != nil {
		return err
	}
	return nil
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
