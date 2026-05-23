package logentry

import "time"

// Status represents the outcome of a cron job execution.
type Status string

const (
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
	StatusTimeout Status = "timeout"
)

// Entry represents a single structured log record for a cron job run.
type Entry struct {
	ID        string            `json:"id"`
	JobName   string            `json:"job_name"`
	Command   string            `json:"command"`
	Status    Status            `json:"status"`
	ExitCode  int               `json:"exit_code"`
	Output    string            `json:"output"`
	Stderr    string            `json:"stderr,omitempty"`
	StartedAt time.Time         `json:"started_at"`
	EndedAt   time.Time         `json:"ended_at"`
	Duration  float64           `json:"duration_seconds"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// New creates a new Entry with the given job name and command.
func New(jobName, command string) *Entry {
	return &Entry{
		JobName:   jobName,
		Command:   command,
		StartedAt: time.Now().UTC(),
		Labels:    make(map[string]string),
	}
}

// Finish marks the entry as complete, computing duration and setting status.
func (e *Entry) Finish(exitCode int, output, stderr string) {
	e.EndedAt = time.Now().UTC()
	e.Duration = e.EndedAt.Sub(e.StartedAt).Seconds()
	e.ExitCode = exitCode
	e.Output = output
	e.Stderr = stderr

	if exitCode == 0 {
		e.Status = StatusSuccess
	} else {
		e.Status = StatusFailure
	}
}

// Failed returns true if the entry represents a failed job run.
func (e *Entry) Failed() bool {
	return e.Status == StatusFailure || e.Status == StatusTimeout
}

// AddLabel attaches a key-value label to the entry.
func (e *Entry) AddLabel(key, value string) {
	if e.Labels == nil {
		e.Labels = make(map[string]string)
	}
	e.Labels[key] = value
}
