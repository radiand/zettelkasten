package notes

import "testing"

import "github.com/stretchr/testify/assert"

func TestLoadNote(t *testing.T) {
	// GIVEN
	testCases := []struct {
		testName string
		body     string
	}{
		{"Fenced block at bottom", "Abcdef.\n```\n$curl -XGET localhost:8080\n```"},
		{"Fenced block on top", "```\nprint('hello')\n```\nAbcdef."},
		{"Fenced block with type", "```python\nprint('hello')\n```\nAbcdef."},
		{"Nothing", ""},
		{"Simple string", "My body is a cage"},
	}

	header := "```toml\n" +
		"title = \"NOTE_TITLE\"\n" +
		"timestamp = \"2024-01-01T01:00:00+01:00\"\n" +
		"uid = \"20240101T000000Z\"\n" +
		"tags = [\"lang:en\"]\n" +
		"referred_from = [\"20200101T000000Z\"]\n" +
		"refers_to = [\"20210101T000000Z\"]\n" +
		"```\n\n"

	for _, tc := range testCases {
		testFunc := func(t *testing.T) {
			file := header + tc.body

			// WHEN
			actual, err := UnmarshallNote(file)
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
				Body: tc.body,
			}

			assert.Equal(t, expected, actual)
		}

		t.Run(tc.testName, testFunc)
	}
}
