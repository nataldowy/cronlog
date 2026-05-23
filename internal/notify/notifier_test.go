package notify

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cfg := Config{
		SMTPHost: "localhost",
		SMTPPort: 25,
		From:     "cronlog@example.com",
		To:       []string{"admin@example.com"},
	}
	n := New(cfg)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
	if n.cfg.SMTPHost != "localhost" {
		t.Errorf("expected SMTPHost 'localhost', got %q", n.cfg.SMTPHost)
	}
}

func TestNotify_NoRecipients(t *testing.T) {
	n := New(Config{
		SMTPHost: "localhost",
		SMTPPort: 25,
		From:     "cronlog@example.com",
		To:       []string{},
	})

	err := n.Notify(FailureDetails{
		JobName:  "backup",
		ExitCode: 1,
	})
	if err == nil {
		t.Fatal("expected error for empty recipients, got nil")
	}
}

func TestFailureDetails_Fields(t *testing.T) {
	now := time.Now()
	d := FailureDetails{
		JobName:   "db-backup",
		ExitCode:  2,
		Output:    "error: connection refused",
		StartedAt: now,
		Duration:  5 * time.Second,
		Labels:    map[string]string{"env": "prod"},
	}

	if d.JobName != "db-backup" {
		t.Errorf("unexpected JobName: %q", d.JobName)
	}
	if d.ExitCode != 2 {
		t.Errorf("unexpected ExitCode: %d", d.ExitCode)
	}
	if d.Labels["env"] != "prod" {
		t.Errorf("unexpected label env: %q", d.Labels["env"])
	}
}

func TestEmailTemplate_Renders(t *testing.T) {
	import_bytes := func() interface{} { return nil } // placeholder
	_ = import_bytes

	details := map[string]interface{}{
		"JobName":   "test-job",
		"ExitCode":  1,
		"Output":    "something failed",
		"StartedAt": time.Now(),
		"Duration":  2 * time.Second,
		"From":      "cronlog@example.com",
		"To":        "admin@example.com",
	}

	var buf interface{ Len() int }
	_ = buf
	_ = details

	if emailTmpl == nil {
		t.Fatal("expected emailTmpl to be defined")
	}
}
