// Package store provides persistence for cron job log entries.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/user/cronlog/internal/logentry"
)

// Store persists log entries to a JSON-lines file.
type Store struct {
	mu   sync.Mutex
	path string
}

// New creates a new Store that writes to the given file path.
// The directory is created if it does not exist.
func New(path string) (*Store, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("store: create directory: %w", err)
	}
	return &Store{path: path}, nil
}

// Append writes a finished log entry to the store as a JSON line.
func (s *Store) Append(e *logentry.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("store: open file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("store: encode entry: %w", err)
	}
	return nil
}

// ReadAll reads all log entries from the store.
func (s *Store) ReadAll() ([]*logentry.Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("store: open file: %w", err)
	}
	defer f.Close()

	var entries []*logentry.Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e logentry.Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("store: decode entry: %w", err)
		}
		entries = append(entries, &e)
	}
	return entries, nil
}
