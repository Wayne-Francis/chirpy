package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

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

func TestTokenvalid(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	tokenString, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	gotID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotID != userID {
		t.Fatal("expected validated token")
	}
}

func TestTokenexpired(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	tokenString, err := MakeJWT(userID, secret, time.Millisecond)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}

}

func TestTokenwrongsecret(t *testing.T) {
	userID := uuid.New()
	tokenString, err := MakeJWT(userID, "secreta", time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	_, err = ValidateJWT(tokenString, "secretb")
	if err == nil {
		t.Fatalf("expected error for wrong secret, got nil")
	}

}
