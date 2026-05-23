// Package runner executes shell commands as cron jobs and records
// structured log entries for each run.
package runner

import (
	"context"
	"os/exec"
	"time"

	"github.com/example/cronlog/internal/logentry"
	"github.com/example/cronlog/internal/store"
)

// Runner executes commands and persists their log entries.
type Runner struct {
	store   *store.Store
	timeout time.Duration
}

// New creates a Runner backed by the given store.
// If timeout is zero, commands are allowed to run indefinitely.
func New(s *store.Store, timeout time.Duration) *Runner {
	return &Runner{store: s, timeout: timeout}
}

// Run executes the named command with the provided arguments.
// It records stdout/stderr output and exit status in a log entry.
func (r *Runner) Run(name string, args ...string) (*logentry.Entry, error) {
	entry := logentry.New(name)

	ctx := context.Background()
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	entry.Output = string(output)

	if ctx.Err() == context.DeadlineExceeded {
		entry.Timeout()
	} else {
		entry.Finish(err)
	}

	if storeErr := r.store.Append(entry); storeErr != nil {
		return entry, storeErr
	}
	return entry, nil
}
