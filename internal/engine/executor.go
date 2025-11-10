package engine

import (
	"context"
	"fmt"
)

// CommandExecutor is responsible for executing commands through an interceptor chain.
type CommandExecutor struct {
	// first is the first interceptor in the chain
	first Interceptor
}

// NewCommandExecutor creates a new command executor with the given interceptors
func NewCommandExecutor(interceptors ...Interceptor) *CommandExecutor {
	if len(interceptors) == 0 {
		panic("at least one interceptor is required")
	}

	// Build the interceptor chain
	for i := 0; i < len(interceptors)-1; i++ {
		interceptors[i].SetNext(interceptors[i+1])
	}

	return &CommandExecutor{
		first: interceptors[0],
	}
}

// Execute runs the command through the interceptor chain
func (e *CommandExecutor) Execute(ctx context.Context, command Command) (interface{}, error) {
	if command == nil {
		return nil, fmt.Errorf("command cannot be nil")
	}

	// Start execution from the first interceptor
	return e.first.Execute(ctx, command, e)
}

// CommandExecutorBuilder helps build a CommandExecutor with default interceptors
type CommandExecutorBuilder struct {
	engine             *Engine
	interceptors       []Interceptor
	enableLogging      bool
	enableTransaction  bool
	enableRetry        bool
	retryAttempts      int
}

// NewCommandExecutorBuilder creates a new builder
func NewCommandExecutorBuilder(engine *Engine) *CommandExecutorBuilder {
	return &CommandExecutorBuilder{
		engine:            engine,
		interceptors:      make([]Interceptor, 0),
		enableLogging:     true,
		enableTransaction: true,
		enableRetry:       false,
		retryAttempts:     3,
	}
}

// WithLogging enables or disables logging interceptor
func (b *CommandExecutorBuilder) WithLogging(enabled bool) *CommandExecutorBuilder {
	b.enableLogging = enabled
	return b
}

// WithTransaction enables or disables transaction interceptor
func (b *CommandExecutorBuilder) WithTransaction(enabled bool) *CommandExecutorBuilder {
	b.enableTransaction = enabled
	return b
}

// WithRetry enables retry interceptor with specified attempts
func (b *CommandExecutorBuilder) WithRetry(enabled bool, attempts int) *CommandExecutorBuilder {
	b.enableRetry = enabled
	b.retryAttempts = attempts
	return b
}

// AddInterceptor adds a custom interceptor
func (b *CommandExecutorBuilder) AddInterceptor(interceptor Interceptor) *CommandExecutorBuilder {
	b.interceptors = append(b.interceptors, interceptor)
	return b
}

// Build creates the CommandExecutor
func (b *CommandExecutorBuilder) Build() *CommandExecutor {
	interceptors := make([]Interceptor, 0)

	// Add logging interceptor first (outermost)
	if b.enableLogging {
		interceptors = append(interceptors, NewLoggingInterceptor())
	}

	// Add retry interceptor
	if b.enableRetry {
		interceptors = append(interceptors, NewRetryInterceptor(b.retryAttempts))
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
