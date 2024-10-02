package common

import "fmt"
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

// AskBool prompts for yes/no answer in console.
func AskBool(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var answer string
	fmt.Scanln(&answer)
	token := strings.ToLower(answer)
	return token == "y" || token == "yes"
}
