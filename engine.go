package flowgo

import (
	"context"

	"github.com/muixstudio/flowgo/api/history"
	"github.com/muixstudio/flowgo/api/repository"
	"github.com/muixstudio/flowgo/api/runtime"
	"github.com/muixstudio/flowgo/api/task"
	"github.com/muixstudio/flowgo/internal/engine"
)

// ProcessEngine is the main entry point for the FlowGo workflow engine.
// It provides access to all core services and manages the engine lifecycle.
type ProcessEngine interface {
	// GetRepositoryService returns the repository service for managing process definitions
	GetRepositoryService() repository.Service

	// GetRuntimeService returns the runtime service for managing process instances
	GetRuntimeService() runtime.Service

	// GetTaskService returns the task service for managing user tasks
	GetTaskService() task.Service

	// GetHistoryService returns the history service for querying historical data
	GetHistoryService() history.Service

	// Start initializes and starts the process engine
	Start(ctx context.Context) error

	// Stop gracefully shuts down the process engine
	Stop(ctx context.Context) error

	// GetName returns the name of this process engine
	GetName() string

	// IsRunning returns whether the engine is currently running
	IsRunning() bool
}

// NewProcessEngine creates a new ProcessEngine with the given configuration.
// This is the primary way to create a process engine instance.
func NewProcessEngine(config *Configuration) (ProcessEngine, error) {
	internalConfig := &engine.Configuration{
		EngineName:     config.EngineName,
		DatabaseDriver: config.DatabaseDriver,
		DatabaseURL:    config.DatabaseURL,
		EnableHistory:  config.EnableHistory,
		EnableAsync:    config.EnableAsync,
		MaxPoolSize:    config.MaxPoolSize,
		IdleTimeout:    config.IdleTimeout,
	}
	return engine.NewEngine(internalConfig)
}

// NewProcessEngineBuilder creates a new builder for constructing a process engine.
// This provides a fluent API for engine configuration.
func NewProcessEngineBuilder() *Builder {
	return &Builder{
		config: DefaultConfiguration(),
	}
}
