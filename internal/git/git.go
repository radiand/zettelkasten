/*
Package git provides just enough integration with git.
*/
package git

// IGit interface provides version control functionalities with git.
type IGit interface {
	Add(paths ...string) error
	Commit(message string) error
	Status() ([]FileStatus, error)
	RootDir() (string, error)
}
