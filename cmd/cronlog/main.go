// Command cronlog wraps a cron job command, captures its output,
// stores a structured log entry, and sends failure notifications.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/example/cronlog/internal/config"
	"github.com/example/cronlog/internal/notify"
	"github.com/example/cronlog/internal/runner"
	"github.com/example/cronlog/internal/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "cronlog: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		return fmt.Errorf("usage: cronlog [options] <command> [args...]")
	}

	// Parse flags: --job=<name> --label=<k=v> before command.
	jobName, command, labels := parseArgs(args)

	s, err := store.New(cfg.StorePath)
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}

	notifier := notify.New(cfg)

	opts := []runner.Option{
		runner.WithLabel("job", jobName),
	}
	for k, v := range labels {
		opts = append(opts, runner.WithLabel(k, v))
	}
	if cfg.TimeoutSeconds > 0 {
		opts = append(opts, runner.WithTimeout(time.Duration(cfg.TimeoutSeconds)*time.Second))
	}

	r := runner.NewWithOptions(command[0], command[1:], s, notifier, opts...)
	return r.Run()
}
