package history

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// PrintSummary writes a human-readable summary table to w.
func PrintSummary(w io.Writer, s *Summary) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	lines := []struct{ k, v string }{
		{"Job", s.JobName},
		{"Total runs", fmt.Sprintf("%d", s.Total)},
		{"Successes", fmt.Sprintf("%d", s.Successes)},
		{"Failures", fmt.Sprintf("%d", s.Failures)},
		{"Last run", formatTime(s)},
		{"Last status", lastStatus(s)},
		{"Avg duration", s.AvgDuration.Round(1e6).String()},
	}
	for _, l := range lines {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", l.k, l.v); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func formatTime(s *Summary) string {
	if s.Total == 0 {
		return "never"
	}
	return s.LastRun.Format("2006-01-02 15:04:05 MST")
}

func lastStatus(s *Summary) string {
	if s.Total == 0 {
		return "n/a"
	}
	return strings.ToUpper(s.LastStatus[:1]) + s.LastStatus[1:]
}
