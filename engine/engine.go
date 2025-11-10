package engine

import (
	"context"

	"github.com/muixstudio/flowgo/history"
	"github.com/muixstudio/flowgo/repository"
	"github.com/muixstudio/flowgo/runtime"
	"github.com/muixstudio/flowgo/task"
)

// ProcessEngine is the main entry point for the FlowGo workflow engine.
// It provides access to all core services and manages the engine lifecycle.
type ProcessEngine interface {
	// GetRepositoryService returns the repository service for managing process definitions
	GetRepositoryService() repository.RepositoryService

	// GetRuntimeService returns the runtime service for managing process instances
	GetRuntimeService() runtime.RuntimeService

	// GetTaskService returns the task service for managing user tasks
	GetTaskService() task.TaskService

	// GetHistoryService returns the history service for querying historical data
	GetHistoryService() history.HistoryService

	// Execute executes a command through the command executor
	//Execute[T any](ctx context.Context, command Command[T]) (T, error)

	// Start initializes and starts the process engine
	Start(ctx context.Context) error

	// Stop gracefully shuts down the process engine
	Stop(ctx context.Context) error

	// GetName returns the name of this process engine
	GetName() string

	// IsRunning returns whether the engine is currently running
	IsRunning() bool
}

// ProcessEngineConfiguration holds the configuration for creating a ProcessEngine
type ProcessEngineConfiguration struct {
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

	// IdleTimeout is the idle timeout for database connections
	IdleTimeout int
}

// DefaultProcessEngineConfiguration returns a configuration with default values
func DefaultProcessEngineConfiguration() *ProcessEngineConfiguration {
	return &ProcessEngineConfiguration{
		EngineName:     "default",
		DatabaseDriver: "postgres",
		EnableHistory:  true,
		EnableAsync:    true,
		MaxPoolSize:    10,
		IdleTimeout:    300,
	}
}

// ProcessEngineBuilder provides a fluent API for building a ProcessEngine
type ProcessEngineBuilder struct {
	config *ProcessEngineConfiguration
}

// NewProcessEngineBuilder creates a new builder with default configuration
func NewProcessEngineBuilder() *ProcessEngineBuilder {
	return &ProcessEngineBuilder{
		config: DefaultProcessEngineConfiguration(),
	}
}

// WithEngineName sets the engine name
func (b *ProcessEngineBuilder) WithEngineName(name string) *ProcessEngineBuilder {
	b.config.EngineName = name
	return b
}

// WithDatabase sets the database configuration
func (b *ProcessEngineBuilder) WithDatabase(driver, url string) *ProcessEngineBuilder {
	b.config.DatabaseDriver = driver
	b.config.DatabaseURL = url
	return b
}

// WithHistory enables or disables history recording
func (b *ProcessEngineBuilder) WithHistory(enabled bool) *ProcessEngineBuilder {
	b.config.EnableHistory = enabled
	return b
}

// WithAsync enables or disables async execution
func (b *ProcessEngineBuilder) WithAsync(enabled bool) *ProcessEngineBuilder {
	b.config.EnableAsync = enabled
	return b
}

// WithPoolSize sets the database connection pool size
func (b *ProcessEngineBuilder) WithPoolSize(size int) *ProcessEngineBuilder {
	b.config.MaxPoolSize = size
	return b
}

// Build creates and returns a new ProcessEngine instance
func (b *ProcessEngineBuilder) Build() (ProcessEngine, error) {
	return NewProcessEngine(b.config)
}

// NewProcessEngine creates a new ProcessEngine with the given configuration
func NewProcessEngine(config *ProcessEngineConfiguration) (ProcessEngine, error) {
	return newProcessEngineImpl(config)
}
