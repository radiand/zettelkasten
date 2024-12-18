package notes

// InMemoryNoteRepository is an implementation of INoteRepository interface and
// stores notes in memory. Use cases: debugging, testing.
type InMemoryNoteRepository struct {
	notes map[string]Note
}

// Get obtains Note.
func (self *InMemoryNoteRepository) Get(uid string) (Note, error) {
	return self.notes[uid], nil
}

// Put saves Note.
func (self *InMemoryNoteRepository) Put(note Note) (string, error) {
	self.notes[note.Header.Uid] = note
	return note.Header.Uid + ".md", nil
}

// List obtains array of saved Notes' Uids.
func (self *InMemoryNoteRepository) List() ([]string, error) {
	keys := []string{}
	for k := range self.notes {
		keys = append(keys, k)
	}
	return keys, nil
}

// NewInMemoryNoteRepository creates new instance of the repository.
func NewInMemoryNoteRepository() *InMemoryNoteRepository {
	return &InMemoryNoteRepository{notes: make(map[string]Note)}
}
