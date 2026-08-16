package monitoring

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/go-sql-driver/mysql"
)

const (
	mysqlDefaultPort  = 3306
	mysqlQueryTimeout = 8 * time.Second
	mysqlMaximumQuery = 8192
)

type MySQLCollection struct {
	Activities    []model.MySQLActivity
	ServerVersion string
}

// MySQLCollector reads the native SHOW FULL PROCESSLIST output. It does not
// require the performance_schema or any third-party MySQL plugin.
type MySQLCollector struct {
	timeout      time.Duration
	maximumQuery int
}

func NewMySQLCollector() *MySQLCollector {
	return &MySQLCollector{timeout: mysqlQueryTimeout, maximumQuery: mysqlMaximumQuery}
}

func (c *MySQLCollector) TestConnection(ctx context.Context, target model.MySQLTarget, password string) (string, error) {
	db, err := c.open(ctx, target, password)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var version string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		return "", fmt.Errorf("read MySQL server version: %w", err)
	}
	return version, nil
}

func (c *MySQLCollector) Collect(ctx context.Context, target model.MySQLTarget, password string) (MySQLCollection, error) {
	db, err := c.open(ctx, target, password)
	if err != nil {
		return MySQLCollection{}, err
	}
	defer db.Close()

	var serverVersion string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&serverVersion); err != nil {
		return MySQLCollection{}, fmt.Errorf("read MySQL server version: %w", err)
	}
	var ownID int64
	if err := db.QueryRowContext(ctx, "SELECT CONNECTION_ID()").Scan(&ownID); err != nil {
		return MySQLCollection{}, fmt.Errorf("read MySQL connection ID: %w", err)
	}

	rows, err := db.QueryContext(ctx, "SHOW FULL PROCESSLIST")
	if err != nil {
		return MySQLCollection{}, fmt.Errorf("read SHOW FULL PROCESSLIST: %w", err)
	}
	defer rows.Close()

	observedAt := time.Now().UTC()
	activities := make([]model.MySQLActivity, 0)
	for rows.Next() {
		var (
			id                    int64
			username, host        sql.NullString
			databaseName, command sql.NullString
			seconds               sql.NullInt64
			state, info           sql.NullString
		)
		if err := rows.Scan(&id, &username, &host, &databaseName, &command, &seconds, &state, &info); err != nil {
			return MySQLCollection{}, fmt.Errorf("scan SHOW FULL PROCESSLIST row: %w", err)
		}
		if id == ownID || strings.EqualFold(strings.TrimSpace(command.String), "sleep") {
			continue
		}

		query := truncateQuery(info.String, c.maximumQuery)
		queryDurationMS := seconds.Int64 * 1000
		queryStart := observedAt.Add(-time.Duration(seconds.Int64) * time.Second)
		stateValue := strings.TrimSpace(state.String)
		if stateValue == "" {
			stateValue = strings.TrimSpace(command.String)
		}
		activities = append(activities, model.MySQLActivity{
			ActivityKey:      fmt.Sprintf("%d:%d", target.ID, id),
			TargetID:         target.ID,
			TargetName:       target.Name,
			ServerVersion:    serverVersion,
			PID:              id,
			DatabaseName:     databaseName.String,
			Username:         username.String,
			ClientAddress:    mysqlClientAddress(host.String),
			Command:          command.String,
			State:            stateValue,
			Query:            query,
			QueryFingerprint: mysqlQueryFingerprint(query),
			QueryDurationMS:  queryDurationMS,
			QueryStart:       &queryStart,
			ObservedAt:       observedAt,
		})
	}
	if err := rows.Err(); err != nil {
		return MySQLCollection{}, fmt.Errorf("iterate SHOW FULL PROCESSLIST: %w", err)
	}

	return MySQLCollection{Activities: activities, ServerVersion: serverVersion}, nil
}

func (c *MySQLCollector) open(ctx context.Context, target model.MySQLTarget, password string) (*sql.DB, error) {
	config, err := mysqlConnectionConfig(target, password)
	if err != nil {
		return nil, err
	}
	connector, err := mysql.NewConnector(config)
	if err != nil {
		return nil, fmt.Errorf("configure MySQL connection: %w", err)
	}
	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to MySQL: %w", err)
	}
	return db, nil
}

func mysqlConnectionConfig(target model.MySQLTarget, password string) (*mysql.Config, error) {
	host := strings.TrimSpace(target.Host)
	if host == "" {
		return nil, fmt.Errorf("MySQL host is required")
	}
	port := target.Port
	if port == 0 {
		port = mysqlDefaultPort
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("MySQL port must be between 1 and 65535")
	}
	username := strings.TrimSpace(target.Username)
	if username == "" {
		return nil, fmt.Errorf("MySQL username is required")
	}
	databaseName := strings.TrimSpace(target.Database)
	if databaseName == "" {
		databaseName = "mysql"
	}
	sslMode := strings.ToLower(strings.TrimSpace(target.SSLMode))
	if sslMode == "" {
		sslMode = "preferred"
	}
	if sslMode != "disable" && sslMode != "preferred" && sslMode != "required" {
		return nil, fmt.Errorf("unsupported MySQL SSL mode %q", sslMode)
	}
	tlsMode := "preferred"
	if sslMode == "disable" {
		tlsMode = "false"
	} else if sslMode == "required" {
		tlsMode = "true"
	}

	return &mysql.Config{
		User:      username,
		Passwd:    password,
		Net:       "tcp",
		Addr:      net.JoinHostPort(host, strconv.Itoa(port)),
		DBName:    databaseName,
		TLSConfig: tlsMode,
		Params: map[string]string{
			"parseTime":    "true",
			"timeout":      mysqlQueryTimeout.String(),
			"readTimeout":  mysqlQueryTimeout.String(),
			"writeTimeout": mysqlQueryTimeout.String(),
		},
	}, nil
}

func mysqlClientAddress(host string) string {
	if address, _, err := net.SplitHostPort(host); err == nil {
		return address
	}
	return host
}

func mysqlQueryFingerprint(query string) string {
	normalized := strings.Join(strings.Fields(strings.ToLower(query)), " ")
	if normalized == "" {
		return ""
	}
	digest := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(digest[:])
}
