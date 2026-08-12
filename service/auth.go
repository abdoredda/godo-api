package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	secret []byte
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{secret: []byte(secret)}
}

func (s *AuthService) GenerateToken(userId int) (string, error) {
	claims := Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("error when trying to use [SignedString] err= %w", err)
	}
	return signedToken, nil
}

func (s *AuthService) ParseToken(token string) (Claims, error) {
	claims := Claims{}
	keyFunc := func(tokenString *jwt.Token) (any, error) {
		if _, ok := tokenString.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tokenString.Header["alg"])
		}
		return s.secret, nil
	}

	_, err := jwt.ParseWithClaims(token, &claims, keyFunc)
	if err != nil {
		return Claims{}, fmt.Errorf("error while parsing the token error = %w", err)
	}

	return claims, nil
}
