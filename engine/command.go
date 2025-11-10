package engine

import "context"

// Command represents an operation that can be executed by the process engine.
// All operations on the engine are encapsulated as commands.
type Command[T any] interface {
	// Execute runs the command and returns a result of type T
	Execute(ctx context.Context, commandContext *CommandContext) (T, error)
}

// CommandExecutor is responsible for executing commands.
// It manages the command execution lifecycle including interceptors and transaction management.
type CommandExecutor interface {
	// Execute runs a command through the interceptor chain
	Execute(ctx context.Context, command Command[any]) (any, error)
}

// CommandContext holds the context information for command execution.
// It provides access to engine services and manages resources during command execution.
type CommandContext struct {
	// Context is the Go context for cancellation and timeout
	Context context.Context

	// Engine is the process engine instance
	Engine *ProcessEngineImpl

	// Session holds the current database session/transaction
	Session interface{}

	// Attributes stores custom attributes for this command execution
	Attributes map[string]interface{}

	// Exception stores any exception that occurred during command execution
	Exception error

	// Result stores the command execution result
	Result interface{}
}

// NewCommandContext creates a new command context
func NewCommandContext(ctx context.Context, engine *ProcessEngineImpl) *CommandContext {
	return &CommandContext{
		Context:    ctx,
		Engine:     engine,
		Attributes: make(map[string]interface{}),
	}
}

// GetAttribute retrieves an attribute from the context
func (c *CommandContext) GetAttribute(key string) interface{} {
	return c.Attributes[key]
}

// SetAttribute sets an attribute in the context
func (c *CommandContext) SetAttribute(key string, value interface{}) {
	c.Attributes[key] = value
}

// HasException returns true if an exception occurred during command execution
func (c *CommandContext) HasException() bool {
	return c.Exception != nil
}

// GetException returns the exception that occurred during command execution
func (c *CommandContext) GetException() error {
	return c.Exception
}

// SetException sets an exception in the command context
func (c *CommandContext) SetException(err error) {
	c.Exception = err
}

// GetResult returns the command execution result
func (c *CommandContext) GetResult() interface{} {
	return c.Result
}

// SetResult sets the command execution result
func (c *CommandContext) SetResult(result interface{}) {
	c.Result = result
}

// Close releases resources associated with this command context
func (c *CommandContext) Close() error {
	// Clean up session/transaction if any
	if c.Session != nil {
		// TODO: Close database session
	}
	return nil
}
