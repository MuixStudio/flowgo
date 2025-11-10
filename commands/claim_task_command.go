package commands

import (
	"context"
	"fmt"

	"github.com/muixstudio/flowgo/engine"
)

// ClaimTaskCommand claims a task for a user
type ClaimTaskCommand struct {
	TaskID string
	UserID string
}

// Execute claims the task
func (c *ClaimTaskCommand) Execute(ctx context.Context, commandContext *engine.CommandContext) (interface{}, error) {
	if c.TaskID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	if c.UserID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	taskService := commandContext.Engine.GetTaskService()

	// Get the task to verify it exists
	task, err := taskService.GetTask(ctx, c.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Check if task is already claimed by another user
	if task.Assignee != "" && task.Assignee != c.UserID {
		return nil, fmt.Errorf("task '%s' is already claimed by user '%s'", c.TaskID, task.Assignee)
	}

	// Check if user is a candidate for this task
	if len(task.CandidateUsers) > 0 {
		isCandidate := false
		for _, candidateUser := range task.CandidateUsers {
			if candidateUser == c.UserID {
				isCandidate = true
				break
			}
		}
		if !isCandidate {
			return nil, fmt.Errorf("user '%s' is not a candidate for task '%s'", c.UserID, c.TaskID)
		}
	}

	// Claim the task
	if err := taskService.Claim(ctx, c.TaskID, c.UserID); err != nil {
		return nil, fmt.Errorf("failed to claim task: %w", err)
	}

	// Record to history if enabled
	if commandContext.Engine.GetConfiguration().EnableHistory {
		// TODO: Record task claim event
	}

	return nil, nil
}

// NewClaimTaskCommand creates a new claim task command
func NewClaimTaskCommand(taskID, userID string) *ClaimTaskCommand {
	return &ClaimTaskCommand{
		TaskID: taskID,
		UserID: userID,
	}
}
