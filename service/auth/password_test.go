package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashedPassword("password123")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Error("expected non-empty hash")
	}

	if hash == "password123" {
		t.Error("expected hash to not equal plain password")
	}
}

func TestComparePasswords(t *testing.T) {
	hash, err := HashedPassword("password123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !ComparePasswords(hash, []byte("password123")) {
		t.Errorf("expected no error, got %v", err)
	}

	if ComparePasswords(hash, []byte("wrongpassword")) {
		t.Error("expected error for wrong password comparison")
	}
}
