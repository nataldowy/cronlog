package runner_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/cronlog/internal/runner"
	"github.com/example/cronlog/internal/store"
)

func newTestRunner(t *testing.T, timeout time.Duration) *runner.Runner {
	t.Helper()
	dir := t.TempDir()
	s, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return runner.New(s, timeout)
}

func TestRun_Success(t *testing.T) {
	r := newTestRunner(t, 0)
	entry, err := r.Run("echo", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !entry.Success {
		t.Errorf("expected success, got failure")
	}
	if entry.Output == "" {
		t.Errorf("expected output, got empty string")
	}
}

func TestRun_Failure(t *testing.T) {
	r := newTestRunner(t, 0)
	entry, err := r.Run("false")
	if err != nil {
		t.Fatalf("unexpected store error: %v", err)
	}
	if entry.Success {
		t.Errorf("expected failure, got success")
	}
}

func TestRun_Timeout(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("skipping timeout test outside CI")
	}
	r := newTestRunner(t, 50*time.Millisecond)
	entry, err := r.Run("sleep", "10")
	if err != nil {
		t.Fatalf("unexpected store error: %v", err)
	}
	if entry.Success {
		t.Errorf("expected timed-out entry to be marked failed")
	}
	if !entry.TimedOut {
		t.Errorf("expected TimedOut to be true")
	}
}

func TestRun_PersistsEntry(t *testing.T) {
	dir := t.TempDir()
	s, _ := store.New(dir)
	r := runner.New(s, 0)

	_, err := r.Run("echo", "persisted")
	if err != nil {
		t.Fatalf("run error: %v", err)
	}

	entries, err := s.ReadAll("echo")
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 stored entry, got %d", len(entries))
	}
}
