package common

import "fmt"

// FmtErrors pretty prints errors that were wrapped by errors.Join.
func FmtErrors(wrapped error) string {
	if errs, ok := wrapped.(unwrappable); ok {
		formatted := "Errors in order from innermost:"
		for idx, err := range errs.Unwrap() {
			formatted = formatted + fmt.Sprintf("\n  %d. (%T) %s", idx+1, err, err.Error())
		}
		return formatted
	}
	return wrapped.Error()
}

type unwrappable interface {
	Unwrap() []error
}
