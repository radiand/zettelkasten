package common

import "fmt"

type disjoinable interface {
	Unwrap() []error
}

type unwrappable interface {
	Unwrap() error
}

// Disjoin recursively unwraps joined errors.
func Disjoin(e error) []error {
	if joined, ok := e.(disjoinable); ok {
		var unwrapped []error
		for _, err := range joined.Unwrap() {
			unwrapped = append(unwrapped, Disjoin(err)...)
		}
		return unwrapped
	}
	return []error{e}
}

// LastError extracts last error from joined one.
func LastError(e error) error {
	if joined, ok := e.(disjoinable); ok {
		if len(joined.Unwrap()) > 0 {
			return LastError(joined.Unwrap()[len(joined.Unwrap())-1])
		}
	}
	return e
}

// FmtErrors pretty prints errors that were wrapped by errors.Join.
func FmtErrors(e error) string {
	disjoint := Disjoin(e)
	if len(disjoint) > 1 {
		formatted := "Errors in order from innermost:"
		for idx, err := range disjoint {
			formatted = formatted + fmt.Sprintf("\n  %d. (%T) %s", idx+1, err, err.Error())
		}
		return formatted
	}
	return e.Error()
}
