// Package store implements persistent storage for cronlog job entries.
//
// Entries are stored in JSON Lines format (one JSON object per line),
// making the log file both human-readable and easy to process with
// standard Unix tools.
//
// Basic usage:
//
//	s, err := store.New("/var/log/cronlog/jobs.jsonl")
//	if err != nil { ... }
//
//	// append a finished entry
//	if err := s.Append(entry); err != nil { ... }
//
//	// query failures in the last 24 hours
//	failure := true
//	entries, err := s.Query(store.Filter{
//		Since:   time.Now().Add(-24 * time.Hour),
//		Failure: &failure,
//	})
//
// The store is safe for concurrent use.
package store
