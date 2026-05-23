package main

import (
	"testing"
)

func TestParseArgs_JobAndCommand(t *testing.T) {
	job, cmd, labels := parseArgs([]string{"--job=backup", "tar", "-czf", "/tmp/out.tar.gz", "/data"})
	if job != "backup" {
		t.Errorf("expected job=backup, got %q", job)
	}
	if len(cmd) != 4 || cmd[0] != "tar" {
		t.Errorf("unexpected command: %v", cmd)
	}
	if len(labels) != 0 {
		t.Errorf("expected no labels, got %v", labels)
	}
}

func TestParseArgs_Labels(t *testing.T) {
	_, _, labels := parseArgs([]string{"--label=env=prod", "--label=region=us-east", "echo", "hi"})
	if labels["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", labels["env"])
	}
	if labels["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %q", labels["region"])
	}
}

func TestParseArgs_NoFlags(t *testing.T) {
	job, cmd, labels := parseArgs([]string{"echo", "hello"})
	if job != "unknown" {
		t.Errorf("expected default job name, got %q", job)
	}
	if len(cmd) != 2 {
		t.Errorf("expected 2 command args, got %d", len(cmd))
	}
	if len(labels) != 0 {
		t.Errorf("expected no labels")
	}
}

func TestParseArgs_Empty(t *testing.T) {
	job, cmd, labels := parseArgs([]string{})
	if job != "unknown" {
		t.Errorf("expected unknown job, got %q", job)
	}
	if len(cmd) != 0 {
		t.Errorf("expected empty command")
	}
	if len(labels) != 0 {
		t.Errorf("expected no labels")
	}
}
