package task

import (
	"context"
	"time"
)

// Task represents a user task in a process
type Task struct {
	ID                  string
	Name                string
	Description         string
	Priority            int
	Owner               string
	Assignee            string
	DueDate             *time.Time
	Category            string
	FormKey             string
	ProcessInstanceID   string
	ProcessDefinitionID string
	ExecutionID         string
	TaskDefinitionKey   string
	CreateTime          time.Time
	ClaimTime           *time.Time
	TenantID            string
	Suspended           bool
	CandidateUsers      []string
	CandidateGroups     []string
}

// Comment represents a comment on a task
type Comment struct {
	ID      string
	TaskID  string
	UserID  string
	Message string
	Time    time.Time
}

// TaskQuery provides a fluent API for querying tasks
type TaskQuery struct {
	taskID                string
	taskName              string
	assignee              string
	owner                 string
	candidateUser         string
	candidateGroup        string
	processInstanceID     string
	processDefinitionKey  string
	suspended             *bool
	active                *bool
	orderBy               string
	ascending             bool
	service               Service
}

// TaskID filters by task ID
func (q *TaskQuery) TaskID(id string) *TaskQuery {
	q.taskID = id
	return q
}

// TaskAssignee filters by assignee
func (q *TaskQuery) TaskAssignee(assignee string) *TaskQuery {
	q.assignee = assignee
	return q
}

// TaskCandidateUser filters by candidate user
func (q *TaskQuery) TaskCandidateUser(userID string) *TaskQuery {
	q.candidateUser = userID
	return q
}

// ProcessInstanceID filters by process instance ID
func (q *TaskQuery) ProcessInstanceID(id string) *TaskQuery {
	q.processInstanceID = id
	return q
}

// Active filters to only active tasks
func (q *TaskQuery) Active() *TaskQuery {
	trueVal := true
	q.active = &trueVal
	return q
}

// OrderByTaskCreateTime orders results by create time
func (q *TaskQuery) OrderByTaskCreateTime() *TaskQuery {
	q.orderBy = "create_time"
	return q
}

// Asc sets ascending order
func (q *TaskQuery) Asc() *TaskQuery {
	q.ascending = true
	return q
}

// List executes the query and returns a list of tasks
func (q *TaskQuery) List(ctx context.Context) ([]*Task, error) {
	// TODO: Implement
	return nil, nil
}
