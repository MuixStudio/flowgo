package history

import (
	"context"
	"time"
)

// HistoricProcessInstance represents a completed or running process instance in history
type HistoricProcessInstance struct {
	ID                       string
	BusinessKey              string
	ProcessDefinitionID      string
	ProcessDefinitionKey     string
	ProcessDefinitionName    string
	ProcessDefinitionVersion int
	StartTime                time.Time
	EndTime                  *time.Time
	DurationInMillis         *int64
	StartUserID              string
	DeleteReason             string
	TenantID                 string
}

// HistoricTaskInstance represents a completed or running task in history
type HistoricTaskInstance struct {
	ID                   string
	ProcessDefinitionID  string
	ProcessDefinitionKey string
	ProcessInstanceID    string
	Name                 string
	Assignee             string
	StartTime            time.Time
	EndTime              *time.Time
	DurationInMillis     *int64
	Priority             int
	TenantID             string
}

// HistoricProcessInstanceQuery provides a fluent API for querying historic process instances
type HistoricProcessInstanceQuery struct {
	processInstanceID     string
	processDefinitionKey  string
	finished              *bool
	startedAfter          *time.Time
	orderBy               string
	ascending             bool
	service               Service
}

// ProcessInstanceID filters by process instance ID
func (q *HistoricProcessInstanceQuery) ProcessInstanceID(id string) *HistoricProcessInstanceQuery {
	q.processInstanceID = id
	return q
}

// ProcessDefinitionKey filters by process definition key
func (q *HistoricProcessInstanceQuery) ProcessDefinitionKey(key string) *HistoricProcessInstanceQuery {
	q.processDefinitionKey = key
	return q
}

// Finished filters to only finished process instances
func (q *HistoricProcessInstanceQuery) Finished() *HistoricProcessInstanceQuery {
	trueVal := true
	q.finished = &trueVal
	return q
}

// OrderByStartTime orders results by start time
func (q *HistoricProcessInstanceQuery) OrderByStartTime() *HistoricProcessInstanceQuery {
	q.orderBy = "start_time"
	return q
}

// Desc sets descending order
func (q *HistoricProcessInstanceQuery) Desc() *HistoricProcessInstanceQuery {
	q.ascending = false
	return q
}

// List executes the query and returns a list of historic process instances
func (q *HistoricProcessInstanceQuery) List(ctx context.Context) ([]*HistoricProcessInstance, error) {
	// TODO: Implement
	return nil, nil
}

// HistoricTaskInstanceQuery provides a fluent API for querying historic task instances
type HistoricTaskInstanceQuery struct {
	taskID              string
	processInstanceID   string
	assignee            string
	finished            *bool
	orderBy             string
	ascending           bool
	service             Service
}

// TaskID filters by task ID
func (q *HistoricTaskInstanceQuery) TaskID(id string) *HistoricTaskInstanceQuery {
	q.taskID = id
	return q
}

// List executes the query and returns a list of historic task instances
func (q *HistoricTaskInstanceQuery) List(ctx context.Context) ([]*HistoricTaskInstance, error) {
	// TODO: Implement
	return nil, nil
}
