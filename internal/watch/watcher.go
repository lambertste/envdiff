package watch

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// FileState holds the last known checksum and modification time of a file.
type FileState struct {
	Path    string
	Checksum string
	ModTime  time.Time
}

// ChangeEvent describes a detected change to a watched file.
type ChangeEvent struct {
	Path    string
	OldHash string
	NewHash string
	At      time.Time
}

// Watch polls the given file paths at the given interval, sending a
// ChangeEvent to the returned channel whenever a file's content changes.
// The caller must close done to stop watching.
func Watch(paths []string, interval time.Duration, done <-chan struct{}) (<-chan ChangeEvent, error) {
	states := make(map[string]FileState, len(paths))
	for _, p := range paths {
		state, err := checksum(p)
		if err != nil {
			return nil, fmt.Errorf("watch: initial stat %q: %w", p, err)
		}
		states[p] = state
	}

	ch := make(chan ChangeEvent, 8)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for _, p := range paths {
					newState, err := checksum(p)
					if err != nil {
						continue
					}
					old := states[p]
					if newState.Checksum != old.Checksum {
						ch <- ChangeEvent{
							Path:    p,
							OldHash: old.Checksum,
							NewHash: newState.Checksum,
							At:      newState.ModTime,
						}
						states[p] = newState
					}
				}
			}
		}
	}()
	return ch, nil
}

func checksum(path string) (FileState, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileState{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return FileState{}, err
	}
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return FileState{}, err
	}
	return FileState{
		Path:     path,
		Checksum: fmt.Sprintf("%x", h.Sum(nil)),
		ModTime:  info.ModTime(),
	}, nil
}
