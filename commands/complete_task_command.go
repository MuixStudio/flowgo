package commands

import (
	"context"
	"fmt"

	"github.com/muixstudio/flowgo/engine"
)

// CompleteTaskCommand completes a user task
type CompleteTaskCommand struct {
	TaskID    string
	Variables map[string]interface{}
}

// Execute completes the task
func (c *CompleteTaskCommand) Execute(ctx context.Context, commandContext *engine.CommandContext) (interface{}, error) {
	if c.TaskID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}

	taskService := commandContext.Engine.GetTaskService()
	runtimeService := commandContext.Engine.GetRuntimeService()

	// Get the task to verify it exists
	task, err := taskService.GetTask(ctx, c.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Check if task is suspended
	if task.Suspended {
		return nil, fmt.Errorf("cannot complete suspended task '%s'", c.TaskID)
	}

	// Set variables on the execution if provided
	if c.Variables != nil && len(c.Variables) > 0 && task.ExecutionID != "" {
		if err := runtimeService.SetVariables(ctx, task.ExecutionID, c.Variables); err != nil {
			return nil, fmt.Errorf("failed to set variables: %w", err)
		}
	}

	// Complete the task
	if err := taskService.Complete(ctx, c.TaskID); err != nil {
		return nil, fmt.Errorf("failed to complete task: %w", err)
	}

	// Record to history if enabled
	if commandContext.Engine.GetConfiguration().EnableHistory {
		// TODO: Record historic task instance completion
	}

	// TODO: Continue process execution after task completion
	// This would involve:
	// 1. Finding the next nodes in the process
	// 2. Evaluating conditions on outgoing edges
	// 3. Creating new tasks or executing service tasks

	return nil, nil
}

// NewCompleteTaskCommand creates a new complete task command
func NewCompleteTaskCommand(taskID string, variables map[string]interface{}) *CompleteTaskCommand {
	return &CompleteTaskCommand{
		TaskID:    taskID,
		Variables: variables,
	}
}
