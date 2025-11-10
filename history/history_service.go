package history

import (
	"context"
	"time"
)

// HistoryService provides operations for querying historical process data.
// This service is responsible for:
// - Querying historical process instances
// - Querying historical tasks
// - Querying historical variables
// - Querying historical activities
// - Deleting historical data
type HistoryService interface {
	// Initialize initializes the history service
	Initialize(ctx context.Context) error

	// Shutdown gracefully shuts down the history service
	Shutdown(ctx context.Context) error

	// CreateHistoricProcessInstanceQuery creates a new historic process instance query
	CreateHistoricProcessInstanceQuery() *HistoricProcessInstanceQuery

	// CreateHistoricTaskInstanceQuery creates a new historic task instance query
	CreateHistoricTaskInstanceQuery() *HistoricTaskInstanceQuery

	// CreateHistoricActivityInstanceQuery creates a new historic activity instance query
	CreateHistoricActivityInstanceQuery() *HistoricActivityInstanceQuery

	// CreateHistoricVariableInstanceQuery creates a new historic variable instance query
	CreateHistoricVariableInstanceQuery() *HistoricVariableInstanceQuery

	// DeleteHistoricProcessInstance deletes a historic process instance
	DeleteHistoricProcessInstance(ctx context.Context, processInstanceID string) error

	// DeleteHistoricTaskInstance deletes a historic task instance
	DeleteHistoricTaskInstance(ctx context.Context, taskID string) error

	// RecordProcessInstance records a process instance to history
	RecordProcessInstance(ctx context.Context, instance *HistoricProcessInstance) error

	// RecordTaskInstance records a task instance to history
	RecordTaskInstance(ctx context.Context, task *HistoricTaskInstance) error

	// RecordActivityInstance records an activity instance to history
	RecordActivityInstance(ctx context.Context, activity *HistoricActivityInstance) error

	// RecordVariableInstance records a variable instance to history
	RecordVariableInstance(ctx context.Context, variable *HistoricVariableInstance) error
}

// HistoricProcessInstance represents a completed or running process instance in history
type HistoricProcessInstance struct {
	ID                   string
	BusinessKey          string
	ProcessDefinitionID  string
	ProcessDefinitionKey string
	ProcessDefinitionName string
	ProcessDefinitionVersion int
	DeploymentID         string
	StartTime            time.Time
	EndTime              *time.Time
	DurationInMillis     *int64
	StartUserID          string
	StartActivityID      string
	EndActivityID        string
	DeleteReason         string
	SuperProcessInstanceID string
	TenantID             string
}

// HistoricTaskInstance represents a completed or running task in history
type HistoricTaskInstance struct {
	ID                  string
	ProcessDefinitionID string
	ProcessDefinitionKey string
	ProcessInstanceID   string
	ExecutionID         string
	Name                string
	Description         string
	TaskDefinitionKey   string
	Owner               string
	Assignee            string
	StartTime           time.Time
	EndTime             *time.Time
	DurationInMillis    *int64
	DeleteReason        string
	Priority            int
	DueDate             *time.Time
	FormKey             string
	Category            string
	TenantID            string
}

// HistoricActivityInstance represents a completed or running activity in history
type HistoricActivityInstance struct {
	ID                  string
	ActivityID          string
	ActivityName        string
	ActivityType        string
	ProcessDefinitionID string
	ProcessInstanceID   string
	ExecutionID         string
	TaskID              string
	Assignee            string
	StartTime           time.Time
	EndTime             *time.Time
	DurationInMillis    *int64
	DeleteReason        string
	TenantID            string
}

// HistoricVariableInstance represents a variable value at a point in history
type HistoricVariableInstance struct {
	ID                  string
	Name                string
	TypeName            string
	Value               interface{}
	ProcessInstanceID   string
	TaskID              string
	CreateTime          time.Time
	LastUpdatedTime     *time.Time
}

// HistoricProcessInstanceQuery provides a fluent API for querying historic process instances
type HistoricProcessInstanceQuery struct {
	processInstanceID        string
	processInstanceBusinessKey string
	processDefinitionID      string
	processDefinitionKey     string
	processDefinitionName    string
	deploymentID             string
	startUserID              string
	superProcessInstanceID   string
	tenantID                 string
	finished                 *bool
	unfinished               *bool
	startedBefore            *time.Time
	startedAfter             *time.Time
	finishedBefore           *time.Time
	finishedAfter            *time.Time
	variableValueEquals      map[string]interface{}
	orderBy                  string
	ascending                bool
	service                  HistoryService
}

