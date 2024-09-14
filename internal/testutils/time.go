package testutils

import "fmt"
import "time"

// Then replaces time.Now().UTC() in tests.
func Then(frozen time.Time) func() time.Time {
	return func() time.Time {
		return frozen
	}
}

// TimeOfPath replaces ModificationTime(path) in tests.
func TimeOfPath(times map[string]time.Time) func(string) (time.Time, error) {
	return func(path string) (time.Time, error) {
		value, ok := times[path]
		if ok {
			return value, nil
		}
		panic(fmt.Sprintf("No modtime specified for path: %s", path))
	}
}
