package forward

import (
	"io"
	"net"
	"net/http"
	"time"
)

type ForwardProxy struct{}

func NewForwardProxy() *ForwardProxy {
	return &ForwardProxy{}
}

func (fp *ForwardProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		fp.handleConnect(w, req)
	} else {
		fp.handleHTTP(w, req)
	}
}

func (fp *ForwardProxy) handleConnect(w http.ResponseWriter, req *http.Request) {
	targetConn, err := net.DialTimeout("tcp", req.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer targetConn.Close()

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	w.WriteHeader(http.StatusOK)

	go transfer(targetConn, clientConn)
	transfer(clientConn, targetConn)
}

func (fp *ForwardProxy) handleHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Scheme == "" {
		if req.Method == "GET" || req.Method == "HEAD" {
			if _, ok := req.Header["Proxy-Authorization"]; !ok {
				req.URL.Scheme = "http"
			}
		}
	}

	if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
		http.Error(w, "Unsupported protocol", http.StatusBadRequest)
		return
	}

	if req.URL.Host == "" {
		http.Error(w, "Missing host", http.StatusBadRequest)
		return
	}

	transport := &http.Transport{}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func transfer(dst net.Conn, src net.Conn) {
	defer dst.Close()
	defer src.Close()
	io.Copy(dst, src)
}
