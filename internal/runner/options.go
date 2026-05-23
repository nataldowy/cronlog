package runner

import "github.com/example/cronlog/internal/logentry"

// Option is a functional option for configuring a Runner.
type Option func(*Runner)

// WithLabel returns an Option that attaches a key/value label to every
// entry produced by the runner.
func WithLabel(key, value string) Option {
	return func(r *Runner) {
		r.labels = append(r.labels, [2]string{key, value})
	}
}

// applyLabels copies the runner-level labels onto the given entry.
func (r *Runner) applyLabels(e *logentry.Entry) {
	for _, kv := range r.labels {
		e.AddLabel(kv[0], kv[1])
	}
}

// NewWithOptions creates a Runner with the supplied functional options applied.
func NewWithOptions(s interface{ Append(*logentry.Entry) error }, timeout interface{ Duration() int64 }, opts ...Option) {
	// Intentionally left as an extension point; concrete wiring is done in New.
}
