package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/muixstudio/flowgo/api/repository"
)

// Service is the internal implementation of repository.Service
type Service struct {
	databaseDriver string
	databaseURL    string
	deployments    map[string]*repository.Deployment
	definitions    map[string]*repository.ProcessDefinition
	mu             sync.RWMutex
}

// NewService creates a new repository service implementation
func NewService(databaseDriver, databaseURL string) *Service {
	return &Service{
		databaseDriver: databaseDriver,
		databaseURL:    databaseURL,
		deployments:    make(map[string]*repository.Deployment),
		definitions:    make(map[string]*repository.ProcessDefinition),
	}
}

// Initialize initializes the repository service
func (s *Service) Initialize(ctx context.Context) error {
	// TODO: Initialize database connection
	return nil
}

// Shutdown gracefully shuts down the repository service
func (s *Service) Shutdown(ctx context.Context) error {
	// TODO: Close database connections
	return nil
}

// CreateDeployment creates a new deployment builder
func (s *Service) CreateDeployment() repository.DeploymentBuilder {
	// TODO: Return proper builder implementation
	return nil
}

// GetDeployment retrieves a deployment by ID
func (s *Service) GetDeployment(ctx context.Context, deploymentID string) (*repository.Deployment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	deployment, exists := s.deployments[deploymentID]
	if !exists {
		return nil, fmt.Errorf("deployment not found: %s", deploymentID)
	}
	return deployment, nil
}

// DeleteDeployment deletes a deployment
func (s *Service) DeleteDeployment(ctx context.Context, deploymentID string, cascade bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.deployments[deploymentID]; !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	if cascade {
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
func (s *Service) CreateProcessDefinitionQuery() repository.ProcessDefinitionQuery {
	// TODO: Return proper query implementation
	return nil
}

// GetProcessDefinition retrieves a process definition by ID
func (s *Service) GetProcessDefinition(ctx context.Context, processDefinitionID string) (*repository.ProcessDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	def, exists := s.definitions[processDefinitionID]
	if !exists {
		return nil, fmt.Errorf("process definition not found: %s", processDefinitionID)
	}
	return def, nil
}

// GetProcessDefinitionByKey retrieves the latest version by key
func (s *Service) GetProcessDefinitionByKey(ctx context.Context, key string) (*repository.ProcessDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var latestDef *repository.ProcessDefinition
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
func (s *Service) SuspendProcessDefinition(ctx context.Context, processDefinitionID string) error {
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
func (s *Service) ActivateProcessDefinition(ctx context.Context, processDefinitionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	def, exists := s.definitions[processDefinitionID]
	if !exists {
		return fmt.Errorf("process definition not found: %s", processDefinitionID)
	}

	def.Suspended = false
	return nil
}

// GetProcessModel retrieves the process model
func (s *Service) GetProcessModel(ctx context.Context, processDefinitionID string) ([]byte, error) {
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

	for _, resource := range deployment.Resources {
		if resource.Name == def.ResourceName {
			return resource.Content, nil
		}
	}

	return nil, fmt.Errorf("resource not found: %s", def.ResourceName)
}

// ValidateProcessDefinition validates a process definition
func (s *Service) ValidateProcessDefinition(ctx context.Context, content []byte) error {
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

	return nil
}

// DeployInternal is called by DeploymentBuilder
func (s *Service) DeployInternal(ctx context.Context, name, category, tenantID string, resources []*repository.Resource) (*repository.Deployment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	deployment := &repository.Deployment{
		ID:         uuid.New().String(),
		Name:       name,
		DeployTime: time.Now(),
		Category:   category,
		TenantID:   tenantID,
		Resources:  resources,
	}

	// Process resources and create process definitions
	for _, resource := range resources {
		resource.ID = uuid.New().String()
		resource.DeploymentID = deployment.ID

		var processData map[string]interface{}
		if err := json.Unmarshal(resource.Content, &processData); err != nil {
			return nil, fmt.Errorf("failed to parse process definition '%s': %w", resource.Name, err)
		}

		if err := s.ValidateProcessDefinition(ctx, resource.Content); err != nil {
			return nil, fmt.Errorf("invalid process definition '%s': %w", resource.Name, err)
		}

		processID, _ := processData["id"].(string)
		processName, _ := processData["name"].(string)
		processDesc, _ := processData["description"].(string)

		version := 1
		for _, existingDef := range s.definitions {
			if existingDef.Key == processID && existingDef.Version >= version {
				version = existingDef.Version + 1
			}
		}

		processDefinition := &repository.ProcessDefinition{
			ID:                   fmt.Sprintf("%s:%d:%s", processID, version, uuid.New().String()),
			Key:                  processID,
			Name:                 processName,
			Description:          processDesc,
			Version:              version,
			Category:             deployment.Category,
			DeploymentID:         deployment.ID,
			ResourceName:         resource.Name,
			TenantID:             deployment.TenantID,
			Suspended:            false,
			HasGraphicalNotation: true,
		}

		s.definitions[processDefinition.ID] = processDefinition
	}

	s.deployments[deployment.ID] = deployment
	return deployment, nil
}
