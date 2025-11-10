package engine

import (
	"context"
	"fmt"
)

// CommandExecutorImpl is the default implementation of CommandExecutor
type CommandExecutorImpl struct {
	// first is the first interceptor in the chain
	first CommandInterceptor

	// last is the last interceptor in the chain (typically the one that executes the command)
	last CommandInterceptor
}

// NewCommandExecutor creates a new command executor with the given interceptors
func NewCommandExecutor(interceptors ...CommandInterceptor) *CommandExecutorImpl {
	if len(interceptors) == 0 {
		panic("at least one interceptor is required")
	}

	// Build the interceptor chain
	for i := 0; i < len(interceptors)-1; i++ {
		interceptors[i].SetNext(interceptors[i+1])
	}

	return &CommandExecutorImpl{
		first: interceptors[0],
		last:  interceptors[len(interceptors)-1],
	}
}

// Execute runs the command through the interceptor chain
func (e *CommandExecutorImpl) Execute(ctx context.Context, command Command[any]) (any, error) {
	if command == nil {
		return nil, fmt.Errorf("command cannot be nil")
	}

	// Start execution from the first interceptor
	return e.first.Execute(ctx, command, e)
}

// CommandInvoker is the final interceptor that actually executes the command
type CommandInvoker struct {
	BaseCommandInterceptor
}

// NewCommandInvoker creates a new command invoker
func NewCommandInvoker() *CommandInvoker {
	return &CommandInvoker{}
}

// Execute actually executes the command
func (i *CommandInvoker) Execute(ctx context.Context, command Command[any], executor CommandExecutor) (any, error) {
	// Get the command context from the context
	commandContext := GetCommandContext(ctx)
	if commandContext == nil {
		return nil, fmt.Errorf("command context not found in context")
	}

	// Execute the command
	return command.Execute(ctx, commandContext)
}

// DefaultCommandExecutorBuilder helps build a CommandExecutor with default interceptors
type DefaultCommandExecutorBuilder struct {
	engine            *ProcessEngineImpl
	interceptors      []CommandInterceptor
	enableLogging     bool
	enableTransaction bool
	enableRetry       bool
	retryAttempts     int
}

// NewDefaultCommandExecutorBuilder creates a new builder
func NewDefaultCommandExecutorBuilder(engine *ProcessEngineImpl) *DefaultCommandExecutorBuilder {
	return &DefaultCommandExecutorBuilder{
		engine:            engine,
		interceptors:      make([]CommandInterceptor, 0),
		enableLogging:     true,
		enableTransaction: true,
		enableRetry:       false,
		retryAttempts:     3,
	}
}

// WithLogging enables or disables logging interceptor
func (b *DefaultCommandExecutorBuilder) WithLogging(enabled bool) *DefaultCommandExecutorBuilder {
	b.enableLogging = enabled
	return b
}

// WithTransaction enables or disables transaction interceptor
func (b *DefaultCommandExecutorBuilder) WithTransaction(enabled bool) *DefaultCommandExecutorBuilder {
	b.enableTransaction = enabled
	return b
}

// WithRetry enables retry interceptor with specified attempts
func (b *DefaultCommandExecutorBuilder) WithRetry(enabled bool, attempts int) *DefaultCommandExecutorBuilder {
	b.enableRetry = enabled
	b.retryAttempts = attempts
	return b
}

// AddInterceptor adds a custom interceptor
func (b *DefaultCommandExecutorBuilder) AddInterceptor(interceptor CommandInterceptor) *DefaultCommandExecutorBuilder {
	b.interceptors = append(b.interceptors, interceptor)
	return b
}

// Build creates the CommandExecutor
func (b *DefaultCommandExecutorBuilder) Build() *CommandExecutorImpl {
	interceptors := make([]CommandInterceptor, 0)

	// Add logging interceptor first (outermost)
	if b.enableLogging {
		interceptors = append(interceptors, NewLoggingInterceptor())
	}

	// Add retry interceptor
	if b.enableRetry {
		interceptors = append(interceptors, NewRetryInterceptor(b.retryAttempts, 0))
	}

	// Add custom interceptors
	interceptors = append(interceptors, b.interceptors...)

	// Add transaction interceptor
	if b.enableTransaction {
		interceptors = append(interceptors, NewTransactionInterceptor())
	}

	// Add context interceptor (must be before invoker)
	interceptors = append(interceptors, NewContextInterceptor(b.engine))

	// Add command invoker last (innermost)
	interceptors = append(interceptors, NewCommandInvoker())

	return NewCommandExecutor(interceptors...)
}
