package notes

// INoteRepository defines interface for repositories - types returning Notes.
type INoteRepository interface {
	Get(uid string) (Note, error)
	Put(note Note) (string, error)
	List() ([]string, error)
}
