package monitoring

import (
	"testing"

	"github.com/gabutlabs/godevopin/internal/model"
)

func TestMySQLConnectionConfigKeepsCredentialsOutOfDSNFormatting(t *testing.T) {
	config, err := mysqlConnectionConfig(model.MySQLTarget{
		Host:     "db.example.test",
		Username: "monitor@example",
		Database: "application",
	}, "p@ss word")
	if err != nil {
		t.Fatalf("mysqlConnectionConfig() error = %v", err)
	}
	if config.Passwd != "p@ss word" {
		t.Fatalf("mysqlConnectionConfig() changed the password: %q", config.Passwd)
	}
	if config.Addr != "db.example.test:3306" || config.DBName != "application" || config.TLSConfig != "preferred" {
		t.Fatalf("unexpected MySQL connection defaults: %+v", config)
	}
}

func TestMySQLActivityHelpers(t *testing.T) {
	if mysqlClientAddress("10.0.0.5:54321") != "10.0.0.5" {
		t.Fatal("mysqlClientAddress() should strip the client port")
	}
	if mysqlClientAddress("localhost") != "localhost" {
		t.Fatal("mysqlClientAddress() should preserve host-only values")
	}
	if mysqlQueryFingerprint("SELECT  * FROM users") != mysqlQueryFingerprint(" select *   from users ") {
		t.Fatal("MySQL query fingerprints should normalize whitespace and case")
	}
}
