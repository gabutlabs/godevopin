package service

import "testing"

func TestPostgreSQLCredentialEncryptionRoundTrip(t *testing.T) {
	service := &postgresqlActivityService{secret: "test-secret"}
	ciphertext, err := service.encrypt("p@ssword with spaces")
	if err != nil {
		t.Fatalf("encrypt() error = %v", err)
	}
	if ciphertext == "p@ssword with spaces" || ciphertext == "" {
		t.Fatalf("encrypt() returned an unsafe or empty value: %q", ciphertext)
	}
	plaintext, err := service.decrypt(ciphertext)
	if err != nil {
		t.Fatalf("decrypt() error = %v", err)
	}
	if plaintext != "p@ssword with spaces" {
		t.Fatalf("decrypt() = %q, want original password", plaintext)
	}
}

func TestNormalizePostgreSQLTargetDefaultsAndValidation(t *testing.T) {
	target, err := normalizeTargetRequest(PostgreSQLTargetRequest{
		Name:     "Production",
		Host:     "db.internal",
		Username: "devopin_monitor",
	}, 0)
	if err != nil {
		t.Fatalf("normalizeTargetRequest() error = %v", err)
	}
	if target.Port != 5432 || target.Database != "postgres" || target.SSLMode != "prefer" || !target.Enabled {
		t.Fatalf("unexpected target defaults: %+v", target)
	}
	if _, err := normalizeTargetRequest(PostgreSQLTargetRequest{Name: "invalid"}, 0); err == nil {
		t.Fatal("normalizeTargetRequest() accepted an incomplete target")
	}
}
