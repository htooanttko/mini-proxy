package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	users := map[string]string{
		"user1": "password1",
		"user2": "password2",
	}
	auth := NewAuth(users)

	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
	}{
		{
			name:           "Valid credentials",
			authorization:  "Basic dXNlcjE6cGFzc3dvcmQx", // user1:password1
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid credentials",
			authorization:  "Basic dXNlcjE6aW52YWxpZA==", // user1:invalid
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing Authorization header",
			authorization:  "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}
		})
	}
}
