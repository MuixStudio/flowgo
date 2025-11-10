package repository

import "context"

// Service provides operations for managing process definitions and deployments.
// This is the public interface exposed to users.
type Service interface {
	// Initialize initializes the repository service
	Initialize(ctx context.Context) error

	// Shutdown gracefully shuts down the repository service
	Shutdown(ctx context.Context) error

	// CreateDeployment creates a new deployment builder
	CreateDeployment() *DeploymentBuilder

	// GetDeployment retrieves a deployment by ID
	GetDeployment(ctx context.Context, deploymentID string) (*Deployment, error)

	// DeleteDeployment deletes a deployment and optionally cascade delete related data
	DeleteDeployment(ctx context.Context, deploymentID string, cascade bool) error

	// CreateProcessDefinitionQuery creates a new process definition query
	CreateProcessDefinitionQuery() *ProcessDefinitionQuery

	// GetProcessDefinition retrieves a process definition by ID
	GetProcessDefinition(ctx context.Context, processDefinitionID string) (*ProcessDefinition, error)

	// GetProcessDefinitionByKey retrieves the latest version of a process definition by key
	GetProcessDefinitionByKey(ctx context.Context, key string) (*ProcessDefinition, error)

	// SuspendProcessDefinition suspends a process definition
	SuspendProcessDefinition(ctx context.Context, processDefinitionID string) error

	// ActivateProcessDefinition activates a suspended process definition
	ActivateProcessDefinition(ctx context.Context, processDefinitionID string) error

	// GetProcessModel retrieves the process model (JSON content) for a process definition
	GetProcessModel(ctx context.Context, processDefinitionID string) ([]byte, error)

	// ValidateProcessDefinition validates a process definition without deploying it
	ValidateProcessDefinition(ctx context.Context, content []byte) error
}
