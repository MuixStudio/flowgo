package commands

import (
	"context"
	"fmt"

	"github.com/muixstudio/flowgo/engine"
	"github.com/muixstudio/flowgo/runtime"
)

// StartProcessInstanceCommand starts a new process instance
type StartProcessInstanceCommand struct {
	ProcessDefinitionID  string
	ProcessDefinitionKey string
	BusinessKey          string
	Variables            map[string]interface{}
}

// Execute starts the process instance
func (c *StartProcessInstanceCommand) Execute(ctx context.Context, commandContext *engine.CommandContext) (*runtime.ProcessInstance, error) {
	runtimeService := commandContext.Engine.GetRuntimeService()
	repoService := commandContext.Engine.GetRepositoryService()

	var instance *runtime.ProcessInstance
	var err error

	// Start by ID or key
	if c.ProcessDefinitionID != "" {
		// Verify process definition exists and is not suspended
		processDef, err := repoService.GetProcessDefinition(ctx, c.ProcessDefinitionID)
		if err != nil {
			return nil, fmt.Errorf("process definition not found: %w", err)
		}
		if processDef.Suspended {
			return nil, fmt.Errorf("cannot start process instance: process definition '%s' is suspended", c.ProcessDefinitionID)
		}

		if c.BusinessKey != "" {
			instance, err = runtimeService.StartProcessInstanceByKeyWithBusinessKey(ctx, processDef.Key, c.BusinessKey, c.Variables)
		} else {
			instance, err = runtimeService.StartProcessInstanceByID(ctx, c.ProcessDefinitionID, c.Variables)
		}
	} else if c.ProcessDefinitionKey != "" {
		// Get latest version
		processDef, err := repoService.GetProcessDefinitionByKey(ctx, c.ProcessDefinitionKey)
		if err != nil {
			return nil, fmt.Errorf("process definition not found with key '%s': %w", c.ProcessDefinitionKey, err)
		}
		if processDef.Suspended {
			return nil, fmt.Errorf("cannot start process instance: process definition '%s' is suspended", c.ProcessDefinitionKey)
		}

		if c.BusinessKey != "" {
			instance, err = runtimeService.StartProcessInstanceByKeyWithBusinessKey(ctx, c.ProcessDefinitionKey, c.BusinessKey, c.Variables)
		} else {
			instance, err = runtimeService.StartProcessInstanceByKey(ctx, c.ProcessDefinitionKey, c.Variables)
		}
	} else {
		return nil, fmt.Errorf("either process definition ID or key must be provided")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to start process instance: %w", err)
	}

	// Record to history if enabled
	if commandContext.Engine.GetConfiguration().EnableHistory {
		// TODO: Record historic process instance
	}

	return instance, nil
}

// NewStartProcessInstanceByKeyCommand creates a command to start a process by key
func NewStartProcessInstanceByKeyCommand(key string, variables map[string]interface{}) *StartProcessInstanceCommand {
	return &StartProcessInstanceCommand{
		ProcessDefinitionKey: key,
		Variables:            variables,
	}
}

// NewStartProcessInstanceByIDCommand creates a command to start a process by ID
func NewStartProcessInstanceByIDCommand(id string, variables map[string]interface{}) *StartProcessInstanceCommand {
	return &StartProcessInstanceCommand{
		ProcessDefinitionID: id,
		Variables:           variables,
	}
}

// NewStartProcessInstanceWithBusinessKeyCommand creates a command to start a process with business key
func NewStartProcessInstanceWithBusinessKeyCommand(key, businessKey string, variables map[string]interface{}) *StartProcessInstanceCommand {
	return &StartProcessInstanceCommand{
		ProcessDefinitionKey: key,
		BusinessKey:          businessKey,
		Variables:            variables,
	}
}
