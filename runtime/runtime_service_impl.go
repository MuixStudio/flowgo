package runtime

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/muixstudio/flowgo/repository"
)

// runtimeServiceImpl is the default implementation of RuntimeService
type runtimeServiceImpl struct {
	repositoryService repository.RepositoryService
	enableAsync       bool
	processInstances  map[string]*ProcessInstance
	executions        map[string]*Execution
	variables         map[string]map[string]interface{} // executionID -> variables
	mu                sync.RWMutex
}

// NewRuntimeService creates a new runtime service
func NewRuntimeService(repositoryService repository.RepositoryService, enableAsync bool) RuntimeService {
	return &runtimeServiceImpl{
		repositoryService: repositoryService,
		enableAsync:       enableAsync,
		processInstances:  make(map[string]*ProcessInstance),
		executions:        make(map[string]*Execution),
		variables:         make(map[string]map[string]interface{}),
	}
}

// Initialize initializes the runtime service
func (s *runtimeServiceImpl) Initialize(ctx context.Context) error {
	// TODO: Initialize async executor if enabled
	return nil
}

// Shutdown gracefully shuts down the runtime service
func (s *runtimeServiceImpl) Shutdown(ctx context.Context) error {
	// TODO: Stop async executor
	return nil
}

// StartProcessInstanceByKey starts a process instance by process definition key
func (s *runtimeServiceImpl) StartProcessInstanceByKey(ctx context.Context, processDefinitionKey string, variables map[string]interface{}) (*ProcessInstance, error) {
	// Get the latest process definition by key
	processDefinition, err := s.repositoryService.GetProcessDefinitionByKey(ctx, processDefinitionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get process definition: %w", err)
	}

	return s.startProcessInstance(ctx, processDefinition, "", variables)
}

// StartProcessInstanceByID starts a process instance by process definition ID
func (s *runtimeServiceImpl) StartProcessInstanceByID(ctx context.Context, processDefinitionID string, variables map[string]interface{}) (*ProcessInstance, error) {
	processDefinition, err := s.repositoryService.GetProcessDefinition(ctx, processDefinitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get process definition: %w", err)
	}

	return s.startProcessInstance(ctx, processDefinition, "", variables)
}

// StartProcessInstanceByKeyWithBusinessKey starts a process instance with a business key
func (s *runtimeServiceImpl) StartProcessInstanceByKeyWithBusinessKey(ctx context.Context, processDefinitionKey, businessKey string, variables map[string]interface{}) (*ProcessInstance, error) {
	processDefinition, err := s.repositoryService.GetProcessDefinitionByKey(ctx, processDefinitionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get process definition: %w", err)
	}

	return s.startProcessInstance(ctx, processDefinition, businessKey, variables)
}

// startProcessInstance is the internal method to start a process instance
func (s *runtimeServiceImpl) startProcessInstance(ctx context.Context, processDefinition *repository.ProcessDefinition, businessKey string, variables map[string]interface{}) (*ProcessInstance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if process definition is suspended
	if processDefinition.Suspended {
		return nil, fmt.Errorf("process definition '%s' is suspended", processDefinition.ID)
	}

	// Create process instance
	processInstance := &ProcessInstance{
		ID:                   uuid.New().String(),
		ProcessDefinitionID:  processDefinition.ID,
		ProcessDefinitionKey: processDefinition.Key,
		ProcessDefinitionName: processDefinition.Name,
		BusinessKey:          businessKey,
		StartTime:            time.Now(),
		TenantID:             processDefinition.TenantID,
		RootProcessInstanceID: "",
	}
	processInstance.RootProcessInstanceID = processInstance.ID

	// Create root execution
	execution := &Execution{
		ID:                processInstance.ID,
		ProcessInstanceID: processInstance.ID,
		IsActive:          true,
		IsScope:           true,
		TenantID:          processDefinition.TenantID,
	}

	// Store process instance and execution
	s.processInstances[processInstance.ID] = processInstance
	s.executions[execution.ID] = execution

	// Initialize variables
	if variables != nil {
		s.variables[execution.ID] = make(map[string]interface{})
		for k, v := range variables {
			s.variables[execution.ID][k] = v
		}
	}

	// TODO: Execute the process (navigate through nodes)
	// This would involve:
	// 1. Finding the start event
	// 2. Creating executions for each path
	// 3. Processing nodes (tasks, gateways, etc.)
	// 4. Managing the execution state

	return processInstance, nil
}

