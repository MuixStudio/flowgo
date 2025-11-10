package engine

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Interceptor intercepts command execution to add cross-cutting concerns.
type Interceptor interface {
	// Execute runs the command, potentially delegating to the next interceptor
	Execute(ctx context.Context, command Command, executor *CommandExecutor) (interface{}, error)

	// GetNext returns the next interceptor in the chain
	GetNext() Interceptor

	// SetNext sets the next interceptor in the chain
	SetNext(next Interceptor)
}

// BaseInterceptor provides a base implementation for interceptors
type BaseInterceptor struct {
	next Interceptor
}

// GetNext returns the next interceptor in the chain
func (i *BaseInterceptor) GetNext() Interceptor {
	return i.next
}

// SetNext sets the next interceptor in the chain
func (i *BaseInterceptor) SetNext(next Interceptor) {
	i.next = next
}

// Execute delegates to the next interceptor
func (i *BaseInterceptor) Execute(ctx context.Context, command Command, executor *CommandExecutor) (interface{}, error) {
	if i.next != nil {
		return i.next.Execute(ctx, command, executor)
	}
	return executor.Execute(ctx, command)
}

// LoggingInterceptor logs command execution
type LoggingInterceptor struct {
	BaseInterceptor
	logger *log.Logger
}

// NewLoggingInterceptor creates a new logging interceptor
func NewLoggingInterceptor() *LoggingInterceptor {
	return &LoggingInterceptor{
		logger: log.Default(),
	}
}

// Execute logs command execution
func (i *LoggingInterceptor) Execute(ctx context.Context, command Command, executor *CommandExecutor) (interface{}, error) {
	commandName := fmt.Sprintf("%T", command)
	i.logger.Printf("[FlowGo] Executing command: %s", commandName)

	start := time.Now()
	result, err := i.next.Execute(ctx, command, executor)
	duration := time.Since(start)

	if err != nil {
		i.logger.Printf("[FlowGo] Command %s failed after %v: %v", commandName, duration, err)
	} else {
		i.logger.Printf("[FlowGo] Command %s completed successfully in %v", commandName, duration)
	}

	return result, err
}

// TransactionInterceptor manages transactions for command execution
type TransactionInterceptor struct {
	BaseInterceptor
}

// NewTransactionInterceptor creates a new transaction interceptor
func NewTransactionInterceptor() *TransactionInterceptor {
	return &TransactionInterceptor{}
}

// Execute wraps command execution in a transaction
func (i *TransactionInterceptor) Execute(ctx context.Context, command Command, executor *CommandExecutor) (interface{}, error) {
	// TODO: Begin transaction
	result, err := i.next.Execute(ctx, command, executor)

	if err != nil {
		// TODO: Rollback transaction
		return nil, err
	}

	// TODO: Commit transaction
	return result, nil
}

// ContextInterceptor manages the CommandContext lifecycle
type ContextInterceptor struct {
	BaseInterceptor
	engine *Engine
}

// NewContextInterceptor creates a new context interceptor
func NewContextInterceptor(engine *Engine) *ContextInterceptor {
	return &ContextInterceptor{
		engine: engine,
	}
}

// Execute creates and manages the CommandContext
func (i *ContextInterceptor) Execute(ctx context.Context, command Command, executor *CommandExecutor) (interface{}, error) {
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
	BaseInterceptor
	maxRetries int
}

// NewRetryInterceptor creates a new retry interceptor
func NewRetryInterceptor(maxRetries int) *RetryInterceptor {
	return &RetryInterceptor{
		maxRetries: maxRetries,
	}
}

// Execute retries command execution on failure
func (i *RetryInterceptor) Execute(ctx context.Context, command Command, executor *CommandExecutor) (interface{}, error) {
	var result interface{}
	var err error

	for attempt := 0; attempt <= i.maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("[FlowGo] Retrying command (attempt %d/%d)", attempt, i.maxRetries)
			time.Sleep(100 * time.Millisecond)
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
	return false
}

// CommandInvoker is the final interceptor that actually executes the command
type CommandInvoker struct {
	BaseInterceptor
}

// NewCommandInvoker creates a new command invoker
func NewCommandInvoker() *CommandInvoker {
	return &CommandInvoker{}
}

// Execute actually executes the command
func (i *CommandInvoker) Execute(ctx context.Context, command Command, executor *CommandExecutor) (interface{}, error) {
	// Get the command context from the context
	commandContext := GetCommandContext(ctx)
	if commandContext == nil {
		return nil, fmt.Errorf("command context not found in context")
	}

	// Execute the command
	return command.Execute(ctx, commandContext)
}
