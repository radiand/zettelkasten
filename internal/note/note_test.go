package note

import "testing"

import "github.com/stretchr/testify/assert"

func TestLoadNote(t *testing.T) {
	// GIVEN
	file := "```toml\n" +
		"title = \"NOTE_TITLE\"\n" +
		"timestamp = \"2024-01-01T01:00:00+01:00\"\n" +
		"uid = \"20240101T000000Z\"\n" +
		"tags = [\"lang:en\"]\n" +
		"referred_from = [\"20200101T000000Z\"]\n" +
		"refers_to = [\"20210101T000000Z\"]\n" +
		"```\n\n" +
		"My body is a cage"

	// WHEN
	actual, err := LoadNote(file)
	assert.Nil(t, err)

	// THEN
	expected := Note{
		Header: Header{
			Title:        "NOTE_TITLE",
			Timestamp:    "2024-01-01T01:00:00+01:00",
			Uid:          "20240101T000000Z",
			Tags:         []string{"lang:en"},
			ReferredFrom: []string{"20200101T000000Z"},
			RefersTo:     []string{"20210101T000000Z"},
		},
		Body: "My body is a cage",
	}

	assert.Equal(t, expected, actual)
}

func TestSaveNote(t *testing.T) {
	// GIVEN
	given := Note{
		Header: Header{
			Title:        "NOTE_TITLE",
			Timestamp:    "2024-01-01T01:00:00+01:00",
			Uid:          "20240101T000000Z",
			Tags:         []string{"lang:en"},
			ReferredFrom: []string{"20200101T000000Z"},
			RefersTo:     []string{"20210101T000000Z"},
		},
		Body: "My body is a cage",
	}

	expected := "```toml\n" +
		"title = \"NOTE_TITLE\"\n" +
		"timestamp = \"2024-01-01T01:00:00+01:00\"\n" +
		"uid = \"20240101T000000Z\"\n" +
		"tags = [\"lang:en\"]\n" +
		"referred_from = [\"20200101T000000Z\"]\n" +
		"refers_to = [\"20210101T000000Z\"]\n" +
		"```\n\n" +
		"My body is a cage"

	// WHEN
	actual, _ := given.ToToml()

	// THEN
	assert.Equal(t, expected, actual)
}

func TestNewNote(t *testing.T) {
	// WHEN
	actual := NewNote()

	// THEN
	assert.Equal(t, actual.Header.Title, "")
	assert.NotEqual(t, actual.Header.Uid, "")
	assert.NotEqual(t, actual.Header.Timestamp, "")
	assert.Equal(t, actual.Header.Tags, []string{})
	assert.Equal(t, actual.Header.ReferredFrom, []string{})
	assert.Equal(t, actual.Header.RefersTo, []string{})
}

func TestArrangeNote(t *testing.T) {
	// GIVEN
	note := NewNote()
	note.Header.Tags = []string{"b", "C", "a"}
	note.Header.ReferredFrom = []string{"20010101T010101Z", "19700101T010101Z"}
	note.Header.RefersTo = []string{"20020202T020202Z", "19700202T020202Z"}

	// WHEN
	note.Arrange()

	// THEN
	expectedTags := []string{"a", "b", "c"}
	expectedReferredFrom := []string{"19700101T010101Z", "20010101T010101Z"}
	expectedRefersTo := []string{"19700202T020202Z", "20020202T020202Z"}
	assert.Equal(t, note.Header.Tags, expectedTags)
	assert.Equal(t, note.Header.ReferredFrom, expectedReferredFrom)
	assert.Equal(t, note.Header.RefersTo, expectedRefersTo)
}
