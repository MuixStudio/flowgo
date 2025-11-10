# FlowGo Command Pattern Architecture

## Overview

FlowGo implements the Command Pattern, similar to Flowable/Activiti, to provide a consistent and extensible way to execute operations on the workflow engine.

## Core Concepts

### 1. Command Interface

```go
type Command[T any] interface {
    Execute(ctx context.Context, commandContext *CommandContext) (T, error)
}
```

Every operation in FlowGo is encapsulated as a `Command`. This includes:
- Deploying process definitions
- Starting process instances
- Completing tasks
- Querying data
- And more...

### 2. CommandExecutor

The `CommandExecutor` is responsible for executing commands through an interceptor chain.

```go
type CommandExecutor interface {
    Execute(ctx context.Context, command Command[any]) (any, error)
}
```

### 3. CommandContext

The `CommandContext` holds the execution context for a command, including:
- Reference to the process engine
- Database session/transaction
- Custom attributes
- Execution result and exceptions

```go
type CommandContext struct {
    Context    context.Context
    Engine     *ProcessEngineImpl
    Session    interface{}
    Attributes map[string]interface{}
    Exception  error
    Result     interface{}
}
```

## Interceptor Chain

Interceptors form a chain of responsibility, allowing cross-cutting concerns to be applied to all commands.

### Built-in Interceptors

#### 1. LoggingInterceptor

Logs command execution with timing information.

```
[FlowGo] Executing command: *commands.DeployCommand
[FlowGo] Command *commands.DeployCommand completed successfully in 15ms
```

#### 2. TransactionInterceptor

Manages database transactions for commands.
- Begins transaction before command execution
- Commits on success
- Rolls back on failure

#### 3. ContextInterceptor

Creates and manages the `CommandContext` lifecycle.
- Creates context before execution
- Stores in Go context for access
- Cleans up resources after execution

#### 4. RetryInterceptor

Provides automatic retry logic for failed commands.
- Configurable retry attempts
- Retry delay between attempts
- Can distinguish retryable vs non-retryable errors

#### 5. CommandInvoker

The final interceptor that actually executes the command.

### Interceptor Chain Order

```
LoggingInterceptor (outermost)
    ↓
RetryInterceptor (if enabled)
    ↓
Custom Interceptors
    ↓
TransactionInterceptor
    ↓
ContextInterceptor
    ↓
CommandInvoker (innermost)
```

## Creating Custom Commands

### Example: Custom Query Command

```go
package commands

import (
    "context"
    "github.com/muixstudio/flowgo/engine"
)

type GetActiveProcessInstancesCommand struct {
    ProcessDefinitionKey string
}

func (c *GetActiveProcessInstancesCommand) Execute(
    ctx context.Context,
    commandContext *engine.CommandContext,
) ([]*runtime.ProcessInstance, error) {
    
    runtimeService := commandContext.Engine.GetRuntimeService()
    
    instances, err := runtimeService.CreateProcessInstanceQuery().
        ProcessDefinitionKey(c.ProcessDefinitionKey).
        Active().
        List(ctx)
    
    if err != nil {
        return nil, err
    }
    
    return instances, nil
}
```

### Using the Command

```go
// Create command
cmd := &GetActiveProcessInstancesCommand{
    ProcessDefinitionKey: "expense-approval",
}

// Execute through engine
result, err := engine.Execute(ctx, cmd)
if err != nil {
    log.Fatal(err)
}

instances := result.([]*runtime.ProcessInstance)
for _, instance := range instances {
    fmt.Printf("Instance: %s\n", instance.ID)
}
```

## Creating Custom Interceptors

### Example: Metrics Interceptor

```go
package interceptors

import (
    "context"
    "time"
    "github.com/muixstudio/flowgo/engine"
)

type MetricsInterceptor struct {
    engine.BaseCommandInterceptor
    metrics *MetricsCollector
}

func (i *MetricsInterceptor) Execute(
    ctx context.Context,
    command engine.Command[any],
    next engine.CommandExecutor,
) (any, error) {
    
    commandName := fmt.Sprintf("%T", command)
    start := time.Now()
    
    result, err := i.next.Execute(ctx, command, next)
    
    duration := time.Since(start)
    
    // Record metrics
    i.metrics.RecordCommandExecution(commandName, duration, err == nil)
    
    return result, err
}
```

### Adding Custom Interceptor

