/*
Package note defines representation and methods of a single zettelkasten note
*/
package note

import "fmt"

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
	return fmt.Sprintf("```toml\n%s```\n\n%s\n", header, note.Body), nil
}

// Arrange enforces unified style of Notes. It modifies Note in place.
func (note *Note) Arrange() {
	note.Header.Arrange()
}

// NewNote creates new Note, dated now.
func NewNote() Note {
	return Note{Header: NewHeader(), Body: ""}
}
