package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// AppConfig serves standard App Configuration
type AppConfig struct {
	// Runnning Application
	AppName AppIdentifier

	// Server Port Address
	Hostname string
	PortAddr string

	// Running Environment
	EnvironmentMode EnvIdentifier

	// Database Configuration
	DatabaseURI  string
	DatabaseName string
	RedisURI     string
}

// AppAddress returns address of hosted app
// which is hostname:port
func (conf *AppConfig) AppAddress() string { return conf.Hostname + ":" + conf.PortAddr }

type (
	// AppIdentifier provides naming for
	// service instance type
	AppIdentifier string

	// EnvIdentifier tags current running environment
	// see constant below
	EnvIdentifier string
)

// Enum for deployment related
const (
	// Applications
	BENJERRY AppIdentifier = "BenJerry"

	// Environments
	DEVELOPMENT EnvIdentifier = "development"
	STAGING     EnvIdentifier = "staging"
	PRODUCTION  EnvIdentifier = "production"
)

// Get application configurations which
// are passed by environment variables
func Get(appID AppIdentifier, host string, port string) AppConfig {
	dbURI := os.Getenv("DB_URI")
	uri, err := url.Parse(dbURI)

	if err != nil {
		fmt.Println("warning: got invalid Database URI:", err)
	}

	dbName := strings.TrimLeft(uri.Path, "/")
	dbURI = dbURI[:len(dbURI)-len(dbName)]

	redisURI := os.Getenv("REDIS_URI")
	if len(redisURI) == 0 {
		redisURI = "redis://localhost:6379"
	}

	env := EnvIdentifier(os.Getenv("ENV_MODE"))
	if len(env) == 0 {
		env = DEVELOPMENT
	}

	return AppConfig{
		AppName:         appID,
		Hostname:        host,
		PortAddr:        port,
		EnvironmentMode: env,
		DatabaseURI:     dbURI,
		DatabaseName:    dbName,
		RedisURI:        redisURI,
	}
}

// PrintConfig display configurations
func PrintConfig(config AppConfig) {
	format := "%s \t\t: \t%s\n"

	fmt.Println("App Configurations:")
	fmt.Println("-----------------------------------------")

	fmt.Printf(format, "App Name", config.AppName)
	fmt.Printf(format, "Hostname", config.Hostname)
	fmt.Printf(format, "Port Address", config.PortAddr)
	fmt.Printf(format, "Environ Mode", config.EnvironmentMode)
	fmt.Printf(format, "Database URI", config.DatabaseURI)
	fmt.Printf(format, "Database Name", config.DatabaseName)
	fmt.Printf(format, "Redis URI", config.RedisURI)

	fmt.Println("-----------------------------------------")
}