// ProcessInstanceID filters by process instance ID
func (q *HistoricProcessInstanceQuery) ProcessInstanceID(id string) *HistoricProcessInstanceQuery {
	q.processInstanceID = id
	return q
}

// ProcessInstanceBusinessKey filters by business key
func (q *HistoricProcessInstanceQuery) ProcessInstanceBusinessKey(businessKey string) *HistoricProcessInstanceQuery {
	q.processInstanceBusinessKey = businessKey
	return q
}

// ProcessDefinitionID filters by process definition ID
func (q *HistoricProcessInstanceQuery) ProcessDefinitionID(id string) *HistoricProcessInstanceQuery {
	q.processDefinitionID = id
	return q
}

// ProcessDefinitionKey filters by process definition key
func (q *HistoricProcessInstanceQuery) ProcessDefinitionKey(key string) *HistoricProcessInstanceQuery {
	q.processDefinitionKey = key
	return q
}

// StartUserID filters by the user who started the process
func (q *HistoricProcessInstanceQuery) StartUserID(userID string) *HistoricProcessInstanceQuery {
	q.startUserID = userID
	return q
}

// Finished filters to only finished process instances
func (q *HistoricProcessInstanceQuery) Finished() *HistoricProcessInstanceQuery {
	trueVal := true
	q.finished = &trueVal
	return q
}

// Unfinished filters to only unfinished process instances
func (q *HistoricProcessInstanceQuery) Unfinished() *HistoricProcessInstanceQuery {
	trueVal := true
	q.unfinished = &trueVal
	return q
}

// StartedBefore filters to process instances started before a specific date
func (q *HistoricProcessInstanceQuery) StartedBefore(date time.Time) *HistoricProcessInstanceQuery {
	q.startedBefore = &date
	return q
}

// StartedAfter filters to process instances started after a specific date
func (q *HistoricProcessInstanceQuery) StartedAfter(date time.Time) *HistoricProcessInstanceQuery {
	q.startedAfter = &date
	return q
}

// FinishedBefore filters to process instances finished before a specific date
func (q *HistoricProcessInstanceQuery) FinishedBefore(date time.Time) *HistoricProcessInstanceQuery {
	q.finishedBefore = &date
	return q
}

// FinishedAfter filters to process instances finished after a specific date
func (q *HistoricProcessInstanceQuery) FinishedAfter(date time.Time) *HistoricProcessInstanceQuery {
	q.finishedAfter = &date
	return q
}

// OrderByProcessInstanceID orders results by process instance ID
func (q *HistoricProcessInstanceQuery) OrderByProcessInstanceID() *HistoricProcessInstanceQuery {
	q.orderBy = "id"
	return q
}

// OrderByStartTime orders results by start time
func (q *HistoricProcessInstanceQuery) OrderByStartTime() *HistoricProcessInstanceQuery {
	q.orderBy = "start_time"
	return q
}

// OrderByEndTime orders results by end time
func (q *HistoricProcessInstanceQuery) OrderByEndTime() *HistoricProcessInstanceQuery {
	q.orderBy = "end_time"
	return q
}

// OrderByDuration orders results by duration
func (q *HistoricProcessInstanceQuery) OrderByDuration() *HistoricProcessInstanceQuery {
	q.orderBy = "duration"
	return q
}

// Asc sets ascending order
func (q *HistoricProcessInstanceQuery) Asc() *HistoricProcessInstanceQuery {
	q.ascending = true
	return q
}

// Desc sets descending order
func (q *HistoricProcessInstanceQuery) Desc() *HistoricProcessInstanceQuery {
	q.ascending = false
	return q
}

// List executes the query and returns a list of historic process instances
func (q *HistoricProcessInstanceQuery) List(ctx context.Context) ([]*HistoricProcessInstance, error) {
	// Will be implemented by the concrete service
	return nil, nil
}

// Count returns the count of matching historic process instances
func (q *HistoricProcessInstanceQuery) Count(ctx context.Context) (int64, error) {
	// Will be implemented by the concrete service
	return 0, nil
}

// HistoricTaskInstanceQuery provides a fluent API for querying historic task instances
type HistoricTaskInstanceQuery struct {
	taskID                string
	processInstanceID     string
	processDefinitionID   string
	processDefinitionKey  string
	executionID           string
	taskDefinitionKey     string
	assignee              string
	owner                 string
	taskName              string
	tenantID              string
	finished              *bool
	unfinished            *bool
	variableValueEquals   map[string]interface{}
	orderBy               string
	ascending             bool
	service               HistoryService
}

