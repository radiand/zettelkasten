package notes

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
func (self *Header) ToToml() (res string, err error) {
	marshalled, err := toml.Marshal(self)
	if err != nil {
		return "", err
	}
	return string(marshalled), nil
}

// Arrange enforces unified style of Headers. It modifies Header in place.
func (self *Header) Arrange() {
	// All tags must be lowercase.
	for i := 0; i < len(self.Tags); i++ {
		self.Tags[i] = strings.ToLower(self.Tags[i])
	}
	sort.Sort(sort.StringSlice(self.Tags))
	sort.Sort(sort.StringSlice(self.ReferredFrom))
	sort.Sort(sort.StringSlice(self.RefersTo))
}

// NewHeader creates new Header.
func NewHeader(when time.Time) Header {
	uid := when.UTC().Format("20060102T150405Z")

	return Header{
		Title:        "",
		Timestamp:    when.Format("2006-01-02T15:04:05-07:00"),
		Uid:          uid,
		Tags:         []string{},
		ReferredFrom: []string{},
		RefersTo:     []string{},
	}
}

// GetUidRegexp creates regexp matching Note Uid, i.e. filenames and references
// of other Notes within Note's body.
func GetUidRegexp() *regexp.Regexp {
	uidPat := `\d{4}\d{2}\d{2}T\d{2}\d{2}\d{2}Z`
	return regexp.MustCompile(uidPat)
}
