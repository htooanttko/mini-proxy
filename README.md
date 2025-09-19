# Custom Load Balancer / Reverse Proxy

A production-like Go-based load balancer and proxy server mimicking NGINX/HAProxy features.

## Features
- Reverse Proxy with SSL Termination
- Forward Proxy
- Load Balancing: Round Robin, Least Connections, IP Hash, Weighted
- Health Checks
- Security: Basic Auth, Rate Limiting
- Caching: In-memory
- Compression: Gzip
- Networking: HTTP/HTTPS listeners

## Usage
1. Configure via `config.yaml` (TBD, for now hardcoded in config.go)
2. `go run cmd/server/main.go`

## Configuration
Edit `internal/config/config.go` for backends, etc.