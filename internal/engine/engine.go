package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/muixstudio/flowgo/api/history"
	"github.com/muixstudio/flowgo/api/repository"
	"github.com/muixstudio/flowgo/api/runtime"
	"github.com/muixstudio/flowgo/api/task"
	internalRepo "github.com/muixstudio/flowgo/internal/repository"
)

// Engine is the internal implementation of ProcessEngine
type Engine struct {
	config            *Configuration
	repositoryService repository.Service
	runtimeService    runtime.Service
	taskService       task.Service
	historyService    history.Service
	commandExecutor   *CommandExecutor
	running           bool
	mu                sync.RWMutex
}

// Configuration holds the engine configuration
type Configuration struct {
	// EngineName is the name of the engine instance
	EngineName string

	// DatabaseDriver is the database driver to use
	DatabaseDriver string

	// DatabaseURL is the connection string
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

// NewEngine creates a new engine implementation
func NewEngine(config *Configuration) (*Engine, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	e := &Engine{
		config:  config,
		running: false,
	}

	// Initialize command executor
	e.commandExecutor = NewCommandExecutorBuilder(e).
		WithLogging(true).
		WithTransaction(true).
		Build()

	// Initialize services
	if err := e.initializeServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	return e, nil
}

// initializeServices initializes all engine services
func (e *Engine) initializeServices() error {
	// Initialize repository service
	repoService := internalRepo.NewService(e.config.DatabaseDriver, e.config.DatabaseURL)
	e.repositoryService = repoService

	// TODO: Initialize other services
	// e.runtimeService = runtime.NewService(e.repositoryService, e.config.EnableAsync)
	// e.taskService = task.NewService(e.runtimeService)
	// if e.config.EnableHistory {
	//     e.historyService = history.NewService(e.config.DatabaseDriver, e.config.DatabaseURL)
	// }

	return nil
}

// GetRepositoryService returns the repository service
func (e *Engine) GetRepositoryService() repository.Service {
	return e.repositoryService
}

// GetRuntimeService returns the runtime service
func (e *Engine) GetRuntimeService() runtime.Service {
	return e.runtimeService
}

// GetTaskService returns the task service
func (e *Engine) GetTaskService() task.Service {
	return e.taskService
}

// GetHistoryService returns the history service
func (e *Engine) GetHistoryService() history.Service {
	return e.historyService
}

// Execute executes a command through the command executor
func (e *Engine) Execute(ctx context.Context, command Command) (interface{}, error) {
	if !e.IsRunning() {
		return nil, fmt.Errorf("engine '%s' is not running", e.config.EngineName)
	}
	return e.commandExecutor.Execute(ctx, command)
}

// Start initializes and starts the process engine
func (e *Engine) Start(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return fmt.Errorf("engine '%s' is already running", e.config.EngineName)
	}

	// TODO: Start all services
	// For now, just mark as running
	e.running = true
	return nil
}

// Stop gracefully shuts down the process engine
func (e *Engine) Stop(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.running {
		return fmt.Errorf("engine '%s' is not running", e.config.EngineName)
	}

	// TODO: Stop all services
	e.running = false
	return nil
}

// GetName returns the name of this process engine
func (e *Engine) GetName() string {
	return e.config.EngineName
}

// IsRunning returns whether the engine is currently running
func (e *Engine) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.running
}

// GetConfiguration returns the engine configuration
func (e *Engine) GetConfiguration() *Configuration {
	return e.config
}
