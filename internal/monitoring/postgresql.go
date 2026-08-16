package monitoring

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	postgresDefaultPort    = 5432
	postgresQueryTimeout   = 8 * time.Second
	postgresMaximumQuery   = 8192
	postgresMinimumVersion = 90600
	postgresBackendVersion = 100000
)

// PostgreSQLCollection contains one native pg_stat_activity snapshot and the
// server version discovered from the same connection.
type PostgreSQLCollection struct {
	Activities    []model.PostgreSQLActivity
	ServerVersion string
	VersionNum    int
}

// PostgreSQLCollector reads PostgreSQL's built-in monitoring views. It does
// not require pg_stat_statements or any other extension.
type PostgreSQLCollector struct {
	timeout      time.Duration
	maximumQuery int
}

func NewPostgreSQLCollector() *PostgreSQLCollector {
	return &PostgreSQLCollector{
		timeout:      postgresQueryTimeout,
		maximumQuery: postgresMaximumQuery,
	}
}

func (c *PostgreSQLCollector) TestConnection(ctx context.Context, target model.PostgreSQLTarget, password string) (string, error) {
	db, err := c.open(ctx, target, password)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var versionNum int
	var version string
	if err := db.QueryRowContext(ctx, "SELECT current_setting('server_version_num')::int, current_setting('server_version')").Scan(&versionNum, &version); err != nil {
		return "", fmt.Errorf("read PostgreSQL server version: %w", err)
	}
	if versionNum < postgresMinimumVersion {
		return version, fmt.Errorf("PostgreSQL %s is older than the supported native monitoring baseline (9.6)", version)
	}
	return version, nil
}

func (c *PostgreSQLCollector) Collect(ctx context.Context, target model.PostgreSQLTarget, password string) (PostgreSQLCollection, error) {
	db, err := c.open(ctx, target, password)
	if err != nil {
		return PostgreSQLCollection{}, err
	}
	defer db.Close()

	var versionNum int
	var serverVersion string
	if err := db.QueryRowContext(ctx, "SELECT current_setting('server_version_num')::int, current_setting('server_version')").Scan(&versionNum, &serverVersion); err != nil {
		return PostgreSQLCollection{}, fmt.Errorf("read PostgreSQL server version: %w", err)
	}
	if versionNum < postgresMinimumVersion {
		return PostgreSQLCollection{}, fmt.Errorf("PostgreSQL %s is older than the supported native monitoring baseline (9.6)", serverVersion)
	}

	query := activityQuery(versionNum)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return PostgreSQLCollection{}, fmt.Errorf("read pg_stat_activity: %w", err)
	}
	defer rows.Close()

	observedAt := time.Now().UTC()
	activities := make([]model.PostgreSQLActivity, 0)
	for rows.Next() {
		activity, err := scanActivity(rows, target, serverVersion, observedAt, versionNum >= postgresBackendVersion, c.maximumQuery)
		if err != nil {
			return PostgreSQLCollection{}, fmt.Errorf("scan pg_stat_activity row: %w", err)
		}
		activities = append(activities, activity)
	}
	if err := rows.Err(); err != nil {
		return PostgreSQLCollection{}, fmt.Errorf("iterate pg_stat_activity: %w", err)
	}

	return PostgreSQLCollection{
		Activities:    activities,
		ServerVersion: serverVersion,
		VersionNum:    versionNum,
	}, nil
}

func (c *PostgreSQLCollector) open(ctx context.Context, target model.PostgreSQLTarget, password string) (*sql.DB, error) {
	dsn, err := postgresDSN(target, password)
	if err != nil {
		return nil, err
	}
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse PostgreSQL connection: %w", err)
	}
	config.RuntimeParams["application_name"] = "devopin"

	db := stdlib.OpenDB(*config)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}
	return db, nil
}

func postgresDSN(target model.PostgreSQLTarget, password string) (string, error) {
	host := strings.TrimSpace(target.Host)
	if host == "" {
		return "", fmt.Errorf("PostgreSQL host is required")
	}
	port := target.Port
	if port == 0 {
		port = postgresDefaultPort
	}
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("PostgreSQL port must be between 1 and 65535")
	}
	databaseName := strings.TrimSpace(target.Database)
	if databaseName == "" {
		databaseName = "postgres"
	}
	username := strings.TrimSpace(target.Username)
	if username == "" {
		return "", fmt.Errorf("PostgreSQL username is required")
	}
	sslMode := strings.ToLower(strings.TrimSpace(target.SSLMode))
	if sslMode == "" {
		sslMode = "prefer"
	}
	switch sslMode {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
	default:
		return "", fmt.Errorf("unsupported PostgreSQL SSL mode %q", sslMode)
	}

	connectionURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
		Path:   "/" + databaseName,
	}
	query := connectionURL.Query()
	query.Set("sslmode", sslMode)
	connectionURL.RawQuery = query.Encode()
	return connectionURL.String(), nil
}

