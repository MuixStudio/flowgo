package task

import (
	"context"
	"time"
)

// TaskService provides operations for managing user tasks.
// This service is responsible for:
// - Querying tasks
// - Claiming and unclaiming tasks
// - Completing tasks
// - Assigning and delegating tasks
// - Managing task variables
type TaskService interface {
	// Initialize initializes the task service
	Initialize(ctx context.Context) error

	// Shutdown gracefully shuts down the task service
	Shutdown(ctx context.Context) error

	// CreateTaskQuery creates a new task query
	CreateTaskQuery() *TaskQuery

	// GetTask retrieves a task by ID
	GetTask(ctx context.Context, taskID string) (*Task, error)

	// NewTask creates a new standalone task (not part of a process)
	NewTask(ctx context.Context, taskID string) (*Task, error)

	// SaveTask saves a standalone task
	SaveTask(ctx context.Context, task *Task) error

	// DeleteTask deletes a task
	DeleteTask(ctx context.Context, taskID string) error

	// Claim assigns a task to a specific user
	Claim(ctx context.Context, taskID, userID string) error

	// Unclaim removes the assignee from a task
	Unclaim(ctx context.Context, taskID string) error

	// Complete completes a task
	Complete(ctx context.Context, taskID string) error

	// CompleteWithVariables completes a task and sets variables
	CompleteWithVariables(ctx context.Context, taskID string, variables map[string]interface{}) error

	// SetAssignee sets the assignee of a task
	SetAssignee(ctx context.Context, taskID, userID string) error

	// SetOwner sets the owner of a task
	SetOwner(ctx context.Context, taskID, userID string) error

	// AddCandidateUser adds a candidate user to a task
	AddCandidateUser(ctx context.Context, taskID, userID string) error

	// AddCandidateGroup adds a candidate group to a task
	AddCandidateGroup(ctx context.Context, taskID, groupID string) error

	// DeleteCandidateUser removes a candidate user from a task
	DeleteCandidateUser(ctx context.Context, taskID, userID string) error

	// DeleteCandidateGroup removes a candidate group from a task
	DeleteCandidateGroup(ctx context.Context, taskID, groupID string) error

	// SetPriority sets the priority of a task
	SetPriority(ctx context.Context, taskID string, priority int) error

	// SetDueDate sets the due date of a task
	SetDueDate(ctx context.Context, taskID string, dueDate time.Time) error

	// GetTaskVariables gets all variables of a task
	GetTaskVariables(ctx context.Context, taskID string) (map[string]interface{}, error)

	// GetTaskVariable gets a specific variable of a task
	GetTaskVariable(ctx context.Context, taskID, variableName string) (interface{}, error)

	// SetTaskVariable sets a variable on a task
	SetTaskVariable(ctx context.Context, taskID, variableName string, value interface{}) error

	// SetTaskVariables sets multiple variables on a task
	SetTaskVariables(ctx context.Context, taskID string, variables map[string]interface{}) error

	// RemoveTaskVariable removes a variable from a task
	RemoveTaskVariable(ctx context.Context, taskID, variableName string) error

	// AddComment adds a comment to a task
	AddComment(ctx context.Context, taskID, message string) (*Comment, error)

	// GetTaskComments gets all comments for a task
	GetTaskComments(ctx context.Context, taskID string) ([]*Comment, error)

	// CreateAttachment creates an attachment for a task
	CreateAttachment(ctx context.Context, taskID, attachmentType, attachmentName, attachmentDescription string, content []byte) (*Attachment, error)

	// GetTaskAttachments gets all attachments for a task
	GetTaskAttachments(ctx context.Context, taskID string) ([]*Attachment, error)

	// DeleteAttachment deletes an attachment
	DeleteAttachment(ctx context.Context, attachmentID string) error
}

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
	ParentTaskID        string
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

// Attachment represents an attachment on a task
type Attachment struct {
	ID                string
	Name              string
	Description       string
	Type              string
	TaskID            string
	ProcessInstanceID string
	URL               string
	Content           []byte
	Time              time.Time
}

// TaskQuery provides a fluent API for querying tasks
type TaskQuery struct {
	taskID               string
	taskName             string
	taskDescription      string
	assignee             string
	owner                string
	candidateUser        string
	candidateGroup       string
	processInstanceID    string
	processDefinitionID  string
	processDefinitionKey string
	executionID          string
	taskDefinitionKey    string
	category             string
	tenantID             string
	suspended            *bool
	active               *bool
	priorityMin          *int
	priorityMax          *int
	dueBefore            *time.Time
	dueAfter             *time.Time
	createdBefore        *time.Time
	createdAfter         *time.Time
	variableValueEquals  map[string]interface{}
	orderBy              string
	ascending            bool
	service              TaskService
}

// TaskID filters by task ID
func (q *TaskQuery) TaskID(id string) *TaskQuery {
	q.taskID = id
	return q
}

// TaskName filters by task name
func (q *TaskQuery) TaskName(name string) *TaskQuery {
	q.taskName = name
	return q
}

// TaskDescription filters by task description
func (q *TaskQuery) TaskDescription(description string) *TaskQuery {
	q.taskDescription = description
	return q
}

// TaskAssignee filters by assignee
func (q *TaskQuery) TaskAssignee(assignee string) *TaskQuery {
	q.assignee = assignee
	return q
}

