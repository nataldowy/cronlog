package store

import (
	"time"

	"github.com/user/cronlog/internal/logentry"
)

// Filter holds optional criteria for querying log entries.
type Filter struct {
	Job     string    // filter by job name (empty means all)
	Since   time.Time // entries with StartedAt after this time
	Failure *bool     // if non-nil, filter by success/failure
}

// Query returns entries from the store matching the given filter.
func (s *Store) Query(f Filter) ([]*logentry.Entry, error) {
	all, err := s.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []*logentry.Entry
	for _, e := range all {
		if f.Job != "" && e.Job != f.Job {
			continue
		}
		if !f.Since.IsZero() && e.StartedAt.Before(f.Since) {
			continue
		}
		if f.Failure != nil && e.Success != !*f.Failure {
			continue
		}
		result = append(result, e)
	}
	return result, nil
}

// LatestByJob returns the most recent entry for each distinct job name.
func (s *Store) LatestByJob() (map[string]*logentry.Entry, error) {
	all, err := s.ReadAll()
	if err != nil {
		return nil, err
	}

	latest := make(map[string]*logentry.Entry)
	for _, e := range all {
		existing, ok := latest[e.Job]
		if !ok || e.StartedAt.After(existing.StartedAt) {
			latest[e.Job] = e
		}
	}
	return latest, nil
}
