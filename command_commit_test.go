package main

import "testing"

import "github.com/radiand/zettelkasten/internal/git"
import "github.com/stretchr/testify/assert"

type Called[T any] struct {
	WasCalled  bool
	CalledWith T
}

type GitMock struct {
	statusReturn []git.FileStatus
	commitMethod Called[string]
}

func NewGitMock() GitMock {
	return GitMock{
		statusReturn: []git.FileStatus{},
		commitMethod: Called[string]{
			WasCalled:  false,
			CalledWith: "",
		},
	}
}

func (self *GitMock) Add(_ ...string) error {
	return nil
}

func (self *GitMock) Commit(message string) error {
	self.commitMethod.WasCalled = true
	self.commitMethod.CalledWith = message
	return nil
}

func (self *GitMock) Status() ([]git.FileStatus, error) {
	return self.statusReturn, nil
}

func TestCommitWhenNoChanges(t *testing.T) {
	// GIVEN
	gitMock := NewGitMock()
	gitMock.statusReturn = []git.FileStatus{}

	cmdCommit := CmdCommit{
		zettelkastenDir: "/tmp", // Does not matter.
		git:             &gitMock,
	}

	// WHEN
	err := cmdCommit.Run()

	// THEN
	assert.Nil(t, err)
	assert.False(t, gitMock.commitMethod.WasCalled)
}

func TestCommitChanges(t *testing.T) {
	// GIVEN
	gitMock := NewGitMock()
	gitMock.statusReturn = []git.FileStatus{
		{Path: "/tmp/a1.txt", Staged: git.Added, Unstaged: git.Unmodified},
		{Path: "/tmp/a2.txt", Staged: git.Added, Unstaged: git.Unmodified},
		{Path: "/tmp/m1.txt", Staged: git.Modified, Unstaged: git.Unmodified},
		{Path: "/tmp/d1.txt", Staged: git.Deleted, Unstaged: git.Unmodified},
	}

	cmdCommit := CmdCommit{
		zettelkastenDir: "/tmp", // Does not matter.
		git:             &gitMock,
	}

	// WHEN
	err := cmdCommit.Run()

	// THEN
	assert.Nil(t, err)
	assert.True(t, gitMock.commitMethod.WasCalled)
	assert.Equal(t, "auto: 2 added, 1 deleted, 1 modified", gitMock.commitMethod.CalledWith)
}