// DeleteProcessInstance deletes a process instance
func (s *runtimeServiceImpl) DeleteProcessInstance(ctx context.Context, processInstanceID, deleteReason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.processInstances[processInstanceID]; !exists {
		return fmt.Errorf("process instance not found: %s", processInstanceID)
	}

	// Delete all executions for this process instance
	for id, exec := range s.executions {
		if exec.ProcessInstanceID == processInstanceID {
			delete(s.executions, id)
			delete(s.variables, id)
		}
	}

	delete(s.processInstances, processInstanceID)
	return nil
}

// SuspendProcessInstance suspends a process instance
func (s *runtimeServiceImpl) SuspendProcessInstance(ctx context.Context, processInstanceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	processInstance, exists := s.processInstances[processInstanceID]
	if !exists {
		return fmt.Errorf("process instance not found: %s", processInstanceID)
	}

	processInstance.Suspended = true
	return nil
}

// ActivateProcessInstance activates a suspended process instance
func (s *runtimeServiceImpl) ActivateProcessInstance(ctx context.Context, processInstanceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	processInstance, exists := s.processInstances[processInstanceID]
	if !exists {
		return fmt.Errorf("process instance not found: %s", processInstanceID)
	}

	processInstance.Suspended = false
	return nil
}

// CreateProcessInstanceQuery creates a new process instance query
func (s *runtimeServiceImpl) CreateProcessInstanceQuery() *ProcessInstanceQuery {
	return &ProcessInstanceQuery{
		service: s,
	}
}

// GetProcessInstance retrieves a process instance by ID
func (s *runtimeServiceImpl) GetProcessInstance(ctx context.Context, processInstanceID string) (*ProcessInstance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	processInstance, exists := s.processInstances[processInstanceID]
	if !exists {
		return nil, fmt.Errorf("process instance not found: %s", processInstanceID)
	}
	return processInstance, nil
}

// SetVariable sets a variable on a process instance
func (s *runtimeServiceImpl) SetVariable(ctx context.Context, executionID, variableName string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.executions[executionID]; !exists {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	if s.variables[executionID] == nil {
		s.variables[executionID] = make(map[string]interface{})
	}

	s.variables[executionID][variableName] = value
	return nil
}

// SetVariables sets multiple variables on a process instance
func (s *runtimeServiceImpl) SetVariables(ctx context.Context, executionID string, variables map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.executions[executionID]; !exists {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	if s.variables[executionID] == nil {
		s.variables[executionID] = make(map[string]interface{})
	}

	for k, v := range variables {
		s.variables[executionID][k] = v
	}
	return nil
}

// GetVariable gets a variable from a process instance
func (s *runtimeServiceImpl) GetVariable(ctx context.Context, executionID, variableName string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.executions[executionID]; !exists {
		return nil, fmt.Errorf("execution not found: %s", executionID)
	}

	if s.variables[executionID] == nil {
		return nil, nil
	}

	return s.variables[executionID][variableName], nil
}

// GetVariables gets all variables from a process instance
func (s *runtimeServiceImpl) GetVariables(ctx context.Context, executionID string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.executions[executionID]; !exists {
		return nil, fmt.Errorf("execution not found: %s", executionID)
	}

	// Return a copy to avoid concurrent modification
	result := make(map[string]interface{})
	if s.variables[executionID] != nil {
		for k, v := range s.variables[executionID] {
			result[k] = v
		}
	}
	return result, nil
}

// RemoveVariable removes a variable from a process instance
func (s *runtimeServiceImpl) RemoveVariable(ctx context.Context, executionID, variableName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.executions[executionID]; !exists {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	if s.variables[executionID] != nil {
		delete(s.variables[executionID], variableName)
	}
	return nil
}

// Signal triggers a signal event
func (s *runtimeServiceImpl) Signal(ctx context.Context, executionID string) error {
	return s.SignalWithVariables(ctx, executionID, nil)
}

// SignalWithVariables triggers a signal event with variables
func (s *runtimeServiceImpl) SignalWithVariables(ctx context.Context, executionID string, variables map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	execution, exists := s.executions[executionID]
	if !exists {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	// Set variables if provided
	if variables != nil {
		if s.variables[executionID] == nil {
			s.variables[executionID] = make(map[string]interface{})
		}
		for k, v := range variables {
			s.variables[executionID][k] = v
		}
	}

	// TODO: Continue execution from this point
	// This would involve finding the next nodes and processing them
	_ = execution

	return nil
}

// CreateExecutionQuery creates a new execution query
func (s *runtimeServiceImpl) CreateExecutionQuery() *ExecutionQuery {
	return &ExecutionQuery{
		service: s,
	}
}
