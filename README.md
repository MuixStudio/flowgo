# FlowGo - Go Workflow Engine

FlowGo is a lightweight, high-performance workflow engine for Go, inspired by Flowable/Activiti but designed specifically for Go applications. It uses JSON-based process definitions instead of BPMN XML.

## Features

- 🚀 **High Performance**: Built with Go's concurrency in mind
- 📝 **JSON-based Process Definitions**: Easy to read and write, no XML
- 🔄 **Complete Process Lifecycle**: Deploy, start, suspend, activate, delete
- 👥 **User Task Management**: Assignment, claiming, delegation
- 📊 **History Tracking**: Complete audit trail of process execution
- 🔌 **Extensible**: Easy to add custom service tasks and expressions
- 🎯 **Event-Driven**: Support for timers, messages, signals
- 🌐 **Multi-Tenancy**: Built-in tenant support

## Architecture

FlowGo follows a service-oriented architecture similar to Activiti/Flowable:

```
┌─────────────────────────────────────┐
│       ProcessEngine (Facade)       │
└────────────┬────────────────────────┘
             │
   ┌─────────┴─────────┐
   │                   │
   ▼                   ▼
┌──────────────┐  ┌──────────────┐
│ Repository   │  │   Runtime    │
│   Service    │  │   Service    │
└──────────────┘  └──────────────┘
   │                   │
   ▼                   ▼
┌──────────────┐  ┌──────────────┐
│    Task      │  │   History    │
│   Service    │  │   Service    │
└──────────────┘  └──────────────┘
```

## Core Services

### ProcessEngine

The main entry point and facade for all services.

```go
engine, err := flowgo.NewProcessEngineBuilder().
    WithEngineName("my-engine").
    WithDatabase("postgres", "postgresql://localhost:5432/flowgo").
    WithHistory(true).
    WithAsync(true).
    Build()

if err != nil {
    log.Fatal(err)
}

ctx := context.Background()
engine.Start(ctx)
defer engine.Stop(ctx)
```

### RepositoryService

Manages process definitions and deployments.

```go
repoService := engine.GetRepositoryService()

// Deploy a process definition
deployment, err := repoService.CreateDeployment().
    Name("My Process").
    AddProcessDefinition("process.json", jsonContent).
    Deploy(ctx)

// Query process definitions
definitions, err := repoService.CreateProcessDefinitionQuery().
    ProcessDefinitionKey("my-process").
    LatestVersion().
    List(ctx)

// Suspend a process definition
err = repoService.SuspendProcessDefinition(ctx, definitionID)
```

### RuntimeService

Manages process instances and executions.

```go
runtimeService := engine.GetRuntimeService()

// Start a process instance
variables := map[string]interface{}{
    "amount": 1000,
    "requester": "john.doe",
}

instance, err := runtimeService.StartProcessInstanceByKey(
    ctx, "expense-approval", variables)

// Query process instances
instances, err := runtimeService.CreateProcessInstanceQuery().
    ProcessDefinitionKey("expense-approval").
    Active().
    List(ctx)

// Set variables
err = runtimeService.SetVariable(ctx, instance.ID, "approved", true)

// Get variables
vars, err := runtimeService.GetVariables(ctx, instance.ID)
```

### TaskService

Manages user tasks.

```go
taskService := engine.GetTaskService()

// Query tasks
tasks, err := taskService.CreateTaskQuery().
    TaskCandidateUser("john.doe").
    Active().
    OrderByTaskPriority().Desc().
    List(ctx)

// Claim a task
err = taskService.Claim(ctx, taskID, "john.doe")

// Add a comment
comment, err := taskService.AddComment(ctx, taskID, "Reviewed and approved")

// Complete a task
variables := map[string]interface{}{
    "approved": true,
    "comment": "Looks good",
}
err = taskService.CompleteWithVariables(ctx, taskID, variables)
```

### HistoryService

Queries historical process data.

