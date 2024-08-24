package note

// InMemoryNoteRepository is an implementation of INoteRepository interface and
// stores notes in memory. Use cases: debugging, testing.
type InMemoryNoteRepository struct {
	notes map[string]Note
}

// Get obtains Note.
func (repo *InMemoryNoteRepository) Get(uid string) (Note, error) {
	return repo.notes[uid], nil
}

// Put saves Note.
func (repo *InMemoryNoteRepository) Put(note Note) error {
	repo.notes[note.Header.Uid] = note
	return nil
}

// List obtains array of saved Notes' Uids.
func (repo *InMemoryNoteRepository) List() ([]string, error) {
	keys := []string{}
	for k := range repo.notes {
		keys = append(keys, k)
	}
	return keys, nil
}

// NewInMemoryNoteRepository creates new instance of the repository.
func NewInMemoryNoteRepository() *InMemoryNoteRepository {
	return &InMemoryNoteRepository{notes: make(map[string]Note)}
}
