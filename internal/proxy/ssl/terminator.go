package ssl

import (
	"log"
	"net/http"
)

type SSLTerminator struct {
	certFile string
	keyFile  string
	handler  http.Handler
}

func NewSSLTerminator(cert, key string, h http.Handler) *SSLTerminator {
	return &SSLTerminator{certFile: cert, keyFile: key, handler: h}
}

func (s *SSLTerminator) ListenAndServe(addr string) error {
	if s.certFile == "" || s.keyFile == "" {
		log.Println("SSL not configured, falling back to HTTP")
		return http.ListenAndServe(addr, s.handler)
	}
	log.Println("Starting HTTPS server")
	return http.ListenAndServeTLS(addr, s.certFile, s.keyFile, s.handler)
}

func (s *SSLTerminator) RedirectToHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			u := *r.URL
			u.Scheme = "https"
			u.Host = r.Host
			http.Redirect(w, r, u.String(), http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *SSLTerminator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
