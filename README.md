# cronlog

Structured log aggregator for cron jobs with failure notifications and history.

## Installation

```bash
go install github.com/yourusername/cronlog@latest
```

## Usage

Wrap any cron job with `cronlog` to capture output, track history, and receive failure notifications:

```bash
# Basic usage
cronlog -- /path/to/your/script.sh

# With a job name and email notification on failure
cronlog --name "daily-backup" --notify user@example.com -- ./backup.sh

# View job history
cronlog history --name "daily-backup"

# List recent runs with status
cronlog list --limit 20
```

### Crontab example

```cron
0 2 * * * cronlog --name "nightly-sync" --notify ops@example.com -- /opt/scripts/sync.sh
```

### Configuration

`cronlog` reads from `~/.cronlog.yaml` or `/etc/cronlog/config.yaml`:

```yaml
notify:
  smtp_host: smtp.example.com
  smtp_port: 587
  from: cronlog@example.com

storage:
  path: /var/log/cronlog
  retention_days: 30
```

## Features

- Structured JSON log capture per job run
- Failure detection with email notifications
- Persistent run history with exit codes and durations
- Simple CLI for querying past executions

## License

MIT