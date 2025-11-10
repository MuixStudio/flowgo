package repository

import (
	"context"
	"fmt"
	"time"
)

// RepositoryService provides operations for managing process definitions and deployments.
// This service is responsible for:
// - Deploying process definitions
// - Querying process definitions
// - Managing process definition lifecycle (suspend/activate)
// - Managing deployments
type RepositoryService interface {
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

// Deployment represents a deployment of process definitions
type Deployment struct {
	ID           string
	Name         string
	DeployTime   time.Time
	Category     string
	TenantID     string
	Resources    []*Resource
}

// Resource represents a resource in a deployment (e.g., process definition file)
type Resource struct {
	ID           string
	Name         string
	DeploymentID string
	Content      []byte
	ContentType  string
}

// ProcessDefinition represents a deployed process definition
type ProcessDefinition struct {
	ID                  string
	Key                 string
	Name                string
	Description         string
	Version             int
	Category            string
	DeploymentID        string
	ResourceName        string
	TenantID            string
	Suspended           bool
	StartFormKey        string
	HasStartFormKey     bool
	HasGraphicalNotation bool
}

// DeploymentBuilder provides a fluent API for creating deployments
type DeploymentBuilder struct {
	name      string
	category  string
	tenantID  string
	resources []*Resource
	service   RepositoryService
}

// Name sets the deployment name
func (b *DeploymentBuilder) Name(name string) *DeploymentBuilder {
	b.name = name
	return b
}

// Category sets the deployment category
func (b *DeploymentBuilder) Category(category string) *DeploymentBuilder {
	b.category = category
	return b
}

// TenantID sets the tenant ID for multi-tenancy
func (b *DeploymentBuilder) TenantID(tenantID string) *DeploymentBuilder {
	b.tenantID = tenantID
	return b
}

// AddResource adds a resource to the deployment
func (b *DeploymentBuilder) AddResource(name string, content []byte) *DeploymentBuilder {
	resource := &Resource{
		Name:    name,
		Content: content,
	}
	b.resources = append(b.resources, resource)
	return b
}

// AddProcessDefinition adds a process definition from JSON content
func (b *DeploymentBuilder) AddProcessDefinition(name string, jsonContent []byte) *DeploymentBuilder {
	return b.AddResource(name, jsonContent)
}

// Deploy executes the deployment
func (b *DeploymentBuilder) Deploy(ctx context.Context) (*Deployment, error) {
	// Cast to implementation type to call internal method
	if impl, ok := b.service.(*repositoryServiceImpl); ok {
		return impl.deployInternal(ctx, b)
	}
	return nil, fmt.Errorf("unsupported service implementation")
}

// ProcessDefinitionQuery provides a fluent API for querying process definitions
type ProcessDefinitionQuery struct {
	processDefinitionID  string
	processDefinitionKey string
	processDefinitionName string
	category             string
	deploymentID         string
	tenantID             string
	version              *int
	latestVersion        bool
	suspended            *bool
	orderBy              string
	ascending            bool
	service              RepositoryService
}

// ProcessDefinitionID filters by process definition ID
func (q *ProcessDefinitionQuery) ProcessDefinitionID(id string) *ProcessDefinitionQuery {
	q.processDefinitionID = id
	return q
}

// ProcessDefinitionKey filters by process definition key
func (q *ProcessDefinitionQuery) ProcessDefinitionKey(key string) *ProcessDefinitionQuery {
	q.processDefinitionKey = key
	return q
}

// ProcessDefinitionName filters by process definition name
func (q *ProcessDefinitionQuery) ProcessDefinitionName(name string) *ProcessDefinitionQuery {
	q.processDefinitionName = name
	return q
}

// Category filters by category
func (q *ProcessDefinitionQuery) Category(category string) *ProcessDefinitionQuery {
	q.category = category
	return q
}

// DeploymentID filters by deployment ID
func (q *ProcessDefinitionQuery) DeploymentID(deploymentID string) *ProcessDefinitionQuery {
	q.deploymentID = deploymentID
	return q
}

// TenantID filters by tenant ID
func (q *ProcessDefinitionQuery) TenantID(tenantID string) *ProcessDefinitionQuery {
	q.tenantID = tenantID
	return q
}

// Version filters by specific version
func (q *ProcessDefinitionQuery) Version(version int) *ProcessDefinitionQuery {
	q.version = &version
	return q
}

// LatestVersion filters to only the latest version of each process definition key
func (q *ProcessDefinitionQuery) LatestVersion() *ProcessDefinitionQuery {
	q.latestVersion = true
	return q
}

// Active filters to only active (non-suspended) process definitions
func (q *ProcessDefinitionQuery) Active() *ProcessDefinitionQuery {
	falseVal := false
	q.suspended = &falseVal
	return q
}

// Suspended filters to only suspended process definitions
func (q *ProcessDefinitionQuery) Suspended() *ProcessDefinitionQuery {
	trueVal := true
	q.suspended = &trueVal
	return q
}

// OrderByProcessDefinitionKey orders results by process definition key
func (q *ProcessDefinitionQuery) OrderByProcessDefinitionKey() *ProcessDefinitionQuery {
	q.orderBy = "key"
	return q
}

// OrderByProcessDefinitionName orders results by process definition name
func (q *ProcessDefinitionQuery) OrderByProcessDefinitionName() *ProcessDefinitionQuery {
	q.orderBy = "name"
	return q
}

// OrderByDeploymentID orders results by deployment ID
func (q *ProcessDefinitionQuery) OrderByDeploymentID() *ProcessDefinitionQuery {
	q.orderBy = "deployment_id"
	return q
}

// Asc sets ascending order
func (q *ProcessDefinitionQuery) Asc() *ProcessDefinitionQuery {
	q.ascending = true
	return q
}

// Desc sets descending order
func (q *ProcessDefinitionQuery) Desc() *ProcessDefinitionQuery {
	q.ascending = false
	return q
}

// List executes the query and returns a list of process definitions
func (q *ProcessDefinitionQuery) List(ctx context.Context) ([]*ProcessDefinition, error) {
	// Will be implemented by the concrete service
	return nil, nil
}

// Count returns the count of matching process definitions
func (q *ProcessDefinitionQuery) Count(ctx context.Context) (int64, error) {
	// Will be implemented by the concrete service
	return 0, nil
}

// SingleResult returns a single process definition or error if not exactly one result
func (q *ProcessDefinitionQuery) SingleResult(ctx context.Context) (*ProcessDefinition, error) {
	// Will be implemented by the concrete service
	return nil, nil
}
