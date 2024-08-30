package note

import "errors"
import "regexp"
import "strings"

import "github.com/BurntSushi/toml"

// LoadNote unmarshalls Note from string. This function expects that Note's
// Header was marshalled to toml string, wrapped in ```toml``` fenced block, as
// commonly done in markdown.
func LoadNote(content string) (res Note, err error) {
	zkRe := regexp.MustCompile("(?s)```toml\n(?P<header>[^`]+)```\n*(?P<body>.*)\n?")
	matched := zkRe.FindStringSubmatch(string(content))
	headerRaw := matched[zkRe.SubexpIndex("header")]
	bodyRaw := strings.TrimSpace(matched[zkRe.SubexpIndex("body")])

	var header Header
	_, err = toml.Decode(headerRaw, &header)
	if err != nil {
		return Note{}, errors.Join(err, errors.New("Cannot unmarshall note"))
	}

	return Note{Header: header, Body: bodyRaw}, nil
}
