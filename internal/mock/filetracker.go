package mock

// FileTracker mocks the service to track already uploaded files.
type FileTracker struct {
	PutFn    func(path string, id string) error
	ExistFn  func(path string) bool
	DeleteFn func(path string) error
}

// Put invokes the mock implementation.
func (t *FileTracker) Put(path string, id string) error {
	return t.PutFn(path, id)
}

// Exist invokes the mock implementation.
func (t *FileTracker) Exist(path string) bool {
	return t.ExistFn(path)
}

// Delete invokes the mock implementation.
func (t *FileTracker) Delete(path string) error {
	return t.DeleteFn(path)
}
