package notes

import "fmt"
import "testing"
import "time"

import "github.com/stretchr/testify/assert"

func TestFindUids(t *testing.T) {
	// GIVEN
	uid1 := "20240101T010101Z"
	uid2 := "20240202T020202Z"
	text := fmt.Sprintf(
		"This note refers to [[%s]].\nSomewhere later it refers to [this](%s)",
		uid1,
		uid2,
	)

	// WHEN
	uids := FindUids(text)

	// THEN
	assert.Equal(t, uids, []string{uid1, uid2})
}

func TestReferences(t *testing.T) {
	// GIVEN
	note1 := NewNote(time.Date(1991, 1, 1, 1, 1, 1, 0, time.UTC))
	note1uid := note1.Header.Uid
	uid11 := "20240101T010101Z"
	note1.Body = fmt.Sprintf("Refers to [[%s]]", uid11)

	note2 := NewNote(time.Date(1992, 2, 2, 2, 2, 2, 0, time.UTC))
	note2uid := note2.Header.Uid
	uid21 := "20240202T020202Z"
	uid22 := "20240303T030303Z"
	note2.Body = fmt.Sprintf("Refers to [[%s]] and [[%s]]", uid21, uid22)

	repository := NewInMemoryNoteRepository()
	repository.Put(note1)
	repository.Put(note2)

	// WHEN
	actual := FindReferences(repository)
	expected := ReferenceMap{
		note1uid: []string{uid11},
		note2uid: []string{uid21, uid22},
	}

	// THEN
	assert.Equal(t, expected, actual)

	// WHEN
	actual = ReverseReferences(actual)
	expected = ReferenceMap{
		uid11: []string{note1uid},
		uid21: []string{note2uid},
		uid22: []string{note2uid},
	}

	// THEN
	assert.Equal(t, expected, actual)
}

// TestReferencesWhenNotesAreNotRelated verifies if notes that do not link
// anywhere do not appear in the references map.
func TestReferencesWhenNotesAreNotRelated(t *testing.T) {
	// GIVEN
	note1 := NewNote(time.Date(1991, 1, 1, 1, 1, 1, 0, time.UTC))
	note2 := NewNote(time.Date(1992, 2, 2, 2, 2, 2, 0, time.UTC))

	repository := NewInMemoryNoteRepository()
	repository.Put(note1)
	repository.Put(note2)

	// WHEN
	actual := FindReferences(repository)

	// THEN
	expected := make(ReferenceMap)
	assert.Equal(t, expected, actual)
}

func TestLinkNotes(t *testing.T) {
	// GIVEN
	note1 := NewNote(time.Date(1991, 1, 1, 1, 1, 1, 0, time.UTC))
	note1uid := note1.Header.Uid
	uid11 := "20240101T010101Z"
	note1.Body = fmt.Sprintf("Refers to [[%s]]", uid11)

	note2 := NewNote(time.Date(1992, 2, 2, 2, 2, 2, 0, time.UTC))
	note2uid := note2.Header.Uid
	uid21 := "20240202T020202Z"
	note2.Body = fmt.Sprintf("Refers to [[%s]] and [[%s]]", uid21, note1uid)

	repository := NewInMemoryNoteRepository()
	repository.Put(note1)
	repository.Put(note2)

	// WHEN
	err := LinkNotes(repository)
	note1, _ = repository.Get(note1uid)
	note2, _ = repository.Get(note2uid)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []string{uid11}, note1.Header.RefersTo)
	assert.Equal(t, []string{note2uid}, note1.Header.ReferredFrom)
	assert.Equal(t, []string{note1uid, uid21}, note2.Header.RefersTo)
}