// TaskOwner filters by owner
func (q *TaskQuery) TaskOwner(owner string) *TaskQuery {
	q.owner = owner
	return q
}

// TaskCandidateUser filters by candidate user
func (q *TaskQuery) TaskCandidateUser(userID string) *TaskQuery {
	q.candidateUser = userID
	return q
}

// TaskCandidateGroup filters by candidate group
func (q *TaskQuery) TaskCandidateGroup(groupID string) *TaskQuery {
	q.candidateGroup = groupID
	return q
}

// ProcessInstanceID filters by process instance ID
func (q *TaskQuery) ProcessInstanceID(id string) *TaskQuery {
	q.processInstanceID = id
	return q
}

// ProcessDefinitionID filters by process definition ID
func (q *TaskQuery) ProcessDefinitionID(id string) *TaskQuery {
	q.processDefinitionID = id
	return q
}

// ProcessDefinitionKey filters by process definition key
func (q *TaskQuery) ProcessDefinitionKey(key string) *TaskQuery {
	q.processDefinitionKey = key
	return q
}

// ExecutionID filters by execution ID
func (q *TaskQuery) ExecutionID(id string) *TaskQuery {
	q.executionID = id
	return q
}

// TaskDefinitionKey filters by task definition key
func (q *TaskQuery) TaskDefinitionKey(key string) *TaskQuery {
	q.taskDefinitionKey = key
	return q
}

// TaskCategory filters by category
func (q *TaskQuery) TaskCategory(category string) *TaskQuery {
	q.category = category
	return q
}

// TenantID filters by tenant ID
func (q *TaskQuery) TenantID(tenantID string) *TaskQuery {
	q.tenantID = tenantID
	return q
}

// Active filters to only active tasks
func (q *TaskQuery) Active() *TaskQuery {
	trueVal := true
	q.active = &trueVal
	return q
}

// Suspended filters to only suspended tasks
func (q *TaskQuery) Suspended() *TaskQuery {
	trueVal := true
	q.suspended = &trueVal
	return q
}

// TaskPriority filters by priority
func (q *TaskQuery) TaskPriority(priority int) *TaskQuery {
	q.priorityMin = &priority
	q.priorityMax = &priority
	return q
}

// TaskPriorityMin filters by minimum priority
func (q *TaskQuery) TaskPriorityMin(minPriority int) *TaskQuery {
	q.priorityMin = &minPriority
	return q
}

// TaskPriorityMax filters by maximum priority
func (q *TaskQuery) TaskPriorityMax(maxPriority int) *TaskQuery {
	q.priorityMax = &maxPriority
	return q
}

// DueBefore filters tasks due before a specific date
func (q *TaskQuery) DueBefore(date time.Time) *TaskQuery {
	q.dueBefore = &date
	return q
}

// DueAfter filters tasks due after a specific date
func (q *TaskQuery) DueAfter(date time.Time) *TaskQuery {
	q.dueAfter = &date
	return q
}

// TaskCreatedBefore filters tasks created before a specific date
func (q *TaskQuery) TaskCreatedBefore(date time.Time) *TaskQuery {
	q.createdBefore = &date
	return q
}

// TaskCreatedAfter filters tasks created after a specific date
func (q *TaskQuery) TaskCreatedAfter(date time.Time) *TaskQuery {
	q.createdAfter = &date
	return q
}

// TaskVariableValueEquals filters by variable value
func (q *TaskQuery) TaskVariableValueEquals(name string, value interface{}) *TaskQuery {
	if q.variableValueEquals == nil {
		q.variableValueEquals = make(map[string]interface{})
	}
	q.variableValueEquals[name] = value
	return q
}

// OrderByTaskID orders results by task ID
func (q *TaskQuery) OrderByTaskID() *TaskQuery {
	q.orderBy = "id"
	return q
}

// OrderByTaskName orders results by task name
func (q *TaskQuery) OrderByTaskName() *TaskQuery {
	q.orderBy = "name"
	return q
}

// OrderByTaskPriority orders results by priority
func (q *TaskQuery) OrderByTaskPriority() *TaskQuery {
	q.orderBy = "priority"
	return q
}

// OrderByTaskCreateTime orders results by create time
func (q *TaskQuery) OrderByTaskCreateTime() *TaskQuery {
	q.orderBy = "create_time"
	return q
}

// OrderByDueDate orders results by due date
func (q *TaskQuery) OrderByDueDate() *TaskQuery {
	q.orderBy = "due_date"
	return q
}

// Asc sets ascending order
func (q *TaskQuery) Asc() *TaskQuery {
	q.ascending = true
	return q
}

// Desc sets descending order
func (q *TaskQuery) Desc() *TaskQuery {
	q.ascending = false
	return q
}

// List executes the query and returns a list of tasks
func (q *TaskQuery) List(ctx context.Context) ([]*Task, error) {
	// Will be implemented by the concrete service
	return nil, nil
}

// Count returns the count of matching tasks
func (q *TaskQuery) Count(ctx context.Context) (int64, error) {
	// Will be implemented by the concrete service
	return 0, nil
}

// SingleResult returns a single task or error if not exactly one result
func (q *TaskQuery) SingleResult(ctx context.Context) (*Task, error) {
	// Will be implemented by the concrete service
	return nil, nil
}
