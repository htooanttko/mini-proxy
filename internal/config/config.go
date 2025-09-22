package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dev-hak/mini-proxy/pkg/models"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return errors.New("invalid duration format")
	}
	*d = Duration(dur)
	return nil
}

type Config struct {
	ListenAddr   string            `json:"listen_addr"`
	TLSCertFile  string            `json:"tls_cert_file"`
	TLSKeyFile   string            `json:"tls_key_file"`
	ProxyType    string            `json:"proxy_type"`    // "reverse" or "forward"
	BalancerType string            `json:"balancer_type"` // "round_robin", "least_conn", "ip_hash", "weighted"
	Backends     []models.Backend  `json:"backends"`
	HealthCheck  HealthCheckConfig `json:"health_check"`
	Security     SecurityConfig    `json:"security"`
	Cache        CacheConfig       `json:"cache"`
}

type HealthCheckConfig struct {
	Interval Duration `json:"interval"`
	Timeout  Duration `json:"timeout"`
	Path     string   `json:"path"`
}

type SecurityConfig struct {
	BasicAuthUsers map[string]string `json:"basic_auth_users"` // username:password
	RateLimit      RateLimitConfig   `json:"rate_limit"`
}

type RateLimitConfig struct {
	RequestsPerMin int      `json:"requests_per_min"`
	BlockedTime    Duration `json:"blocked_time"`
}

type CacheConfig struct {
	DefaultExpiration Duration `json:"default_expiration"`
	CleanupInterval   Duration `json:"cleanup_interval"`
}

var DefaultConfig = Config{
	ListenAddr:   ":8080",       // Default listener address
	ProxyType:    "reverse",     // Default proxy type
	TLSCertFile:  "cert.pem",    // Default certificate file
	TLSKeyFile:   "key.pem",     // Default key file
	BalancerType: "round_robin", // Default balancer type
	Backends: []models.Backend{
		{URL: "http://backend1:3000", Weight: 1, MaxConnections: 100},
		{URL: "http://backend2:3000", Weight: 2, MaxConnections: 100},
	},
	HealthCheck: HealthCheckConfig{
		Interval: Duration(30 * time.Second),
		Timeout:  Duration(5 * time.Second),
		Path:     "/health",
	},
	Security: SecurityConfig{
		BasicAuthUsers: map[string]string{
			"user": "password", // Default user
		},
		RateLimit: RateLimitConfig{
			RequestsPerMin: 60,                        // Default 60 requests per minute
			BlockedTime:    Duration(1 * time.Minute), // Default block time 1 minute
		},
	},
	Cache: CacheConfig{
		DefaultExpiration: Duration(5 * time.Minute),
		CleanupInterval:   Duration(1 * time.Minute),
	},
}

func LoadConfig() *Config {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	templateFlag := flag.Bool("template", false, "Generate a default config.json template")
	flag.Parse()

	if *templateFlag {
		generateTemplate()
		os.Exit(0)
	}

	return LoadConfigWithPath(*configPath)
}

func LoadConfigWithPath(path string) *Config {
	if path == "" {
		path = "config.json"
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", path, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("Failed to parse %s: %v", path, err)
	}
	return &cfg
}

func generateTemplate() {
	file, err := os.Create("config.json")
	if err != nil {
		log.Fatalf("Failed to create config.json template: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(DefaultConfig); err != nil {
		log.Fatalf("Failed to write config.json template: %v", err)
	}
	log.Println("config.json template generated successfully")
}

func (c *Config) Validate() error {
	if c.ListenAddr == "" {
		return fmt.Errorf("listen addr required")
	}
	if len(c.Backends) == 0 {
		return fmt.Errorf("at least one backend required")
	}
	return nil
}
