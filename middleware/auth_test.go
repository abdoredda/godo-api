package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"task-manager-api/service"
)

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	// Arrange
	authService := service.NewAuthService("test-secret")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := AuthMiddleware(dummyHandler, authService)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()

	// Act
	protectedHandler.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
