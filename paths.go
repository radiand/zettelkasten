package main

import "os"
import "strings"

func expandHomeDir(text string) string {
	home, _ := os.UserHomeDir()
	tokens := []string{"~", "$HOME"}
	for _, token := range tokens {
		text = strings.Replace(text, token, home, 1)
	}
	return text
}