func activityQuery(versionNum int) string {
	columns := `
		pid,
		COALESCE(datname, ''),
		COALESCE(usename, ''),
		COALESCE(application_name, ''),
		COALESCE(client_addr::text, ''),
		COALESCE(client_port, 0),
		backend_start,
		xact_start,
		query_start,
		state_change,
		COALESCE(wait_event_type, ''),
		COALESCE(wait_event, ''),
		COALESCE(state, ''),
		COALESCE(query, ''),
		COALESCE((pg_blocking_pids(pid))[1], 0) AS blocking_pid`
	if versionNum >= postgresBackendVersion {
		columns += `,
		COALESCE(backend_type, '')`
	}
	return `WITH activity AS (SELECT` + columns + `
	FROM pg_stat_activity
	WHERE pid <> pg_backend_pid()
	  AND state IS DISTINCT FROM 'idle')
	SELECT activity.*,
	       COALESCE(blocker.datname, ''),
	       COALESCE(blocker.usename, ''),
	       COALESCE(blocker.state, ''),
	       COALESCE(blocker.query, '')
	FROM activity
	LEFT JOIN pg_stat_activity blocker ON blocker.pid = activity.blocking_pid
	ORDER BY activity.query_start NULLS LAST, activity.pid`
}

func scanActivity(rows *sql.Rows, target model.PostgreSQLTarget, serverVersion string, observedAt time.Time, hasBackendType bool, maximumQuery int) (model.PostgreSQLActivity, error) {
	var (
		pid, clientPort, blockingPID                      int64
		datname, username, applicationName, clientAddress string
		backendStart, transactionStart, queryStart        sql.NullTime
		stateChange                                       sql.NullTime
		waitEventType, waitEvent, state, query            string
		blockingDatabase, blockingUsername, blockingState string
		blockingQuery                                     string
		backendType                                       sql.NullString
	)
	args := []any{
		&pid,
		&datname,
		&username,
		&applicationName,
		&clientAddress,
		&clientPort,
		&backendStart,
		&transactionStart,
		&queryStart,
		&stateChange,
		&waitEventType,
		&waitEvent,
		&state,
		&query,
		&blockingPID,
	}
	if hasBackendType {
		args = append(args, &backendType)
	}
	args = append(args, &blockingDatabase, &blockingUsername, &blockingState, &blockingQuery)
	if err := rows.Scan(args...); err != nil {
		return model.PostgreSQLActivity{}, err
	}

	query = truncateQuery(query, maximumQuery)
	blockingQuery = truncateQuery(blockingQuery, maximumQuery)
	activityKey := fmt.Sprintf("%d:%d:%d", target.ID, pid, timestampValue(backendStart))
	queryStartTime := nullableTimeValue(queryStart)
	durationMS := int64(0)
	if queryStartTime != nil {
		if elapsed := observedAt.Sub(*queryStartTime); elapsed > 0 {
			durationMS = elapsed.Milliseconds()
		}
	}

	return model.PostgreSQLActivity{
		ActivityKey:      activityKey,
		TargetID:         target.ID,
		TargetName:       target.Name,
		ServerVersion:    serverVersion,
		PID:              pid,
		DatabaseName:     datname,
		Username:         username,
		ApplicationName:  applicationName,
		ClientAddress:    clientAddress,
		ClientPort:       int(clientPort),
		BackendType:      backendType.String,
		BackendStart:     nullableTimeValue(backendStart),
		TransactionStart: nullableTimeValue(transactionStart),
		QueryStart:       queryStartTime,
		StateChange:      nullableTimeValue(stateChange),
		WaitEventType:    waitEventType,
		WaitEvent:        waitEvent,
		BlockingPID:      blockingPID,
		BlockingDatabase: blockingDatabase,
		BlockingUsername: blockingUsername,
		BlockingState:    blockingState,
		BlockingQuery:    blockingQuery,
		State:            state,
		Query:            query,
		QueryFingerprint: queryFingerprint(query),
		QueryDurationMS:  durationMS,
		ObservedAt:       observedAt,
	}, nil
}

func nullableTimeValue(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	timestamp := value.Time.UTC()
	return &timestamp
}

func timestampValue(value sql.NullTime) int64 {
	if !value.Valid {
		return 0
	}
	return value.Time.UnixNano()
}

func truncateQuery(query string, maximum int) string {
	query = strings.TrimSpace(query)
	if maximum <= 0 || len([]rune(query)) <= maximum {
		return query
	}
	return string([]rune(query)[:maximum])
}

func queryFingerprint(query string) string {
	normalized := strings.Join(strings.Fields(strings.ToLower(query)), " ")
	if normalized == "" {
		return ""
	}
	digest := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(digest[:])
}
