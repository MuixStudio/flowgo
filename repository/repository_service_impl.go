package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// repositoryServiceImpl is the default implementation of RepositoryService
type repositoryServiceImpl struct {
	databaseDriver string
	databaseURL    string
	deployments    map[string]*Deployment
	definitions    map[string]*ProcessDefinition
	mu             sync.RWMutex
}

// NewRepositoryService creates a new repository service
func NewRepositoryService(databaseDriver, databaseURL string) RepositoryService {
	return &repositoryServiceImpl{
		databaseDriver: databaseDriver,
		databaseURL:    databaseURL,
		deployments:    make(map[string]*Deployment),
		definitions:    make(map[string]*ProcessDefinition),
	}
}

// Initialize initializes the repository service
func (s *repositoryServiceImpl) Initialize(ctx context.Context) error {
	// TODO: Initialize database connection and create tables if needed
	return nil
}

// Shutdown gracefully shuts down the repository service
func (s *repositoryServiceImpl) Shutdown(ctx context.Context) error {
	// TODO: Close database connections
	return nil
}

// CreateDeployment creates a new deployment builder
func (s *repositoryServiceImpl) CreateDeployment() *DeploymentBuilder {
	return &DeploymentBuilder{
		service:   s,
		resources: make([]*Resource, 0),
	}
}

// GetDeployment retrieves a deployment by ID
func (s *repositoryServiceImpl) GetDeployment(ctx context.Context, deploymentID string) (*Deployment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	deployment, exists := s.deployments[deploymentID]
	if !exists {
		return nil, fmt.Errorf("deployment not found: %s", deploymentID)
	}
	return deployment, nil
}

// DeleteDeployment deletes a deployment
func (s *repositoryServiceImpl) DeleteDeployment(ctx context.Context, deploymentID string, cascade bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	deployment, exists := s.deployments[deploymentID]
	if !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	if cascade {
		// Delete all process definitions related to this deployment
		for id, def := range s.definitions {
			if def.DeploymentID == deploymentID {
				delete(s.definitions, id)
			}
		}
	}

	delete(s.deployments, deploymentID)
	return nil
}

// CreateProcessDefinitionQuery creates a new process definition query
func (s *repositoryServiceImpl) CreateProcessDefinitionQuery() *ProcessDefinitionQuery {
	return &ProcessDefinitionQuery{
		service: s,
	}
}

// GetProcessDefinition retrieves a process definition by ID
func (s *repositoryServiceImpl) GetProcessDefinition(ctx context.Context, processDefinitionID string) (*ProcessDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	def, exists := s.definitions[processDefinitionID]
	if !exists {
		return nil, fmt.Errorf("process definition not found: %s", processDefinitionID)
	}
	return def, nil
}

// GetProcessDefinitionByKey retrieves the latest version of a process definition by key
func (s *repositoryServiceImpl) GetProcessDefinitionByKey(ctx context.Context, key string) (*ProcessDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var latestDef *ProcessDefinition
	for _, def := range s.definitions {
		if def.Key == key {
			if latestDef == nil || def.Version > latestDef.Version {
				latestDef = def
			}
		}
	}

	if latestDef == nil {
		return nil, fmt.Errorf("process definition not found with key: %s", key)
	}
	return latestDef, nil
}

// SuspendProcessDefinition suspends a process definition
func (s *repositoryServiceImpl) SuspendProcessDefinition(ctx context.Context, processDefinitionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	def, exists := s.definitions[processDefinitionID]
	if !exists {
		return fmt.Errorf("process definition not found: %s", processDefinitionID)
	}

	def.Suspended = true
	return nil
}

// ActivateProcessDefinition activates a suspended process definition
func (s *repositoryServiceImpl) ActivateProcessDefinition(ctx context.Context, processDefinitionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	def, exists := s.definitions[processDefinitionID]
	if !exists {
		return fmt.Errorf("process definition not found: %s", processDefinitionID)
	}

	def.Suspended = false
	return nil
}

