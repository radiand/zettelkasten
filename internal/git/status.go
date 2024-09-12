package git

import "bytes"
import "fmt"

// Status is used as a type for enums.
type Status int8

// Status enum.
const (
	Added              Status = iota
	Copied             Status = iota
	Deleted            Status = iota
	Ignored            Status = iota
	Modified           Status = iota
	Renamed            Status = iota
	TypeChanged        Status = iota
	Unmodified         Status = iota
	Untracked          Status = iota
	UpdatedButUnmerged Status = iota
)

// FileStatus carries git status of a path.
type FileStatus struct {
	Path     string
	Staged   Status
	Unstaged Status
}

func charToStatus(char byte) Status {
	switch char {
	case 'A':
		return Added
	case 'C':
		return Copied
	case 'D':
		return Deleted
	case '!':
		return Ignored
	case 'M':
		return Modified
	case 'R':
		return Renamed
	case 'T':
		return TypeChanged
	case ' ':
		return Unmodified
	case '?':
		return Untracked
	case 'U':
		return UpdatedButUnmerged
	}
	panic("Unknown git status identifier '" + string(char) + "'")
}

func readGitStatusPorcelain(data []byte) ([]FileStatus, error) {
	files := []FileStatus{}
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(line) < 1 {
			continue
		}
		x := charToStatus(line[0])
		y := charToStatus(line[1])
		path := line[3:]
		if bytes.Contains(path, []byte("->")) {
			return []FileStatus{}, fmt.Errorf("Rename is not supported (invalid line: '%s')", line)
		}

		files = append(files, FileStatus{Path: string(path), Staged: x, Unstaged: y})
	}
	return files, nil
}
