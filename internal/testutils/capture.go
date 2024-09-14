package testutils

// Capture is used to store info whether method was called.
type Capture[T any] struct {
	WasCalled  bool
	CalledWith T
}
