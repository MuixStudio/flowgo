package commands

import (
	"context"
	"fmt"

	"github.com/muixstudio/flowgo/engine"
	"github.com/muixstudio/flowgo/repository"
)

// DeployCommand deploys a process definition
type DeployCommand struct {
	DeploymentName   string
	Category         string
	TenantID         string
	ResourceName     string
	ResourceContent  []byte
}

// Execute deploys the process definition
func (c *DeployCommand) Execute(ctx context.Context, commandContext *engine.CommandContext) (*repository.Deployment, error) {
	if c.ResourceContent == nil || len(c.ResourceContent) == 0 {
		return nil, fmt.Errorf("resource content cannot be empty")
	}

	if c.ResourceName == "" {
		return nil, fmt.Errorf("resource name cannot be empty")
	}

	// Get repository service
	repoService := commandContext.Engine.GetRepositoryService()

	// Create deployment
	deployment, err := repoService.CreateDeployment().
		Name(c.DeploymentName).
		Category(c.Category).
		TenantID(c.TenantID).
		AddProcessDefinition(c.ResourceName, c.ResourceContent).
		Deploy(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to deploy process definition: %w", err)
	}

	return deployment, nil
}

// NewDeployCommand creates a new deploy command
func NewDeployCommand(name, resourceName string, content []byte) *DeployCommand {
	return &DeployCommand{
		DeploymentName:  name,
		ResourceName:    resourceName,
		ResourceContent: content,
	}
}
