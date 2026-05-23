// Package notify provides failure notification support for cron job log entries.
package notify

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
	"time"
)

// Config holds the configuration for the email notifier.
type Config struct {
	SMTPHost   string
	SMTPPort   int
	From       string
	To         []string
	Username   string
	Password   string
}

// Notifier sends failure notifications.
type Notifier struct {
	cfg Config
}

// New creates a new Notifier with the given configuration.
func New(cfg Config) *Notifier {
	return &Notifier{cfg: cfg}
}

// FailureDetails contains information about a failed cron job.
type FailureDetails struct {
	JobName   string
	ExitCode  int
	Output    string
	StartedAt time.Time
	Duration  time.Duration
	Labels    map[string]string
}

var emailTmpl = template.Must(template.New("failure").Parse(`Subject: [cronlog] Job "{{.JobName}}" failed
From: {{.From}}
To: {{.To}}

Cron job failure detected.

Job:      {{.JobName}}
Exit Code: {{.ExitCode}}
Started:  {{.StartedAt.Format "2006-01-02 15:04:05"}}
Duration: {{.Duration}}

Output:
{{.Output}}
`))

// Notify sends an email notification for a failed job.
func (n *Notifier) Notify(details FailureDetails) error {
	if len(n.cfg.To) == 0 {
		return fmt.Errorf("notify: no recipients configured")
	}

	var buf bytes.Buffer
	err := emailTmpl.Execute(&buf, map[string]interface{}{
		"JobName":   details.JobName,
		"ExitCode":  details.ExitCode,
		"Output":    details.Output,
		"StartedAt": details.StartedAt,
		"Duration":  details.Duration,
		"From":      n.cfg.From,
		"To":        n.cfg.To[0],
	})
	if err != nil {
		return fmt.Errorf("notify: render template: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", n.cfg.SMTPHost, n.cfg.SMTPPort)
	var auth smtp.Auth
	if n.cfg.Username != "" {
		auth = smtp.PlainAuth("", n.cfg.Username, n.cfg.Password, n.cfg.SMTPHost)
	}

	if err := smtp.SendMail(addr, auth, n.cfg.From, n.cfg.To, buf.Bytes()); err != nil {
		return fmt.Errorf("notify: send mail: %w", err)
	}
	return nil
}
