# Custom Load Balancer / Reverse Proxy in Go

A production-like Go-based load balancer and proxy server.


## Installation

### With Go Direct
To install the `mini-proxy` tool directly using Go, run the following command:

```bash
go install github.com/dev-hak/mini-proxy/cmd/mini-proxy@latest
```

### macOS
Download and install `mini-proxy` for AMD64 using cURL:

```bash
curl -sL https://github.com/dev-hak/mini-proxy/releases/download/v0.1.0/mini-proxy_0.1.0_macOS_.amd64.tar.gz \
  | tar xz
sudo mv mini-proxy /usr/local/bin/
```

Download and install `mini-proxy` for ARM64 using cURL:

```bash
curl -sL https://github.com/dev-hak/mini-proxy/releases/download/v0.1.0/mini-proxy_0.1.0_macOS_.arm64.tar.gz \
  | tar xz
sudo mv mini-proxy /usr/local/bin/
```



### Windows
Install `mini-proxy` for AMD64 with PowerShell:

```powershell
Invoke-WebRequest https://github.com/dev-hak/mini-proxy/releases/download/v0.1.0/mini-proxy_0.1.0_Windows_.amd64.tar.gz -OutFile mini-proxy.zip
Expand-Archive mini-proxy.zip -DestinationPath .
Move-Item mini-proxy.exe C:\Windows\System32\
```

Install `mini-proxy` for ARM64 with PowerShell:

```powershell
Invoke-WebRequest https://github.com/dev-hak/mini-proxy/releases/download/v0.1.0/mini-proxy_0.1.0_Windows_.arm64.tar.gz -OutFile mini-proxy.zip
Expand-Archive mini-proxy.zip -DestinationPath .
Move-Item mini-proxy.exe C:\Windows\System32\
```

### Linux
Download and install `mini-proxy` for AMD64 using cURL:

```bash
curl -sL https://github.com/dev-hak/mini-proxy/releases/download/v0.1.0/mini-proxy_0.1.0_Linux_.amd64.tar.gz \
  | tar xz
sudo mv mini-proxy /usr/local/bin/
```
Download and install `mini-proxy` for ARM64 using cURL:

```bash
curl -sL https://github.com/dev-hak/mini-proxy/releases/download/v0.1.0/mini-proxy_0.1.0_Linux_.arm64.tar.gz \
  | tar xz
sudo mv mini-proxy /usr/local/bin/
```


## Features
- Reverse Proxy with SSL/TLS Termination
- Forward Proxy
- Load Balancing: Round Robin, Least Connections, IP Hash, Weighted
- Health Checks
- Security: Basic Auth, Rate Limiting
- Caching: In-memory
- Compression: Gzip
- Networking: HTTP/HTTPS listeners

## Usage
1. Run `mini-proxy --template` for configuration template.
2. Configure via `config.json` file.
3. Run `mini-proxy --config config.json`.
4. Access via defined listeners.

## Notes

- SSL/TLS certificates are required for HTTPS termination.

- Rate limiting and basic authentication can be adjusted in config.json.
