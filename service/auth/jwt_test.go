package auth

import "testing"

func TestCreateJWT(t *testing.T) {
	secret := []byte("testsecret")

	token, err := CreateJWT(secret, 1)
	if err != nil {
		t.Errorf("failed to create JWT: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}
}
