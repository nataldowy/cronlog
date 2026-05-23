package main

import (
	"strings"
)

// parseArgs splits cronlog's own flags from the wrapped command.
// Supported flags (must appear before the command):
//
//	--job=<name>    logical job name stored with the log entry
//	--label=<k=v>   arbitrary key=value label (repeatable)
func parseArgs(args []string) (jobName string, command []string, labels map[string]string) {
	labels = make(map[string]string)
	jobName = "unknown"

	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "--job=") {
			jobName = strings.TrimPrefix(arg, "--job=")
			i++
			continue
		}
		if strings.HasPrefix(arg, "--label=") {
			kv := strings.TrimPrefix(arg, "--label=")
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 {
				labels[parts[0]] = parts[1]
			}
			i++
			continue
		}
		// First non-flag argument is the start of the command.
		break
	}

	command = args[i:]
	return jobName, command, labels
}
