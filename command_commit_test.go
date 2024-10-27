package main

import "testing"
import "time"

import "github.com/radiand/zettelkasten/internal/git"
import "github.com/radiand/zettelkasten/internal/testutils"
import "github.com/stretchr/testify/assert"

func TestCommitWhenNoChanges(t *testing.T) {
	// GIVEN
	gitMock := git.NewMockGit()
	gitMock.StatusReturns.Enqueue([]git.FileStatus{})

	cmdCommit := CmdCommit{
		dirs:       []string{"/tmp"}, // Does not matter.
		gitFactory: func(string) git.IGit { return &gitMock },
	}

	// WHEN
	err := cmdCommit.Run()

	// THEN
	assert.Nil(t, err)
	assert.False(t, gitMock.CommitCapture.WasCalled)
}

func TestCommitChanges(t *testing.T) {
	// GIVEN
	gitMock := git.NewMockGit()
	gitMock.StatusReturns.Enqueue(
		[]git.FileStatus{
			{Path: "/tmp/a1.txt", Staged: git.Added, Unstaged: git.Unmodified},
			{Path: "/tmp/a2.txt", Staged: git.Added, Unstaged: git.Unmodified},
			{Path: "/tmp/m1.txt", Staged: git.Modified, Unstaged: git.Unmodified},
			{Path: "/tmp/d1.txt", Staged: git.Deleted, Unstaged: git.Unmodified},
		},
	)

	cmdCommit := CmdCommit{
		dirs:       []string{"/tmp"}, // Does not matter.
		gitFactory: func(string) git.IGit { return &gitMock },
	}

	// WHEN
	err := cmdCommit.Run()

	// THEN
	assert.Nil(t, err)
	assert.True(t, gitMock.CommitCapture.WasCalled)
	assert.Equal(t, "auto: 2 added, 1 deleted, 1 modified", gitMock.CommitCapture.CalledWith)
}

func TestCommitOldEnough(t *testing.T) {
	// GIVEN
	gitMock := git.NewMockGit()
	gitMock.RootDirReturns = "/virtual"

	// First git.Status call is before git.Add, so all paths are unstaged now.
	gitMock.StatusReturns.Enqueue(
		[]git.FileStatus{
			{Path: "zettelkasten/old.txt", Staged: git.Unmodified, Unstaged: git.Modified},
			{Path: "zettelkasten/new.txt", Staged: git.Unmodified, Unstaged: git.Modified},
		},
	)

	// Second git.Status call is after git.Add.
	gitMock.StatusReturns.Enqueue(
		[]git.FileStatus{
			{Path: "zettelkasten/old.txt", Staged: git.Modified, Unstaged: git.Unmodified},
			{Path: "zettelkasten/new.txt", Staged: git.Unmodified, Unstaged: git.Modified},
		},
	)

	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	pathModTimes := map[string]time.Time{
		"/virtual/zettelkasten/old.txt": t0,
		"/virtual/zettelkasten/new.txt": t0.Add(time.Second * 60),
	}

	cooldown, _ := time.ParseDuration("60s")
	cmdCommit := CmdCommit{
		dirs:       []string{"/virtual/zettelkasten"},
		gitFactory: func(string) git.IGit { return &gitMock },
		nowtime:    testutils.Then(t0.Add(time.Second * 61)),
		modtime:    testutils.TimeOfPath(pathModTimes),
		cooldown:   cooldown,
	}

	// WHEN
	err := cmdCommit.Run()

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []string{"/virtual/zettelkasten", ":!zettelkasten/new.txt"}, gitMock.AddCapture.CalledWith)
	assert.Equal(t, "auto: 1 modified", gitMock.CommitCapture.CalledWith)
}