// TaskID filters by task ID
func (q *HistoricTaskInstanceQuery) TaskID(id string) *HistoricTaskInstanceQuery {
	q.taskID = id
	return q
}

// ProcessInstanceID filters by process instance ID
func (q *HistoricTaskInstanceQuery) ProcessInstanceID(id string) *HistoricTaskInstanceQuery {
	q.processInstanceID = id
	return q
}

// TaskAssignee filters by assignee
func (q *HistoricTaskInstanceQuery) TaskAssignee(assignee string) *HistoricTaskInstanceQuery {
	q.assignee = assignee
	return q
}

// TaskOwner filters by owner
func (q *HistoricTaskInstanceQuery) TaskOwner(owner string) *HistoricTaskInstanceQuery {
	q.owner = owner
	return q
}

// Finished filters to only finished tasks
func (q *HistoricTaskInstanceQuery) Finished() *HistoricTaskInstanceQuery {
	trueVal := true
	q.finished = &trueVal
	return q
}

// Unfinished filters to only unfinished tasks
func (q *HistoricTaskInstanceQuery) Unfinished() *HistoricTaskInstanceQuery {
	trueVal := true
	q.unfinished = &trueVal
	return q
}

// List executes the query and returns a list of historic task instances
func (q *HistoricTaskInstanceQuery) List(ctx context.Context) ([]*HistoricTaskInstance, error) {
	// Will be implemented by the concrete service
	return nil, nil
}

// Count returns the count of matching historic task instances
func (q *HistoricTaskInstanceQuery) Count(ctx context.Context) (int64, error) {
	// Will be implemented by the concrete service
	return 0, nil
}

// HistoricActivityInstanceQuery provides a fluent API for querying historic activity instances
type HistoricActivityInstanceQuery struct {
	activityID          string
	activityType        string
	processInstanceID   string
	processDefinitionID string
	executionID         string
	finished            *bool
	orderBy             string
	ascending           bool
	service             HistoryService
}

// ActivityID filters by activity ID
func (q *HistoricActivityInstanceQuery) ActivityID(id string) *HistoricActivityInstanceQuery {
	q.activityID = id
	return q
}

// ActivityType filters by activity type
func (q *HistoricActivityInstanceQuery) ActivityType(activityType string) *HistoricActivityInstanceQuery {
	q.activityType = activityType
	return q
}

// ProcessInstanceID filters by process instance ID
func (q *HistoricActivityInstanceQuery) ProcessInstanceID(id string) *HistoricActivityInstanceQuery {
	q.processInstanceID = id
	return q
}

// Finished filters to only finished activities
func (q *HistoricActivityInstanceQuery) Finished() *HistoricActivityInstanceQuery {
	trueVal := true
	q.finished = &trueVal
	return q
}

// List executes the query and returns a list of historic activity instances
func (q *HistoricActivityInstanceQuery) List(ctx context.Context) ([]*HistoricActivityInstance, error) {
	// Will be implemented by the concrete service
	return nil, nil
}

// Count returns the count of matching historic activity instances
func (q *HistoricActivityInstanceQuery) Count(ctx context.Context) (int64, error) {
	// Will be implemented by the concrete service
	return 0, nil
}

// HistoricVariableInstanceQuery provides a fluent API for querying historic variable instances
type HistoricVariableInstanceQuery struct {
	variableName        string
	processInstanceID   string
	taskID              string
	orderBy             string
	ascending           bool
	service             HistoryService
}

// VariableName filters by variable name
func (q *HistoricVariableInstanceQuery) VariableName(name string) *HistoricVariableInstanceQuery {
	q.variableName = name
	return q
}

// ProcessInstanceID filters by process instance ID
func (q *HistoricVariableInstanceQuery) ProcessInstanceID(id string) *HistoricVariableInstanceQuery {
	q.processInstanceID = id
	return q
}

// TaskID filters by task ID
func (q *HistoricVariableInstanceQuery) TaskID(id string) *HistoricVariableInstanceQuery {
	q.taskID = id
	return q
}

// List executes the query and returns a list of historic variable instances
func (q *HistoricVariableInstanceQuery) List(ctx context.Context) ([]*HistoricVariableInstance, error) {
	// Will be implemented by the concrete service
	return nil, nil
}

// Count returns the count of matching historic variable instances
func (q *HistoricVariableInstanceQuery) Count(ctx context.Context) (int64, error) {
	// Will be implemented by the concrete service
	return 0, nil
}
