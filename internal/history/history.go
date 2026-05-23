// Package history provides utilities for querying and summarising
// past cron job log entries from the store.
package history

import (
	"fmt"
	"time"

	"github.com/example/cronlog/internal/logentry"
	"github.com/example/cronlog/internal/store"
)

// Summary holds aggregated statistics for a single job over a time window.
type Summary struct {
	JobName      string
	Total        int
	Successes    int
	Failures     int
	LastRun      time.Time
	LastStatus   string
	AvgDuration  time.Duration
}

// Reader wraps a store and exposes history queries.
type Reader struct {
	st *store.Store
}

// New returns a Reader backed by the given store.
func New(st *store.Store) *Reader {
	return &Reader{st: st}
}

// Summarise returns a Summary for the named job using all stored entries.
func (r *Reader) Summarise(jobName string) (*Summary, error) {
	entries, err := r.st.ReadAll(jobName)
	if err != nil {
		return nil, fmt.Errorf("history: read entries for %q: %w", jobName, err)
	}

	s := &Summary{JobName: jobName}
	var totalDuration time.Duration

	for _, e := range entries {
		s.Total++
		if e.Success {
			s.Successes++
		} else {
			s.Failures++
		}
		if e.FinishedAt.After(s.LastRun) {
			s.LastRun = e.FinishedAt
			if e.Success {
				s.LastStatus = "success"
			} else {
				s.LastStatus = "failure"
			}
		}
		totalDuration += e.Duration()
	}

	if s.Total > 0 {
		s.AvgDuration = totalDuration / time.Duration(s.Total)
	}

	return s, nil
}

// Recent returns up to n most recent entries for the named job.
func (r *Reader) Recent(jobName string, n int) ([]*logentry.Entry, error) {
	all, err := r.st.ReadAll(jobName)
	if err != nil {
		return nil, fmt.Errorf("history: recent entries for %q: %w", jobName, err)
	}
	if len(all) <= n {
		return all, nil
	}
	return all[len(all)-n:], nil
}
