package filetracker

import (
	"strconv"
	"strings"
	"time"
)

// TrackedFile represents a tracked file in the repository.
type TrackedFile struct {
	ModTime time.Time
	Hash    string
	ID      string
}

// NewTrackedFile returns a TrackedFile with the specified values
func NewTrackedFile(value string) TrackedFile {
	parts := strings.SplitN(value, "|", 3)

	modTime := time.Time{}
	hash := ""
	id := ""

	if len(parts) >= 2 {
		unixTime, err := strconv.ParseInt(parts[0], 10, 64)
		if err == nil {
			modTime = time.Unix(0, unixTime)
		}
		hash = parts[1]
	} else {
		hash = parts[0]
	}

	if len(parts) == 3 {
		id = parts[2]
	}

	return TrackedFile{
		Hash:    hash,
		ModTime: modTime,
		ID:      id,
	}
}

func (tf TrackedFile) String() string {
	if tf.ModTime.IsZero() {
		return tf.Hash
	} else {
		return strconv.FormatInt(tf.ModTime.UnixNano(), 10) + "|" + tf.Hash
	}
}
