package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthService_GenerateAndParseToken(t *testing.T) {
	// Arrange
	jwtSecret := "test-secret"

	authService := NewAuthService(jwtSecret)
	userID := 5
	// Act
	token, err := authService.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := authService.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	// Assert

	if claims.UserID != userID {
		t.Errorf("got %v, want %v", claims.UserID, userID)
	}
}

func TestAuthService_ParseTokenOnlyAcceptHS256(t *testing.T) {
	jwtSecret := "test-secret"
	authService := NewAuthService(jwtSecret)

	claims := Claims{
		UserID: 5,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		t.Fatalf("error when trying to use [SignedString] err= %v", err)
	}

	_, err = authService.ParseToken(signedToken)
	if err == nil {
		t.Error("expected ParseToken to reject HS512 token")
	}
}

func TestAuthService_ParseTokenOnlyRejectExpiredTokens(t *testing.T) {
	jwtSecret := "test-secret"
	authService := NewAuthService(jwtSecret)

	claims := Claims{
		UserID: 5,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		t.Fatalf("error when trying to use [SignedString] err= %v", err)
	}

	_, err = authService.ParseToken(signedToken)
	if err == nil {
		t.Error("expected ParseToken to reject expired tokens")
	}
}
