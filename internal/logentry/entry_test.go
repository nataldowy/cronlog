package logentry

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	before := time.Now().UTC()
	e := New("backup", "/usr/bin/backup.sh")
	after := time.Now().UTC()

	if e.JobName != "backup" {
		t.Errorf("expected JobName 'backup', got %q", e.JobName)
	}
	if e.Command != "/usr/bin/backup.sh" {
		t.Errorf("expected Command '/usr/bin/backup.sh', got %q", e.Command)
	}
	if e.StartedAt.Before(before) || e.StartedAt.After(after) {
		t.Error("StartedAt is not within expected range")
	}
	if e.Labels == nil {
		t.Error("Labels map should be initialized")
	}
}

func TestFinish_Success(t *testing.T) {
	e := New("cleanup", "rm -rf /tmp/cache")
	e.Finish(0, "done", "")

	if e.Status != StatusSuccess {
		t.Errorf("expected status %q, got %q", StatusSuccess, e.Status)
	}
	if e.Failed() {
		t.Error("Failed() should return false for exit code 0")
	}
	if e.Duration <= 0 {
		t.Error("Duration should be positive")
	}
	if e.Output != "done" {
		t.Errorf("expected output 'done', got %q", e.Output)
	}
}

func TestFinish_Failure(t *testing.T) {
	e := New("sync", "rsync -av /src /dst")
	e.Finish(1, "", "rsync: connection refused")

	if e.Status != StatusFailure {
		t.Errorf("expected status %q, got %q", StatusFailure, e.Status)
	}
	if !e.Failed() {
		t.Error("Failed() should return true for non-zero exit code")
	}
	if e.Stderr != "rsync: connection refused" {
		t.Errorf("unexpected stderr: %q", e.Stderr)
	}
}

func TestFinish_SetsFinishedAt(t *testing.T) {
	e := New("report", "generate-report")
	before := time.Now().UTC()
	e.Finish(0, "ok", "")
	after := time.Now().UTC()

	if e.FinishedAt.Before(before) || e.FinishedAt.After(after) {
		t.Error("FinishedAt is not within expected range")
	}
}

func TestAddLabel(t *testing.T) {
	e := New("report", "generate-report")
	e.AddLabel("env", "production")
	e.AddLabel("host", "server-01")

	if e.Labels["env"] != "production" {
		t.Errorf("expected label env=production, got %q", e.Labels["env"])
	}
	if e.Labels["host"] != "server-01" {
		t.Errorf("expected label host=server-01, got %q", e.Labels["host"])
	}
}

func TestTimeout_Failed(t *testing.T) {
	e := New("slow-job", "sleep 9999")
	e.Status = StatusTimeout

	if !e.Failed() {
		t.Error("Failed() should return true for timeout status")
	}
}
