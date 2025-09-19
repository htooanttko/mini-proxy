package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/dev-hak/mini-proxy/pkg/models"
)

type Config struct {
	ListenAddr   string
	TLSCertFile  string
	TLSKeyFile   string
	ProxyType    string // "reverse" or "forward"
	BalancerType string // "round_robin", "least_conn", "ip_hash", "weighted"
	Backends     []models.Backend
	HealthCheck  models.HealthCheckConfig
	Security     SecurityConfig
	Cache        CacheConfig
}

type SecurityConfig struct {
	BasicAuthUsers map[string]string // username:password
	RateLimit      RateLimitConfig
}

type RateLimitConfig struct {
	RequestsPerMin int
	BlockedTime    time.Duration
}

type CacheConfig struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

var DefaultConfig = Config{
	ListenAddr:   ":8080",       // Default listener address
	TLSCertFile:  "cert.pem",    // Default certificate file
	TLSKeyFile:   "key.pem",     // Default key file
	BalancerType: "round_robin", // Default balancer type
	Backends: []models.Backend{
		{URL: "http://backend1:3000", Weight: 1, MaxConnections: 100},
		{URL: "http://backend2:3000", Weight: 2, MaxConnections: 100},
	},
	HealthCheck: models.HealthCheckConfig{
		Interval: time.Second * 30,
	},
	Cache: CacheConfig{
		DefaultExpiration: time.Minute * 5,
		CleanupInterval:   time.Minute,
	},
}

func LoadConfig() *Config {
	// Check for config.json or config.yaml
	if _, err := os.Stat("config.json"); err == nil {
		file, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("Failed to open config.json: %v", err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		var cfg Config
		if err := decoder.Decode(&cfg); err != nil {
			log.Fatalf("Failed to parse config.json: %v", err)
		}
		return &cfg
	} else if _, err := os.Stat("config.yaml"); err == nil {
		file, err := os.Open("config.yaml")
		if err != nil {
			log.Fatalf("Failed to open config.yaml: %v", err)
		}
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		var cfg Config
		if err := decoder.Decode(&cfg); err != nil {
			log.Fatalf("Failed to parse config.yaml: %v", err)
		}
		return &cfg
	}

	// Fallback to default config
	log.Println("No config file found, using default configuration")
	return &DefaultConfig
}

func LoadConfigWithPath(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", path, err)
	}
	defer file.Close()

	if strings.HasSuffix(path, ".json") {
		decoder := json.NewDecoder(file)
		var cfg Config
		if err := decoder.Decode(&cfg); err != nil {
			log.Fatalf("Failed to parse %s: %v", path, err)
		}
		return &cfg
	} else if strings.HasSuffix(path, ".yaml") {
		decoder := yaml.NewDecoder(file)
		var cfg Config
		if err := decoder.Decode(&cfg); err != nil {
			log.Fatalf("Failed to parse %s: %v", path, err)
		}
		return &cfg
	}

	log.Fatalf("Unsupported config file format: %s", path)
	return nil
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
