package flowgo

// Configuration holds the configuration for creating a ProcessEngine.
// All fields are exported to allow direct manipulation if needed.
type Configuration struct {
	// EngineName is the name of the engine instance
	EngineName string

	// DatabaseDriver is the database driver to use (e.g., "postgres", "mysql")
	DatabaseDriver string

	// DatabaseURL is the connection string for the database
	DatabaseURL string

	// EnableHistory determines if history data should be recorded
	EnableHistory bool

	// EnableAsync determines if async executors should be enabled
	EnableAsync bool

	// MaxPoolSize is the maximum number of database connections
	MaxPoolSize int

	// IdleTimeout is the idle timeout for database connections (in seconds)
	IdleTimeout int
}

// DefaultConfiguration returns a configuration with sensible default values.
func DefaultConfiguration() *Configuration {
	return &Configuration{
		EngineName:     "default",
		DatabaseDriver: "postgres",
		EnableHistory:  true,
		EnableAsync:    true,
		MaxPoolSize:    10,
		IdleTimeout:    300,
	}
}

// Builder provides a fluent API for building a ProcessEngine.
type Builder struct {
	config *Configuration
}

// WithEngineName sets the engine name.
func (b *Builder) WithEngineName(name string) *Builder {
	b.config.EngineName = name
	return b
}

// WithDatabase sets the database configuration.
func (b *Builder) WithDatabase(driver, url string) *Builder {
	b.config.DatabaseDriver = driver
	b.config.DatabaseURL = url
	return b
}

// WithHistory enables or disables history recording.
func (b *Builder) WithHistory(enabled bool) *Builder {
	b.config.EnableHistory = enabled
	return b
}

// WithAsync enables or disables async execution.
func (b *Builder) WithAsync(enabled bool) *Builder {
	b.config.EnableAsync = enabled
	return b
}

// WithPoolSize sets the database connection pool size.
func (b *Builder) WithPoolSize(size int) *Builder {
	b.config.MaxPoolSize = size
	return b
}

// Build creates and returns a new ProcessEngine instance.
func (b *Builder) Build() (ProcessEngine, error) {
	return NewProcessEngine(b.config)
}
