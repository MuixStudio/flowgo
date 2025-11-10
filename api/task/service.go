package task

import (
	"context"
	"time"
)

// Service provides operations for managing user tasks.
type Service interface {
	// Initialize initializes the task service
	Initialize(ctx context.Context) error

	// Shutdown gracefully shuts down the task service
	Shutdown(ctx context.Context) error

	// CreateTaskQuery creates a new task query
	CreateTaskQuery() *TaskQuery

	// GetTask retrieves a task by ID
	GetTask(ctx context.Context, taskID string) (*Task, error)

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

	// AddCandidateUser adds a candidate user to a task
	AddCandidateUser(ctx context.Context, taskID, userID string) error

	// AddCandidateGroup adds a candidate group to a task
	AddCandidateGroup(ctx context.Context, taskID, groupID string) error

	// SetPriority sets the priority of a task
	SetPriority(ctx context.Context, taskID string, priority int) error

	// AddComment adds a comment to a task
	AddComment(ctx context.Context, taskID, message string) (*Comment, error)

	// GetTaskComments gets all comments for a task
	GetTaskComments(ctx context.Context, taskID string) ([]*Comment, error)
}
