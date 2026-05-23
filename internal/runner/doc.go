// Package runner provides the core execution layer for cronlog.
//
// It wraps os/exec to run arbitrary shell commands, captures their combined
// output, and records structured [logentry.Entry] values via the store package.
//
// Basic usage:
//
//	s, _ := store.New("/var/log/cronlog")
//	r := runner.New(s, 30*time.Second)
//	entry, err := r.Run("pg_dump", "-Fc", "mydb")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if !entry.Success {
//		// send notification, etc.
//	}
//
// A zero timeout disables the deadline so commands may run as long as needed.
package runner
