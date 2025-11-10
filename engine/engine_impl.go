package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/muixstudio/flowgo/history"
	"github.com/muixstudio/flowgo/repository"
	"github.com/muixstudio/flowgo/runtime"
	"github.com/muixstudio/flowgo/task"
)

// ProcessEngineImpl is the default implementation of ProcessEngine
type ProcessEngineImpl struct {
	config            *ProcessEngineConfiguration
	repositoryService repository.RepositoryService
	runtimeService    runtime.RuntimeService
	taskService       task.TaskService
	historyService    history.HistoryService
	commandExecutor   CommandExecutor
	running           bool
	mu                sync.RWMutex
}

// newProcessEngineImpl creates a new process engine implementation
func newProcessEngineImpl(config *ProcessEngineConfiguration) (*ProcessEngineImpl, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	engine := &ProcessEngineImpl{
		config:  config,
		running: false,
	}

	// Initialize command executor (one instance for all commands)
	engine.commandExecutor = NewDefaultCommandExecutorBuilder(engine).
		WithLogging(true).
		WithTransaction(true).
		Build()

	// Initialize services
	if err := engine.initializeServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	return engine, nil
}

// initializeServices initializes all engine services
func (e *ProcessEngineImpl) initializeServices() error {
	// Initialize repository service
	e.repositoryService = repository.NewRepositoryService(e.config.DatabaseDriver, e.config.DatabaseURL)

	// Initialize runtime service
	e.runtimeService = runtime.NewRuntimeService(e.repositoryService, e.config.EnableAsync)

	// Initialize task service
	e.taskService = task.NewTaskService(e.runtimeService)

	// Initialize history service (if enabled)
	if e.config.EnableHistory {
		e.historyService = history.NewHistoryService(e.config.DatabaseDriver, e.config.DatabaseURL)
	} else {
		e.historyService = history.NewNoOpHistoryService()
	}

	return nil
}

// GetRepositoryService returns the repository service
func (e *ProcessEngineImpl) GetRepositoryService() repository.RepositoryService {
	return e.repositoryService
}

// GetRuntimeService returns the runtime service
func (e *ProcessEngineImpl) GetRuntimeService() runtime.RuntimeService {
	return e.runtimeService
}

// GetTaskService returns the task service
func (e *ProcessEngineImpl) GetTaskService() task.TaskService {
	return e.taskService
}

// GetHistoryService returns the history service
func (e *ProcessEngineImpl) GetHistoryService() history.HistoryService {
	return e.historyService
}

// GetCommandExecutor returns the command executor
func (e *ProcessEngineImpl) GetCommandExecutor() CommandExecutor {
	return e.commandExecutor
}

// ExecuteCommand executes a command through the command executor
// This method accepts Command[any] and returns any (requires type assertion by caller)
func (e *ProcessEngineImpl) ExecuteCommand(ctx context.Context, command Command[any]) (any, error) {
	if !e.IsRunning() {
		return nil, fmt.Errorf("engine '%s' is not running", e.config.EngineName)
	}
	return e.commandExecutor.Execute(ctx, command)
}

// Start initializes and starts the process engine
func (e *ProcessEngineImpl) Start(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return fmt.Errorf("engine '%s' is already running", e.config.EngineName)
	}

	// Start all services
	if err := e.repositoryService.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to start repository service: %w", err)
	}

	if err := e.runtimeService.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to start runtime service: %w", err)
	}

	if err := e.taskService.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to start task service: %w", err)
	}

	if e.config.EnableHistory {
		if err := e.historyService.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to start history service: %w", err)
		}
	}

	e.running = true
	return nil
}

// Stop gracefully shuts down the process engine
func (e *ProcessEngineImpl) Stop(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.running {
		return fmt.Errorf("engine '%s' is not running", e.config.EngineName)
	}

	// Stop all services in reverse order
	if e.config.EnableHistory {
		if err := e.historyService.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to stop history service: %w", err)
		}
	}

	if err := e.taskService.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop task service: %w", err)
	}

	if err := e.runtimeService.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop runtime service: %w", err)
	}

	if err := e.repositoryService.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop repository service: %w", err)
	}

	e.running = false
	return nil
}

// GetName returns the name of this process engine
func (e *ProcessEngineImpl) GetName() string {
	return e.config.EngineName
}

// IsRunning returns whether the engine is currently running
func (e *ProcessEngineImpl) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.running
}

// GetConfiguration returns the engine configuration
func (e *ProcessEngineImpl) GetConfiguration() *ProcessEngineConfiguration {
	return e.config
}
