package upload

// UploadFolderJob represents a job to upload all photos from the specified folder
type UploadFolderJob struct {
	FileTracker FileTracker

	SourceFolder string
	CreateAlbums string
	Filter       FileFilterer
}

// FileTracker represents a service to track already uploaded files.
type FileTracker interface {
	Put(file string, id string) error
	Exist(file string) bool
	Delete(file string) error
}

// FileFilterer represents a way to implement include/exclude files filtering.
type FileFilterer interface {
	IsAllowed(path string) bool
	IsExcluded(path string) bool
}
