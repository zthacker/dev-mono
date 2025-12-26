package config

import (
	"time"
)

// Config holds all server configuration.
// Load from environment variables, config file, or flags.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	NATS     NATSConfig
	Game     GameConfig
}

type ServerConfig struct {
	// Network binding
	BindAddr string // e.g., ":8085"

	// Identity
	ServerID uint32   // Unique ID for this server instance
	ZoneIDs  []uint32 // Which zones this server handles

	// Limits
	MaxConnections int
	MaxPlayers     int
}

type DatabaseConfig struct {
	// MySQL connection
	Host     string
	Port     int
	User     string
	Password string
	Database string

	// Connection pool
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type NATSConfig struct {
	URL string // e.g., "nats://localhost:4222"

	// Reconnection
	MaxReconnects int
	ReconnectWait time.Duration
}

type GameConfig struct {
	// Tick rate
	TickRate int // Ticks per second (default 20)

	// View distance
	ViewDistance float32 // Units (default 100)

	// Combat
	MaxLevel uint8
}

// TODO: Implement configuration loading:
//
// func Load() (*Config, error)
//   Priority (highest to lowest):
//   1. Environment variables (REALM_SERVER_BIND_ADDR, etc.)
//   2. Config file (config.yaml or config.json)
//   3. Defaults
//
// func LoadFromFile(path string) (*Config, error)
//
// func (c *Config) Validate() error
//   - Check required fields
//   - Validate ranges

// DefaultConfig returns sensible defaults for development.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			BindAddr:       ":8085",
			ServerID:       1,
			ZoneIDs:        []uint32{1}, // Default zone
			MaxConnections: 1000,
			MaxPlayers:     500,
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            3306,
			User:            "realm",
			Password:        "realm",
			Database:        "realm_characters",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		NATS: NATSConfig{
			URL:           "nats://localhost:4222",
			MaxReconnects: 10,
			ReconnectWait: 2 * time.Second,
		},
		Game: GameConfig{
			TickRate:     20,
			ViewDistance: 100,
			MaxLevel:     60,
		},
	}
}

// DSN returns the MySQL connection string.
func (c *DatabaseConfig) DSN() string {
	// TODO: Implement
	// return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
	//     c.User, c.Password, c.Host, c.Port, c.Database)
	return ""
}
