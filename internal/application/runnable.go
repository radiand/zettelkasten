package application

// Runnable is a common interface for Commands and Queries.
type Runnable interface {
	// Run returns the result of procedure or error. Treat (string, error) pair
	// as stdout and stderr.
	Run() (string, error)
}
