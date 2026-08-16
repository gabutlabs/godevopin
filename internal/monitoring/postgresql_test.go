package monitoring

import (
	"strings"
	"testing"

	"github.com/gabutlabs/godevopin/internal/model"
)

func TestPostgreSQLDSNUsesEscapedCredentialsAndDefaults(t *testing.T) {
	dsn, err := postgresDSN(model.PostgreSQLTarget{
		Host:     "db.example.test",
		Username: "monitor@example",
		Database: "application",
	}, "p@ss word")
	if err != nil {
		t.Fatalf("postgresDSN() error = %v", err)
	}
	if !strings.Contains(dsn, "sslmode=prefer") || !strings.Contains(dsn, "5432") {
		t.Fatalf("postgresDSN() = %q, missing default connection settings", dsn)
	}
	if strings.Contains(dsn, "p@ss word") {
		t.Fatalf("postgresDSN() exposed the raw password: %q", dsn)
	}
}

func TestActivityQueryIsVersionAware(t *testing.T) {
	legacyQuery := activityQuery(90600)
	if strings.Contains(legacyQuery, "backend_type") {
		t.Fatal("legacy PostgreSQL query should not select backend_type")
	}
	modernQuery := activityQuery(100000)
	if !strings.Contains(modernQuery, "backend_type") {
		t.Fatal("modern PostgreSQL query should select backend_type")
	}
	if !strings.Contains(modernQuery, "state IS DISTINCT FROM 'idle'") {
		t.Fatal("activity query should exclude plain idle sessions")
	}
	if !strings.Contains(modernQuery, "pg_blocking_pids(pid)") {
		t.Fatal("activity query should include the first blocking PID")
	}
	if !strings.Contains(modernQuery, "blocker.query") {
		t.Fatal("activity query should include the blocking query")
	}
}

func TestTruncateQueryAndFingerprint(t *testing.T) {
	query := truncateQuery("abcdefghij", 5)
	if query != "abcde" {
		t.Fatalf("truncateQuery() = %q", query)
	}
	if queryFingerprint("SELECT  * FROM users") != queryFingerprint(" select *   from users ") {
		t.Fatal("query fingerprints should normalize whitespace and case")
	}
}
