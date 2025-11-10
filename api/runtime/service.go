package runtime

import "context"

// Service provides operations for managing process instances and executions.
type Service interface {
	// Initialize initializes the runtime service
	Initialize(ctx context.Context) error

	// Shutdown gracefully shuts down the runtime service
	Shutdown(ctx context.Context) error

	// StartProcessInstanceByKey starts a process instance by process definition key
	StartProcessInstanceByKey(ctx context.Context, processDefinitionKey string, variables map[string]interface{}) (*ProcessInstance, error)

	// StartProcessInstanceByID starts a process instance by process definition ID
	StartProcessInstanceByID(ctx context.Context, processDefinitionID string, variables map[string]interface{}) (*ProcessInstance, error)

	// StartProcessInstanceByKeyWithBusinessKey starts a process instance with a business key
	StartProcessInstanceByKeyWithBusinessKey(ctx context.Context, processDefinitionKey, businessKey string, variables map[string]interface{}) (*ProcessInstance, error)

	// DeleteProcessInstance deletes a process instance
	DeleteProcessInstance(ctx context.Context, processInstanceID, deleteReason string) error

	// SuspendProcessInstance suspends a process instance
	SuspendProcessInstance(ctx context.Context, processInstanceID string) error

	// ActivateProcessInstance activates a suspended process instance
	ActivateProcessInstance(ctx context.Context, processInstanceID string) error

	// CreateProcessInstanceQuery creates a new process instance query
	CreateProcessInstanceQuery() *ProcessInstanceQuery

	// GetProcessInstance retrieves a process instance by ID
	GetProcessInstance(ctx context.Context, processInstanceID string) (*ProcessInstance, error)

	// SetVariable sets a variable on a process instance
	SetVariable(ctx context.Context, executionID, variableName string, value interface{}) error

	// SetVariables sets multiple variables on a process instance
	SetVariables(ctx context.Context, executionID string, variables map[string]interface{}) error

	// GetVariable gets a variable from a process instance
	GetVariable(ctx context.Context, executionID, variableName string) (interface{}, error)

	// GetVariables gets all variables from a process instance
	GetVariables(ctx context.Context, executionID string) (map[string]interface{}, error)

	// RemoveVariable removes a variable from a process instance
	RemoveVariable(ctx context.Context, executionID, variableName string) error

	// Signal triggers a signal event
	Signal(ctx context.Context, executionID string) error

	// SignalWithVariables triggers a signal event with variables
	SignalWithVariables(ctx context.Context, executionID string, variables map[string]interface{}) error

	// CreateExecutionQuery creates a new execution query
	CreateExecutionQuery() *ExecutionQuery
}
