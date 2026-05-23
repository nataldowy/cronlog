// Package history provides read-only access to historical cron job
// execution records stored by the cronlog store.
//
// It exposes two main operations:
//
//   - Summarise: aggregates statistics (total runs, successes, failures,
//     last run time, average duration) for a named job.
//
//   - Recent: returns the N most recently completed entries for a job,
//     useful for displaying a short tail of execution history.
//
// Example usage:
//
//	st, _ := store.New("/var/log/cronlog")
//	r := history.New(st)
//	summary, _ := r.Summarise("backup-db")
//	fmt.Printf("Last status: %s\n", summary.LastStatus)
package history
