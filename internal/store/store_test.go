package store_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/cronlog/internal/logentry"
	"github.com/user/cronlog/internal/store"
)

func newTestEntry(job string, success bool) *logentry.Entry {
	e := logentry.New(job)
	e.Finish(success, "output line")
	return e
}

func TestAppendAndReadAll(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logs", "cronlog.jsonl")

	s, err := store.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	e1 := newTestEntry("backup", true)
	e2 := newTestEntry("cleanup", false)

	if err := s.Append(e1); err != nil {
		t.Fatalf("Append e1: %v", err)
	}
	if err := s.Append(e2); err != nil {
		t.Fatalf("Append e2: %v", err)
	}

	entries, err := s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Job != "backup" {
		t.Errorf("expected job 'backup', got %q", entries[0].Job)
	}
	if entries[1].Job != "cleanup" {
		t.Errorf("expected job 'cleanup', got %q", entries[1].Job)
	}
}

func TestReadAll_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cronlog.jsonl")

	s, err := store.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	entries, err := s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on empty store: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestAppend_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "cronlog.jsonl")

	s, err := store.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	e := newTestEntry("test-job", true)
	if err := s.Append(e); err != nil {
		t.Fatalf("Append: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist at %s: %v", path, err)
	}
	_ = time.Now() // suppress unused import
}
