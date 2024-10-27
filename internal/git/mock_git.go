package git

import "github.com/radiand/zettelkasten/internal/testutils"

// MockGit replaces IGit in tests.
type MockGit struct {
	StatusReturns  testutils.Cycle[[]FileStatus]
	AddCapture     testutils.Capture[[]string]
	CommitCapture  testutils.Capture[string]
	RootDirReturns string
}

// NewMockGit creates new, empty instance of MockGit.
func NewMockGit() MockGit {
	return MockGit{
		StatusReturns: testutils.NewCycle[[]FileStatus](),
		AddCapture:    testutils.Capture[[]string]{},
		CommitCapture: testutils.Capture[string]{
			WasCalled:  false,
			CalledWith: "",
		},
		RootDirReturns: "/root",
	}
}

// Add captures calls to IGit.Add().
func (self *MockGit) Add(paths ...string) error {
	self.AddCapture.WasCalled = true
	for _, path := range paths {
		self.AddCapture.CalledWith = append(self.AddCapture.CalledWith, path)
	}
	return nil
}

// Commit captures calls to IGit.Commit().
func (self *MockGit) Commit(message string) error {
	self.CommitCapture.WasCalled = true
	self.CommitCapture.CalledWith = message
	return nil
}

// Status mocks IGit.Status() and returns consecutive values of self.StatusReturns.
func (self *MockGit) Status() ([]FileStatus, error) {
	return self.StatusReturns.Next(), nil
}

// RootDir mocks IGit.Status() and constantly returns value of self.RootDirReturns.
func (self *MockGit) RootDir() (string, error) {
	return self.RootDirReturns, nil
}
