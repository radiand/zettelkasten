package main

import "testing"
import "time"

import "github.com/radiand/zettelkasten/internal/git"
import "github.com/radiand/zettelkasten/internal/testutils"
import "github.com/stretchr/testify/assert"

type GitMock struct {
	statusReturns testutils.Cycle[[]git.FileStatus]
	addCapture    testutils.Capture[[]string]
	commitCapture testutils.Capture[string]
}

func NewGitMock() GitMock {
	return GitMock{
		statusReturns: testutils.NewCycle[[]git.FileStatus](),
		addCapture:    testutils.Capture[[]string]{},
		commitCapture: testutils.Capture[string]{
			WasCalled:  false,
			CalledWith: "",
		},
	}
}

func (self *GitMock) Add(paths ...string) error {
	self.addCapture.WasCalled = true
	for _, path := range paths {
		self.addCapture.CalledWith = append(self.addCapture.CalledWith, path)
	}
	return nil
}

func (self *GitMock) Commit(message string) error {
	self.commitCapture.WasCalled = true
	self.commitCapture.CalledWith = message
	return nil
}

func (self *GitMock) Status() ([]git.FileStatus, error) {
	return self.statusReturns.Next(), nil
}

func TestCommitWhenNoChanges(t *testing.T) {
	// GIVEN
	gitMock := NewGitMock()
	gitMock.statusReturns.Enqueue([]git.FileStatus{})

	cmdCommit := CmdCommit{
		zettelkastenDir: "/tmp", // Does not matter.
		git:             &gitMock,
	}

	// WHEN
	err := cmdCommit.Run()

	// THEN
	assert.Nil(t, err)
	assert.False(t, gitMock.commitCapture.WasCalled)
}

func TestCommitChanges(t *testing.T) {
	// GIVEN
	gitMock := NewGitMock()
	gitMock.statusReturns.Enqueue(
		[]git.FileStatus{
			{Path: "/tmp/a1.txt", Staged: git.Added, Unstaged: git.Unmodified},
			{Path: "/tmp/a2.txt", Staged: git.Added, Unstaged: git.Unmodified},
			{Path: "/tmp/m1.txt", Staged: git.Modified, Unstaged: git.Unmodified},
			{Path: "/tmp/d1.txt", Staged: git.Deleted, Unstaged: git.Unmodified},
		},
	)

	cmdCommit := CmdCommit{
		zettelkastenDir: "/tmp", // Does not matter.
		git:             &gitMock,
	}

	// WHEN
	err := cmdCommit.Run()

	// THEN
	assert.Nil(t, err)
	assert.True(t, gitMock.commitCapture.WasCalled)
	assert.Equal(t, "auto: 2 added, 1 deleted, 1 modified", gitMock.commitCapture.CalledWith)
}

func TestCommitOldEnough(t *testing.T) {
	// GIVEN
	gitMock := NewGitMock()

	// First git.Status call is before git.Add, so all paths are unstaged now.
	gitMock.statusReturns.Enqueue(
		[]git.FileStatus{
			{Path: "old.txt", Staged: git.Unmodified, Unstaged: git.Modified},
			{Path: "new.txt", Staged: git.Unmodified, Unstaged: git.Modified},
		},
	)

	// Second git.Status call is after git.Add.
	gitMock.statusReturns.Enqueue(
		[]git.FileStatus{
			{Path: "old.txt", Staged: git.Modified, Unstaged: git.Unmodified},
			{Path: "new.txt", Staged: git.Unmodified, Unstaged: git.Modified},
		},
	)

	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	pathModTimes := map[string]time.Time{
		"/virtual/old.txt": t0,
		"/virtual/new.txt": t0.Add(time.Second * 60),
	}
	cmdCommit := CmdCommit{
		zettelkastenDir: "/virtual",
		git:             &gitMock,
		nowtime:         testutils.Then(t0.Add(time.Second * 61)),
		modtime:         testutils.TimeOfPath(pathModTimes),
		olderThanSec:    60,
	}

	// WHEN
	err := cmdCommit.Run()

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []string{"/virtual/old.txt"}, gitMock.addCapture.CalledWith)
	assert.Equal(t, "auto: 1 modified", gitMock.commitCapture.CalledWith)
}
