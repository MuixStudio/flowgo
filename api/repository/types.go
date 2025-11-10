package repository

import (
	"context"
	"fmt"
	"time"
)

// Deployment represents a deployment of process definitions
type Deployment struct {
	ID         string
	Name       string
	DeployTime time.Time
	Category   string
	TenantID   string
	Resources  []*Resource
}

// Resource represents a resource in a deployment
type Resource struct {
	ID           string
	Name         string
	DeploymentID string
	Content      []byte
	ContentType  string
}

// ProcessDefinition represents a deployed process definition
type ProcessDefinition struct {
	ID                   string
	Key                  string
	Name                 string
	Description          string
	Version              int
	Category             string
	DeploymentID         string
	ResourceName         string
	TenantID             string
	Suspended            bool
	StartFormKey         string
	HasStartFormKey      bool
	HasGraphicalNotation bool
}

// DeploymentBuilder provides a fluent API for creating deployments
type DeploymentBuilder struct {
	name      string
	category  string
	tenantID  string
	resources []*Resource
	service   Service
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

// TenantID sets the tenant ID
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
	// This will be implemented by calling internal implementation
	if b.service == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	// TODO: Call internal implementation
	return nil, fmt.Errorf("not implemented")
}

// ProcessDefinitionQuery provides a fluent API for querying process definitions
type ProcessDefinitionQuery struct {
	processDefinitionID   string
	processDefinitionKey  string
	processDefinitionName string
	category              string
	deploymentID          string
	tenantID              string
	version               *int
	latestVersion         bool
	suspended             *bool
	orderBy               string
	ascending             bool
	service               Service
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

// LatestVersion filters to only the latest version
func (q *ProcessDefinitionQuery) LatestVersion() *ProcessDefinitionQuery {
	q.latestVersion = true
	return q
}

// Active filters to only active process definitions
func (q *ProcessDefinitionQuery) Active() *ProcessDefinitionQuery {
	falseVal := false
	q.suspended = &falseVal
	return q
}

// List executes the query and returns a list of process definitions
func (q *ProcessDefinitionQuery) List(ctx context.Context) ([]*ProcessDefinition, error) {
	// TODO: Call internal implementation
	return nil, nil
}

// Count returns the count of matching process definitions
func (q *ProcessDefinitionQuery) Count(ctx context.Context) (int64, error) {
	// TODO: Call internal implementation
	return 0, nil
}
