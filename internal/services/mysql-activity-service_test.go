package service

import "testing"

func TestMySQLCredentialEncryptionRoundTrip(t *testing.T) {
	service := &mysqlActivityService{secret: "test-secret"}
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

func TestNormalizeMySQLTargetDefaultsAndValidation(t *testing.T) {
	target, err := normalizeMySQLTargetRequest(MySQLTargetRequest{
		Name:     "Production",
		Host:     "db.internal",
		Username: "devopin_monitor",
	}, 0)
	if err != nil {
		t.Fatalf("normalizeMySQLTargetRequest() error = %v", err)
	}
	if target.Port != 3306 || target.Database != "mysql" || target.SSLMode != "preferred" || !target.Enabled {
		t.Fatalf("unexpected target defaults: %+v", target)
	}
	if _, err := normalizeMySQLTargetRequest(MySQLTargetRequest{Name: "invalid"}, 0); err == nil {
		t.Fatal("normalizeMySQLTargetRequest() accepted an incomplete target")
	}
}