```go
historyService := engine.GetHistoryService()

// Query completed process instances
historicProcesses, err := historyService.CreateHistoricProcessInstanceQuery().
    Finished().
    ProcessDefinitionKey("expense-approval").
    StartedAfter(time.Now().AddDate(0, -1, 0)). // Last month
    OrderByStartTime().Desc().
    List(ctx)

// Query historical tasks
historicTasks, err := historyService.CreateHistoricTaskInstanceQuery().
    Finished().
    TaskAssignee("john.doe").
    List(ctx)
```

## Process Definition Format

FlowGo uses JSON instead of BPMN XML. Here's a simple example:

```json
{
  "id": "expense-approval",
  "name": "Expense Approval Process",
  "version": 1,
  "nodes": [
    {
      "id": "start",
      "type": "startEvent",
      "name": "Start"
    },
    {
      "id": "submit-expense",
      "type": "userTask",
      "name": "Submit Expense Report",
      "properties": {
        "assignee": "${requester}",
        "formKey": "expense-form"
      }
    },
    {
      "id": "approve-expense",
      "type": "userTask",
      "name": "Approve Expense",
      "properties": {
        "candidateGroups": ["managers"],
        "priority": 8
      }
    },
    {
      "id": "end",
      "type": "endEvent",
      "name": "End"
    }
  ],
  "edges": [
    {
      "id": "flow1",
      "source": "start",
      "target": "submit-expense"
    },
    {
      "id": "flow2",
      "source": "submit-expense",
      "target": "approve-expense"
    },
    {
      "id": "flow3",
      "source": "approve-expense",
      "target": "end"
    }
  ]
}
```

See `schema/README.md` for complete documentation on the JSON format.

## Node Types

### Events
- **startEvent**: Process start
- **endEvent**: Process end
- **intermediateEvent**: Timer, message, signal events
- **boundaryEvent**: Events attached to activities

### Tasks
- **userTask**: Manual task requiring human interaction
- **serviceTask**: Automated task executing business logic
- **scriptTask**: Execute script code
- **callActivity**: Call another process
- **subProcess**: Embedded subprocess

### Gateways
- **exclusiveGateway**: XOR - choose one path
- **parallelGateway**: AND - execute all paths
- **inclusiveGateway**: OR - execute multiple paths based on conditions
- **eventBasedGateway**: Wait for events

## Expression Language

FlowGo supports expressions for dynamic behavior:

```json
{
  "assignee": "${processInitiator}",
  "dueDate": "${now() + duration('P2D')}",
  "condition": "${amount > 1000}"
}
```

Supported expression types:
- Variable references: `${variableName}`
- Comparisons: `${amount > 1000}`
- Logical operations: `${condition1 && condition2}`
- Function calls: `${now()}`, `${duration('P2D')}`

## Installation

```bash
go get github.com/muixstudio/flowgo
```

## Quick Start

1. Create a process definition JSON file
2. Initialize the process engine
3. Deploy the process definition
4. Start process instances
5. Complete tasks

See `examples/basic_usage.go` for a complete example.

## Project Structure

```
flowgo/
├── engine.go                 # ProcessEngine interface
├── engine_impl.go            # ProcessEngine implementation
├── repository/               # Repository service
│   ├── repository_service.go
│   └── repository_service_impl.go
├── runtime/                  # Runtime service
│   ├── runtime_service.go
│   └── runtime_service_impl.go
├── task/                     # Task service
│   ├── task_service.go
│   └── task_service_impl.go
├── history/                  # History service
│   ├── history_service.go
│   └── history_service_impl.go
├── schema/                   # JSON Schema definitions
│   ├── process_definition.schema.json
│   └── README.md
└── examples/                 # Example processes and code
    ├── leave_approval.json
    └── basic_usage.go
```

## Roadmap

- [ ] Complete process execution engine
- [ ] Expression language evaluation
- [ ] Database persistence layer
- [ ] Async job executor
- [ ] REST API
- [ ] Process designer UI
- [ ] More event types (message, signal, timer)
- [ ] Compensation and error handling
- [ ] Process migration
- [ ] Multi-instance activities

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## Acknowledgments

Inspired by:
- [Flowable](https://flowable.org/)
- [Activiti](https://www.activiti.org/)
- [Camunda](https://camunda.com/)
