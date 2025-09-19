package security

import (
	"net/http"
	"strings"
)

type Auth struct {
	users map[string]string
}

func NewAuth(users map[string]string) *Auth {
	return &Auth{users: users}
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(auth, "Basic ") {
			http.Error(w, "Invalid auth", http.StatusUnauthorized)
			return
		}
		// Decode base64 (simplified, use proper b64 in prod)
		payload := strings.TrimPrefix(auth, "Basic ")
		// Assume payload is "user:pass" base64 encoded
		// In prod, use golang.org/x/crypto/bcrypt or similar
		for u, p := range a.users {
			if payload == u+":"+p { // Raw for simplicity
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	})
}
