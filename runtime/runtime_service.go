package runtime

import (
	"context"
	"time"

	"github.com/muixstudio/flowgo/repository"
)

// RuntimeService provides operations for managing process instances and executions.
// This service is responsible for:
// - Starting process instances
// - Managing process instance lifecycle
// - Setting and retrieving process variables
// - Signaling events
// - Managing executions
type RuntimeService interface {
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

// ProcessInstance represents a running or completed process instance
type ProcessInstance struct {
	ID                   string
	ProcessDefinitionID  string
	ProcessDefinitionKey string
	ProcessDefinitionName string
	BusinessKey          string
	StartTime            time.Time
	EndTime              *time.Time
	StartUserID          string
	Suspended            bool
	TenantID             string
	RootProcessInstanceID string
	ParentProcessInstanceID string
}

// Execution represents an execution (thread of control) within a process instance
type Execution struct {
	ID                  string
	ProcessInstanceID   string
	ParentID            string
	ActivityID          string
	IsActive            bool
	IsConcurrent        bool
	IsScope             bool
	IsEventScope        bool
	Suspended           bool
	TenantID            string
}

// ProcessInstanceQuery provides a fluent API for querying process instances
type ProcessInstanceQuery struct {
	processInstanceID        string
	processInstanceBusinessKey string
	processDefinitionID      string
	processDefinitionKey     string
	processDefinitionName    string
	superProcessInstanceID   string
	subProcessInstanceID     string
	startUserID              string
	tenantID                 string
	suspended                *bool
	active                   *bool
	variableValueEquals      map[string]interface{}
	orderBy                  string
	ascending                bool
	service                  RuntimeService
}

// ProcessInstanceID filters by process instance ID
func (q *ProcessInstanceQuery) ProcessInstanceID(id string) *ProcessInstanceQuery {
	q.processInstanceID = id
	return q
}

// ProcessInstanceBusinessKey filters by business key
func (q *ProcessInstanceQuery) ProcessInstanceBusinessKey(businessKey string) *ProcessInstanceQuery {
	q.processInstanceBusinessKey = businessKey
	return q
}

// ProcessDefinitionID filters by process definition ID
func (q *ProcessInstanceQuery) ProcessDefinitionID(id string) *ProcessInstanceQuery {
	q.processDefinitionID = id
	return q
}

// ProcessDefinitionKey filters by process definition key
func (q *ProcessInstanceQuery) ProcessDefinitionKey(key string) *ProcessInstanceQuery {
	q.processDefinitionKey = key
	return q
}

// StartUserID filters by the user who started the process
func (q *ProcessInstanceQuery) StartUserID(userID string) *ProcessInstanceQuery {
	q.startUserID = userID
	return q
}

// SuperProcessInstanceID filters by super process instance ID
func (q *ProcessInstanceQuery) SuperProcessInstanceID(id string) *ProcessInstanceQuery {
	q.superProcessInstanceID = id
	return q
}

// SubProcessInstanceID filters by sub process instance ID
func (q *ProcessInstanceQuery) SubProcessInstanceID(id string) *ProcessInstanceQuery {
	q.subProcessInstanceID = id
	return q
}

// TenantID filters by tenant ID
func (q *ProcessInstanceQuery) TenantID(tenantID string) *ProcessInstanceQuery {
	q.tenantID = tenantID
	return q
}

// Active filters to only active process instances
func (q *ProcessInstanceQuery) Active() *ProcessInstanceQuery {
	trueVal := true
	q.active = &trueVal
	return q
}

// Suspended filters to only suspended process instances
func (q *ProcessInstanceQuery) Suspended() *ProcessInstanceQuery {
	trueVal := true
	q.suspended = &trueVal
	return q
}

// VariableValueEquals filters by variable value
func (q *ProcessInstanceQuery) VariableValueEquals(name string, value interface{}) *ProcessInstanceQuery {
	if q.variableValueEquals == nil {
		q.variableValueEquals = make(map[string]interface{})
	}
	q.variableValueEquals[name] = value
	return q
}

// OrderByProcessInstanceID orders results by process instance ID
func (q *ProcessInstanceQuery) OrderByProcessInstanceID() *ProcessInstanceQuery {
	q.orderBy = "id"
	return q
}

// OrderByProcessDefinitionKey orders results by process definition key
func (q *ProcessInstanceQuery) OrderByProcessDefinitionKey() *ProcessInstanceQuery {
	q.orderBy = "process_definition_key"
	return q
}

// OrderByStartTime orders results by start time
func (q *ProcessInstanceQuery) OrderByStartTime() *ProcessInstanceQuery {
	q.orderBy = "start_time"
	return q
}

// Asc sets ascending order
func (q *ProcessInstanceQuery) Asc() *ProcessInstanceQuery {
	q.ascending = true
	return q
}

// Desc sets descending order
func (q *ProcessInstanceQuery) Desc() *ProcessInstanceQuery {
	q.ascending = false
	return q
}

// List executes the query and returns a list of process instances
func (q *ProcessInstanceQuery) List(ctx context.Context) ([]*ProcessInstance, error) {
	// Will be implemented by the concrete service
	return nil, nil
}

// Count returns the count of matching process instances
func (q *ProcessInstanceQuery) Count(ctx context.Context) (int64, error) {
	// Will be implemented by the concrete service
	return 0, nil
}

// SingleResult returns a single process instance or error if not exactly one result
func (q *ProcessInstanceQuery) SingleResult(ctx context.Context) (*ProcessInstance, error) {
	// Will be implemented by the concrete service
	return nil, nil
}

// ExecutionQuery provides a fluent API for querying executions
type ExecutionQuery struct {
	executionID         string
	processInstanceID   string
	processDefinitionID string
	processDefinitionKey string
	activityID          string
	parentID            string
	tenantID            string
	active              *bool
	orderBy             string
	ascending           bool
	service             RuntimeService
}

// ExecutionID filters by execution ID
func (q *ExecutionQuery) ExecutionID(id string) *ExecutionQuery {
	q.executionID = id
	return q
}

// ProcessInstanceID filters by process instance ID
func (q *ExecutionQuery) ProcessInstanceID(id string) *ExecutionQuery {
	q.processInstanceID = id
	return q
}

// ProcessDefinitionID filters by process definition ID
func (q *ExecutionQuery) ProcessDefinitionID(id string) *ExecutionQuery {
	q.processDefinitionID = id
	return q
}

// ActivityID filters by activity ID
func (q *ExecutionQuery) ActivityID(activityID string) *ExecutionQuery {
	q.activityID = activityID
	return q
}

// ParentID filters by parent execution ID
func (q *ExecutionQuery) ParentID(parentID string) *ExecutionQuery {
	q.parentID = parentID
	return q
}

// Active filters to only active executions
func (q *ExecutionQuery) Active() *ExecutionQuery {
	trueVal := true
	q.active = &trueVal
	return q
}

// List executes the query and returns a list of executions
func (q *ExecutionQuery) List(ctx context.Context) ([]*Execution, error) {
	// Will be implemented by the concrete service
	return nil, nil
}

// Count returns the count of matching executions
func (q *ExecutionQuery) Count(ctx context.Context) (int64, error) {
	// Will be implemented by the concrete service
	return 0, nil
}
