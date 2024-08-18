/*
Package note defines representation and methods of a single zettelkasten note
*/
package note

import "fmt"
import "os"
import "regexp"
import "slices"
import "sort"
import "strings"
import "time"

import "github.com/BurntSushi/toml"

// Header is a metadata put on top of the note. It is marshalled as a toml
// block.
type Header struct {
	Title        string   `toml:"title"`
	Timestamp    string   `toml:"timestamp"`
	Uid          string   `toml:"uid"` // revive:disable
	Tags         []string `toml:"tags"`
	ReferredFrom []string `toml:"referred_from"`
	RefersTo     []string `toml:"refers_to"`
}

// Equal checks equality of two Headers, i.e. same values and same order of
// them.
func (lhs *Header) Equal(rhs Header) bool {
	titlesEq := lhs.Title == rhs.Title
	timestampEq := lhs.Timestamp == rhs.Timestamp
	uidEq := lhs.Uid == rhs.Uid
	tagsEq := slices.Equal(lhs.Tags, rhs.Tags)
	refFromEq := slices.Equal(lhs.ReferredFrom, rhs.ReferredFrom)
	refToEq := slices.Equal(lhs.RefersTo, rhs.RefersTo)
	return titlesEq && timestampEq && uidEq && tagsEq && refFromEq && refToEq
}

// ToToml marshalls Header.
func (header *Header) ToToml() (res string, err error) {
	marshalled, err := toml.Marshal(header)
	if err != nil {
		return "", err
	}
	return string(marshalled), nil
}

// Arrange enforces unified style of Headers. It modifies Header in place.
func (header *Header) Arrange() {
	// All tags must be lowercase.
	for i := 0; i < len(header.Tags); i++ {
		header.Tags[i] = strings.ToLower(header.Tags[i])
	}
	sort.Sort(sort.StringSlice(header.Tags))
	sort.Sort(sort.StringSlice(header.ReferredFrom))
	sort.Sort(sort.StringSlice(header.RefersTo))
}

// Note is a single zettelkasten note.
type Note struct {
	Header Header
	Body   string
}

// Equal checks equality of two Notes.
func (lhs *Note) Equal(rhs Note) bool {
	headerEq := lhs.Header.Equal(rhs.Header)
	bodyEq := lhs.Body == rhs.Body
	return headerEq && bodyEq
}

// ToToml marshalls Note.
func (note *Note) ToToml() (res string, err error) {
	header, err := note.Header.ToToml()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("```toml\n%s```\n\n%s", header, note.Body), nil
}

// Arrange enforces unified style of Notes. It modifies Note in place.
func (note *Note) Arrange() {
	note.Header.Arrange()
}

// NewNote creates new Note, dated now.
func NewNote() Note {
	now := time.Now()
	uid := now.UTC().Format("20060102T150405Z")

	header := Header{
		Title:        "",
		Timestamp:    now.Format("2006-01-02T15:04:05-07:00"),
		Uid:          uid,
		Tags:         []string{},
		ReferredFrom: []string{},
		RefersTo:     []string{},
	}

	return Note{Header: header, Body: ""}
}

// LoadNote unmarshalls Note from string. This function expects that Note's
// Header was marshalled to toml string, wrapped in ```toml``` fenced block, as
// commonly done in markdown.
func LoadNote(content string) (res Note, err error) {
	zkRe := regexp.MustCompile("(?s)```toml\n(?P<header>.*)```\n*(?P<body>.*)\n?")
	matched := zkRe.FindStringSubmatch(string(content))
	headerRaw := matched[zkRe.SubexpIndex("header")]
	bodyRaw := strings.TrimSpace(matched[zkRe.SubexpIndex("body")])

	var header Header
	_, err = toml.Decode(headerRaw, &header)
	if err != nil {
		return Note{}, err
	}

	return Note{Header: header, Body: bodyRaw}, nil
}

// LoadNoteFromFile unmarshalls Note from file pointed by a given path
// argument.
func LoadNoteFromFile(path string) (res Note, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Note{}, err
	}
	return LoadNote(string(content))
}
