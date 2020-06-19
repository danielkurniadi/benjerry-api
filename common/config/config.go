package config

import "os"

// Config serves standard App Configuration
type Config struct {
	// Runnning Application
	AppName AppIdentifier

	// Server Port Address
	PortAddr string

	// Running Environment
	EnvironmentMode EnvIdentifier

	// Database Configuration
	DatabaseURI   string
	SigningSecret string
}

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
func Get(appID AppIdentifier, port string) *Config {
	return &Config{
		AppName:         appID,
		PortAddr:        port,
		EnvironmentMode: EnvIdentifier(os.Getenv("ENV_MODE")),
		DatabaseURI:     os.Getenv("DB_URI"),
		SigningSecret:   os.Getenv("SIGNING_SECRET"),
	}
}
