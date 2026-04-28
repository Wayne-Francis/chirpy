package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "mySecretPassword123"

	// act
	hash, err := HashPassword(password)

	// assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == "" {
		t.Fatal("expected a hash, got empty string")
	}
}

func TestHashPasswordMatch(t *testing.T) {
	password := "mySecretPassword123"
	false_password := "thisiswrong"
	// act
	hash, err := HashPassword(password)
	// assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == "" {
		t.Fatal("expected a hash, got empty string")
	}
	goodResult, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !goodResult {
		t.Fatal("passwords should match")
	}
	badMatch, err := CheckPasswordHash(false_password, hash)
	if badMatch {
		t.Fatal("Passwords should not match")
	}
}
