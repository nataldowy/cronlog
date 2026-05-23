// Package logentry defines the core data structure for cronlog.
//
// An Entry captures all relevant metadata about a single cron job execution:
// the job name, command run, start and end times, exit code, captured output,
// and an optional set of key-value labels for filtering and grouping.
//
// Typical usage:
//
//	e := logentry.New("daily-backup", "/usr/local/bin/backup.sh")
//	// ... run the job ...
//	e.Finish(exitCode, stdout, stderr)
//	if e.Failed() {
//		// trigger notification
//	}
//
// Status values:
//
//	StatusSuccess - job exited with code 0
//	StatusFailure - job exited with a non-zero code
//	StatusTimeout - job was killed due to exceeding a time limit
package logentry
