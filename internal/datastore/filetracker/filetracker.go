package filetracker

import (
	"fmt"
	"os"
	"time"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

var (
	// ErrItemNotFound is the expected error if the item is not found.
	ErrItemNotFound = fmt.Errorf("item was not found")
)

// FileTracker allows to track already uploaded files in a repository.
type FileTracker struct {
	repo Repository

	// Hasher allows to change the way that hashes are calculated. Uses xxHash32Hasher{} by default.
	// Useful for testing.
	Hasher Hasher

	logger log.Logger
}

// Hasher is a Hasher to get the value of the file.
type Hasher interface {
	Hash(file string) (string, error)
}

// Repository is the repository where to track already uploaded files.
type Repository interface {
	// Get It returns ErrItemNotFound if the repo does not contains the key.
	Get(key string) (TrackedFile, error)
	Put(key string, item TrackedFile) error
	Delete(key string) error
	Close() error
}

// New returns a FileTracker using specified repo.
func New(r Repository) *FileTracker {
	return &FileTracker{
		repo:   r,
		Hasher: xxHash32Hasher{},
		logger: log.GetInstance(),
	}
}

// Put marks a file as already uploaded to prevent re-uploads.
func (ft FileTracker) Put(file string, id string) error {
	modTime := time.Now()
	hash := file

	if id == "" {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}

		modTime = fileInfo.ModTime()
		fileHash, err := ft.Hasher.Hash(file)
		if err != nil {
			return err
		}

		hash = fileHash
	}

	item := TrackedFile{
		ModTime: modTime,
		Hash:    hash,
		ID:      id,
	}

	return ft.repo.Put(file, item)
}

func (ft FileTracker) Get(file string) (TrackedFile, error) {
	// Get returns ErrItemNotFound if the repo does not contains the key.
	item, err := ft.repo.Get(file)
	if err != nil {
		return TrackedFile{}, err
	}

	return item, nil
}

// Exist checks if the file was already uploaded.
// Exist compares the last modification time of the file against the one in the repository.
// Last time modification comparison tries to reduce the number of times where the hash comparison
// is needed.
// In case that last modification time has changed (or it doesn't exist - retro compatibility),
// it compares a hash of the content of the file against the one in the repository.
func (ft FileTracker) Exist(file string) bool {
	// Get returns ErrItemNotFound if the repo does not contains the key.
	item, err := ft.repo.Get(file)
	if err != nil {
		return false
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		ft.logger.Debugf("Error retrieving file info for '%s' (%s).", file, err)
		return false
	}

	if item.ModTime.Equal(fileInfo.ModTime()) {
		ft.logger.Debugf("File modification time has not changed for '%s'.", file)
		return true
	}

	hash, err := ft.Hasher.Hash(file)
	if err != nil {
		return false
	}

	// checks if the file is the same (equal value)
	if item.Hash == hash {
		ft.logger.Debugf("File hash has not changed for '%s'.", file)

		// updates file marker with mtime to speed up comparison on next run
		item.ModTime = fileInfo.ModTime()
		if err = ft.repo.Put(file, item); err != nil {
			ft.logger.Debugf("Error updating marker for '%s' with modification time (%s).", file, err)
		}

		return true
	}

	return false
}

// Delete un-marks a file as already uploaded.
func (ft FileTracker) Delete(file string) error {
	return ft.repo.Delete(file)
}

// Close closes the file tracker repository.
// No operation could be done after that.
func (ft FileTracker) Close() error {
	return ft.repo.Close()
}
