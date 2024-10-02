package common

import "os"
import "strings"

// ExpandHomeDir expands $HOME dir variables with actual, full path to $HOME.
func ExpandHomeDir(text string) string {
	home, _ := os.UserHomeDir()
	tokens := []string{"~", "$HOME"}
	for _, token := range tokens {
		text = strings.Replace(text, token, home, 1)
	}
	return text
}

// Exists returns whether the given file or directory exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
