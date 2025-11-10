package engine

import (
	"context"
	"fmt"
	"log"
	"time"
)

// CommandInterceptor intercepts command execution to add cross-cutting concerns.
// Interceptors form a chain of responsibility pattern.
type CommandInterceptor interface {
	// Execute runs the command, potentially delegating to the next interceptor
	Execute(ctx context.Context, command Command[any], executor CommandExecutor) (any, error)

	// GetNext returns the next interceptor in the chain
	GetNext() CommandInterceptor

	// SetNext sets the next interceptor in the chain
	SetNext(next CommandInterceptor)
}

// BaseCommandInterceptor provides a base implementation for interceptors
type BaseCommandInterceptor struct {
	next CommandInterceptor
}

// GetNext returns the next interceptor in the chain
func (i *BaseCommandInterceptor) GetNext() CommandInterceptor {
	return i.next
}

// SetNext sets the next interceptor in the chain
func (i *BaseCommandInterceptor) SetNext(next CommandInterceptor) {
	i.next = next
}

// Execute delegates to the next interceptor
func (i *BaseCommandInterceptor) Execute(ctx context.Context, command Command[any], executor CommandExecutor) (any, error) {
	if i.next != nil {
		return i.next.Execute(ctx, command, executor)
	}
	return executor.Execute(ctx, command)
}

// LoggingInterceptor logs command execution
type LoggingInterceptor struct {
	BaseCommandInterceptor
	logger *log.Logger
}

// NewLoggingInterceptor creates a new logging interceptor
func NewLoggingInterceptor() *LoggingInterceptor {
	return &LoggingInterceptor{
		logger: log.Default(),
	}
}

// Execute logs command execution
func (i *LoggingInterceptor) Execute(ctx context.Context, command Command[any], executor CommandExecutor) (any, error) {
	commandName := fmt.Sprintf("%T", command)
	i.logger.Printf("[FlowGo] Executing command: %s", commandName)

	start := time.Now()
	result, err := i.next.Execute(ctx, command, executor)
	duration := time.Since(start)

	if err != nil {
		i.logger.Printf("[FlowGo] Command %s failed after %v: %v", commandName, duration, err)
		return nil, err
	}

	i.logger.Printf("[FlowGo] Command %s completed successfully in %v", commandName, duration)
	return result, nil
}

// TransactionInterceptor manages transactions for command execution
type TransactionInterceptor struct {
	BaseCommandInterceptor
}

// NewTransactionInterceptor creates a new transaction interceptor
func NewTransactionInterceptor() *TransactionInterceptor {
	return &TransactionInterceptor{}
}

// Execute wraps command execution in a transaction
func (i *TransactionInterceptor) Execute(ctx context.Context, command Command[any], executor CommandExecutor) (any, error) {
	// TODO: Begin transaction
	// tx, err := beginTransaction()
	// if err != nil {
	//     return nil, err
	// }

	result, err := i.next.Execute(ctx, command, executor)

	if err != nil {
		// TODO: Rollback transaction
		// tx.Rollback()
		return nil, err
	}

	// TODO: Commit transaction
	// if err := tx.Commit(); err != nil {
	//     return nil, err
	// }

	return result, nil
}

// ContextInterceptor manages the CommandContext lifecycle
type ContextInterceptor struct {
	BaseCommandInterceptor
	engine *ProcessEngineImpl
}

// NewContextInterceptor creates a new context interceptor
func NewContextInterceptor(engine *ProcessEngineImpl) *ContextInterceptor {
	return &ContextInterceptor{
		engine: engine,
	}
}

// Execute creates and manages the CommandContext
func (i *ContextInterceptor) Execute(ctx context.Context, command Command[any], executor CommandExecutor) (any, error) {
	// Create command context
	commandContext := NewCommandContext(ctx, i.engine)
	defer commandContext.Close()

	// Store in context for access by command
	ctx = context.WithValue(ctx, commandContextKey, commandContext)

	// Execute command with context
	result, err := command.Execute(ctx, commandContext)
	if err != nil {
		commandContext.SetException(err)
		return nil, err
	}

	commandContext.SetResult(result)
	return result, nil
}

// RetryInterceptor provides retry logic for failed commands
type RetryInterceptor struct {
	BaseCommandInterceptor
	maxRetries int
	retryDelay time.Duration
}

// NewRetryInterceptor creates a new retry interceptor
func NewRetryInterceptor(maxRetries int, retryDelay time.Duration) *RetryInterceptor {
	return &RetryInterceptor{
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// Execute retries command execution on failure
func (i *RetryInterceptor) Execute(ctx context.Context, command Command[any], executor CommandExecutor) (any, error) {
	var result any
	var err error

	for attempt := 0; attempt <= i.maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("[FlowGo] Retrying command (attempt %d/%d)", attempt, i.maxRetries)
			time.Sleep(i.retryDelay)
		}

		result, err = i.next.Execute(ctx, command, executor)
		if err == nil {
			return result, nil
		}

		// Check if error is retryable
		if !isRetryableError(err) {
			break
		}
	}

	return nil, fmt.Errorf("command failed after %d retries: %w", i.maxRetries, err)
}

// isRetryableError checks if an error should trigger a retry
func isRetryableError(err error) bool {
	// TODO: Implement proper error classification
	// For now, we don't retry any errors
	return false
}

// commandContextKey is the key for storing CommandContext in context.Context
type contextKey string

const commandContextKey contextKey = "commandContext"

// GetCommandContext retrieves the CommandContext from a context.Context
func GetCommandContext(ctx context.Context) *CommandContext {
	if commandContext, ok := ctx.Value(commandContextKey).(*CommandContext); ok {
		return commandContext
	}
	return nil
}
