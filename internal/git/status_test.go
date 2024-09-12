package git

import "testing"

import "github.com/stretchr/testify/assert"

func TestReadingPorcelainV1(t *testing.T) {
	cmdOut := "MM a.txt\n" +
		"A  b.txt\n" +
		"?? c.txt\n"

	actual, err := readGitStatusPorcelain([]byte(cmdOut))
	assert.Nil(t, err)

	expected := []FileStatus{
		{Staged: Modified, Unstaged: Modified, Path: "a.txt"},
		{Staged: Added, Unstaged: Unmodified, Path: "b.txt"},
		{Staged: Untracked, Unstaged: Untracked, Path: "c.txt"},
	}
	assert.Equal(t, expected, actual)
}