// GetProcessModel retrieves the process model for a process definition
func (s *repositoryServiceImpl) GetProcessModel(ctx context.Context, processDefinitionID string) ([]byte, error) {
	s.mu.RLock()
	def, exists := s.definitions[processDefinitionID]
	if !exists {
		s.mu.RUnlock()
		return nil, fmt.Errorf("process definition not found: %s", processDefinitionID)
	}

	deployment, exists := s.deployments[def.DeploymentID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("deployment not found: %s", def.DeploymentID)
	}

	// Find the resource with the matching name
	for _, resource := range deployment.Resources {
		if resource.Name == def.ResourceName {
			return resource.Content, nil
		}
	}

	return nil, fmt.Errorf("resource not found: %s", def.ResourceName)
}

// ValidateProcessDefinition validates a process definition without deploying it
func (s *repositoryServiceImpl) ValidateProcessDefinition(ctx context.Context, content []byte) error {
	// Parse the JSON content
	var processData map[string]interface{}
	if err := json.Unmarshal(content, &processData); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Validate required fields
	if _, ok := processData["id"]; !ok {
		return fmt.Errorf("process definition must have an 'id' field")
	}
	if _, ok := processData["name"]; !ok {
		return fmt.Errorf("process definition must have a 'name' field")
	}
	if _, ok := processData["nodes"]; !ok {
		return fmt.Errorf("process definition must have a 'nodes' field")
	}
	if _, ok := processData["edges"]; !ok {
		return fmt.Errorf("process definition must have an 'edges' field")
	}

	// TODO: Add more comprehensive validation
	// - Validate node types
	// - Validate edge connections
	// - Validate required properties per node type
	// - Check for cycles
	// - Ensure start and end events exist

	return nil
}

// deployInternal is called by DeploymentBuilder to execute the deployment
func (s *repositoryServiceImpl) deployInternal(ctx context.Context, builder *DeploymentBuilder) (*Deployment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create deployment
	deployment := &Deployment{
		ID:         uuid.New().String(),
		Name:       builder.name,
		DeployTime: time.Now(),
		Category:   builder.category,
		TenantID:   builder.tenantID,
		Resources:  builder.resources,
	}

	// Process each resource to create process definitions
	for _, resource := range builder.resources {
		resource.ID = uuid.New().String()
		resource.DeploymentID = deployment.ID

		// Parse process definition from JSON
		var processData map[string]interface{}
		if err := json.Unmarshal(resource.Content, &processData); err != nil {
			return nil, fmt.Errorf("failed to parse process definition '%s': %w", resource.Name, err)
		}

		// Validate process definition
		if err := s.ValidateProcessDefinition(ctx, resource.Content); err != nil {
			return nil, fmt.Errorf("invalid process definition '%s': %w", resource.Name, err)
		}

		// Extract process definition details
		processID, _ := processData["id"].(string)
		processName, _ := processData["name"].(string)
		processDesc, _ := processData["description"].(string)

		// Calculate version - find existing versions with the same key
		version := 1
		for _, existingDef := range s.definitions {
			if existingDef.Key == processID && existingDef.Version >= version {
				version = existingDef.Version + 1
			}
		}

		// Create process definition
		processDefinition := &ProcessDefinition{
			ID:                  fmt.Sprintf("%s:%d:%s", processID, version, uuid.New().String()),
			Key:                 processID,
			Name:                processName,
			Description:         processDesc,
			Version:             version,
			Category:            deployment.Category,
			DeploymentID:        deployment.ID,
			ResourceName:        resource.Name,
			TenantID:            deployment.TenantID,
			Suspended:           false,
			HasGraphicalNotation: true,
		}

		s.definitions[processDefinition.ID] = processDefinition
	}

	s.deployments[deployment.ID] = deployment
	return deployment, nil
}
