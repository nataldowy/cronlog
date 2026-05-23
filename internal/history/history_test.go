package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/cronlog/internal/history"
	"github.com/example/cronlog/internal/logentry"
	"github.com/example/cronlog/internal/store"
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "cronlog-history-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func appendEntry(t *testing.T, st *store.Store, job string, success bool, dur time.Duration) {
	t.Helper()
	e := logentry.New(job, []string{"echo", "hi"})
	time.Sleep(1 * time.Millisecond) // ensure FinishedAt > StartedAt
	e.Finish(success, "", nil)
	_ = dur // duration is recorded internally; we rely on real elapsed time
	if err := st.Append(e); err != nil {
		t.Fatalf("Append: %v", err)
	}
}

func TestSummarise_Empty(t *testing.T) {
	st := newTestStore(t)
	r := history.New(st)
	s, err := r.Summarise("myjob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
}

func TestSummarise_Counts(t *testing.T) {
	st := newTestStore(t)
	appendEntry(t, st, "myjob", true, 0)
	appendEntry(t, st, "myjob", true, 0)
	appendEntry(t, st, "myjob", false, 0)

	r := history.New(st)
	s, err := r.Summarise("myjob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Total != 3 {
		t.Errorf("Total: want 3, got %d", s.Total)
	}
	if s.Successes != 2 {
		t.Errorf("Successes: want 2, got %d", s.Successes)
	}
	if s.Failures != 1 {
		t.Errorf("Failures: want 1, got %d", s.Failures)
	}
	if s.LastStatus != "failure" {
		t.Errorf("LastStatus: want failure, got %s", s.LastStatus)
	}
}

func TestRecent_LimitRespected(t *testing.T) {
	st := newTestStore(t)
	for i := 0; i < 5; i++ {
		appendEntry(t, st, "batch", true, 0)
	}
	r := history.New(st)
	entries, err := r.Recent("batch", 3)
	if err != nil {
		t.Fatalf("Recent: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("want 3 entries, got %d", len(entries))
	}
}

func TestRecent_FewerThanN(t *testing.T) {
	st := newTestStore(t)
	appendEntry(t, st, "sparse", false, 0)
	r := history.New(st)
	entries, err := r.Recent("sparse", 10)
	if err != nil {
		t.Fatalf("Recent: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("want 1 entry, got %d", len(entries))
	}
}
