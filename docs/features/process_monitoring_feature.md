# Process Monitoring Feature

Devopin provides an htop-style process view at `/processes`.

## Live process view

The authenticated `GET /api/processes/live` endpoint returns the current process snapshot. It supports:

- `search`: matches process name, username, or PID.
- `sort`: `cpu`, `memory`, `pid`, `name`, or `threads`.
- `limit`: maximum 500 records; the default is 200.

Process identity uses the PID and process start time together, so a reused PID is treated as a new process.

## Historical storage

Historical samples are stored in the `process_metrics` table in `metrics.db`. The background worker persists a maximum of 100 processes per sample, selecting the top CPU and memory consumers. The live endpoint is not a database write path.

The authenticated `GET /api/processes/history` endpoint accepts:

- `process_key`: the process identity returned by the live endpoint.
- `filter`: `1h`, `6h`, `12h`, `1d`, `7d`, or `30d`.

History is aggregated into time buckets before it is returned to the frontend.

## Retention

Process samples older than 30 days are deleted by the background worker. Cleanup runs at startup and then hourly. The cutoff is evaluated in UTC and records older than the cutoff are removed.

The process worker samples historical data every 30 seconds by default. Individual processes can disappear or deny access while being collected; those failures are skipped without failing the complete snapshot.