```go
metrics := NewMetricsCollector()
metricsInterceptor := &MetricsInterceptor{metrics: metrics}

executor := engine.NewDefaultCommandExecutorBuilder(engineImpl).
    AddInterceptor(metricsInterceptor).
    WithLogging(true).
    WithTransaction(true).
    Build()
```

## Available Commands

### Repository Commands

- `DeployCommand` - Deploy process definitions
- `GetProcessDefinitionCommand` - Get process definition by ID
- `SuspendProcessDefinitionCommand` - Suspend process definition
- `ActivateProcessDefinitionCommand` - Activate process definition

### Runtime Commands

- `StartProcessInstanceCommand` - Start process instance
- `DeleteProcessInstanceCommand` - Delete process instance
- `SetVariablesCommand` - Set process variables
- `SignalExecutionCommand` - Signal execution

### Task Commands

- `ClaimTaskCommand` - Claim a task
- `CompleteTaskCommand` - Complete a task
- `DelegateTaskCommand` - Delegate a task
- `AddCommentCommand` - Add task comment

### History Commands

- `GetHistoricProcessInstancesCommand` - Query historic processes
- `GetHistoricTasksCommand` - Query historic tasks
- `DeleteHistoricDataCommand` - Clean up history

## Benefits of Command Pattern

### 1. Consistency

All operations follow the same execution model, making the codebase consistent and predictable.

### 2. Separation of Concerns

Business logic (in commands) is separated from cross-cutting concerns (in interceptors).

### 3. Transaction Management

Automatic transaction management ensures data consistency.

### 4. Logging and Monitoring

Centralized logging and metrics collection for all operations.

### 5. Error Handling

Consistent error handling and retry logic.

### 6. Testability

Commands can be easily tested in isolation.

```go
func TestDeployCommand(t *testing.T) {
    // Create mock engine and context
    engine := NewMockEngine()
    ctx := engine.NewCommandContext(context.Background())
    
    // Create and execute command
    cmd := NewDeployCommand("test", "test.json", content)
    result, err := cmd.Execute(context.Background(), ctx)
    
    // Assert results
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 7. Extensibility

Easy to add new commands and interceptors without modifying existing code.

### 8. Composition

Commands can be composed to create complex operations.

```go
type ComplexCommand struct {
    subCommands []engine.Command[any]
}

func (c *ComplexCommand) Execute(
    ctx context.Context,
    commandContext *engine.CommandContext,
) (interface{}, error) {
    
    results := make([]interface{}, 0)
    
    for _, subCmd := range c.subCommands {
        result, err := commandContext.Engine.Execute(ctx, subCmd)
        if err != nil {
            return nil, err
        }
        results = append(results, result)
    }
    
    return results, nil
}
```

## Comparison with Flowable/Activiti

| Feature | Flowable/Activiti | FlowGo |
|---------|-------------------|--------|
| Command Interface | `Command<T>` | `Command[T]` (Go generics) |
| Context | `CommandContext` | `CommandContext` |
| Executor | `CommandExecutor` | `CommandExecutor` |
| Interceptors | `CommandInterceptor` | `CommandInterceptor` |
| Transaction | Spring `@Transactional` | `TransactionInterceptor` |
| Logging | SLF4J | Go `log` package |

## Best Practices

### 1. Keep Commands Focused

Each command should do one thing well.

✅ Good:
```go
DeployCommand
StartProcessInstanceCommand
CompleteTaskCommand
```

❌ Bad:
```go
DeployAndStartProcessCommand  // Too much in one command
```

### 2. Use CommandContext

Store command-specific data in `CommandContext.Attributes`.

```go
commandContext.SetAttribute("userId", "john.doe")
commandContext.SetAttribute("timestamp", time.Now())
```

### 3. Handle Errors Properly

Return descriptive errors from commands.

```go
if err != nil {
    return nil, fmt.Errorf("failed to deploy process '%s': %w", name, err)
}
```

### 4. Test Commands Independently

Write unit tests for each command.

### 5. Document Custom Commands

Provide clear documentation for custom commands.

## Performance Considerations

- Interceptors add overhead - only enable what you need
- Use connection pooling for database operations
- Consider async execution for long-running commands
- Cache frequently accessed data in `CommandContext`

## Future Enhancements

- [ ] Command queue for async execution
- [ ] Command history/audit trail
- [ ] Command scheduling
- [ ] Distributed command execution
- [ ] Command compensation/rollback
