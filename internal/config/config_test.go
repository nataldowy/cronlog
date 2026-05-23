package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/cronlog/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "cronlog.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTemp(t, "store:\n  path: \"\"\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Store.Path != "/var/log/cronlog" {
		t.Errorf("expected default store path, got %q", cfg.Store.Path)
	}
	if cfg.Notify.SMTPPort != 25 {
		t.Errorf("expected default smtp port 25, got %d", cfg.Notify.SMTPPort)
	}
	if cfg.Runner.DefaultTimeout != 24*time.Hour {
		t.Errorf("expected default timeout 24h, got %v", cfg.Runner.DefaultTimeout)
	}
}

func TestLoad_FullConfig(t *testing.T) {
	yaml := `
store:
  path: /tmp/logs
notify:
  smtp_host: mail.example.com
  smtp_port: 587
  from: cronlog@example.com
  recipients:
    - ops@example.com
runner:
  default_timeout: 1h
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Store.Path != "/tmp/logs" {
		t.Errorf("store path: got %q", cfg.Store.Path)
	}
	if cfg.Notify.SMTPHost != "mail.example.com" {
		t.Errorf("smtp host: got %q", cfg.Notify.SMTPHost)
	}
	if len(cfg.Notify.Recipients) != 1 || cfg.Notify.Recipients[0] != "ops@example.com" {
		t.Errorf("recipients: got %v", cfg.Notify.Recipients)
	}
	if cfg.Runner.DefaultTimeout != time.Hour {
		t.Errorf("default timeout: got %v", cfg.Runner.DefaultTimeout)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/cronlog.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
