package config

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

const (
	dockerComposePgConn = "postgres://fediverse:fediverse@localhost:5432/fediverse"
	viteDevServerURL    = "http://localhost:5173"
)

// Config contains the configuration for the application.
// Shared between the API and the crawler.
type Config struct {
	Environment Environment `help:"Environment to run in." enum:"development,production" default:"development" env:"ENVIRONMENT"`
	LogLevel    LogLevel    `help:"Log level." enum:"debug,info,warn,error" default:"info" env:"LOG_LEVEL"`

	// DB
	PgConn string `help:"Postgres connection string." env:"PG_CONN"`

	// API
	FrontendURL string `help:"URL of the frontend." env:"FRONTEND_URL"`
}

// SetDevelopmentDefaults sets the defaults for the development environment.
func (c *Config) SetDevelopmentDefaults() {
	if c.Environment != EnvironmentDevelopment {
		return
	}

	if c.PgConn == "" {
		c.PgConn = dockerComposePgConn
	}

	if c.FrontendURL == "" {
		c.FrontendURL = viteDevServerURL
	}
}
