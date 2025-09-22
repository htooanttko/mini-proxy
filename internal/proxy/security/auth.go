package security

import (
	"encoding/base64"
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
		// Decode base64
		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
		if err != nil {
			http.Error(w, "Invalid auth", http.StatusUnauthorized)
			return
		}
		// Split on colon
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Invalid auth", http.StatusUnauthorized)
			return
		}
		// Validate user:pass
		if a.users[pair[0]] != pair[1] {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
