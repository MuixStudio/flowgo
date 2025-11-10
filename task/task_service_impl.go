package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/muixstudio/flowgo/commands"
	"github.com/muixstudio/flowgo/engine"
	"github.com/muixstudio/flowgo/runtime"
)

// taskServiceImpl is the default implementation of TaskService
type taskServiceImpl struct {
	runtimeService runtime.RuntimeService
	executor       engine.CommandExecutor
	tasks          map[string]*Task
	comments       map[string][]*Comment             // taskID -> comments
	attachments    map[string][]*Attachment          // taskID -> attachments
	variables      map[string]map[string]interface{} // taskID -> variables
	mu             sync.RWMutex
}

// NewTaskService creates a new task service
func NewTaskService(runtimeService runtime.RuntimeService) TaskService {
	cmdExec := engine.NewCommandExecutor()
	return &taskServiceImpl{
		runtimeService: runtimeService,
		executor:       cmdExec,
		tasks:          make(map[string]*Task),
		comments:       make(map[string][]*Comment),
		attachments:    make(map[string][]*Attachment),
		variables:      make(map[string]map[string]interface{}),
	}
}

// Initialize initializes the task service
func (s *taskServiceImpl) Initialize(ctx context.Context) error {
	return nil
}

// Shutdown gracefully shuts down the task service
func (s *taskServiceImpl) Shutdown(ctx context.Context) error {
	return nil
}

// CreateTaskQuery creates a new task query
func (s *taskServiceImpl) CreateTaskQuery() *TaskQuery {
	return &TaskQuery{
		service: s,
	}
}

// GetTask retrieves a task by ID
func (s *taskServiceImpl) GetTask(ctx context.Context, taskID string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, err := s.executor.Execute(ctx, &commands.StartProcessInstanceCommand{})
	return task, err
}

// NewTask creates a new standalone task
func (s *taskServiceImpl) NewTask(ctx context.Context, taskID string) (*Task, error) {
	if taskID == "" {
		taskID = uuid.New().String()
	}

	task := &Task{
		ID:         taskID,
		CreateTime: time.Now(),
		Priority:   5, // Default priority
	}
	return task, nil
}

// SaveTask saves a standalone task
func (s *taskServiceImpl) SaveTask(ctx context.Context, task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	s.tasks[task.ID] = task
	return nil
}

// DeleteTask deletes a task
func (s *taskServiceImpl) DeleteTask(ctx context.Context, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskID]; !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	delete(s.tasks, taskID)
	delete(s.comments, taskID)
	delete(s.attachments, taskID)
	delete(s.variables, taskID)
	return nil
}

// Claim assigns a task to a specific user
func (s *taskServiceImpl) Claim(ctx context.Context, taskID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if task.Assignee != "" && task.Assignee != userID {
		return fmt.Errorf("task is already claimed by another user: %s", task.Assignee)
	}

	now := time.Now()
	task.Assignee = userID
	task.ClaimTime = &now
	return nil
}

// Unclaim removes the assignee from a task
func (s *taskServiceImpl) Unclaim(ctx context.Context, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Assignee = ""
	task.ClaimTime = nil
	return nil
}

// Complete completes a task
func (s *taskServiceImpl) Complete(ctx context.Context, taskID string) error {
	return s.CompleteWithVariables(ctx, taskID, nil)
}

// CompleteWithVariables completes a task and sets variables
func (s *taskServiceImpl) CompleteWithVariables(ctx context.Context, taskID string, variables map[string]interface{}) error {
	s.mu.Lock()
	task, exists := s.tasks[taskID]
	s.mu.Unlock()

	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Set variables on the execution
	if variables != nil && task.ExecutionID != "" {
		if err := s.runtimeService.SetVariables(ctx, task.ExecutionID, variables); err != nil {
			return fmt.Errorf("failed to set variables: %w", err)
		}
	}

	// TODO: Signal the execution to continue
	if task.ExecutionID != "" {
		if err := s.runtimeService.Signal(ctx, task.ExecutionID); err != nil {
			return fmt.Errorf("failed to signal execution: %w", err)
		}
	}

	// Delete the task
	s.mu.Lock()
	delete(s.tasks, taskID)
	s.mu.Unlock()

	return nil
}

// SetAssignee sets the assignee of a task
func (s *taskServiceImpl) SetAssignee(ctx context.Context, taskID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Assignee = userID
	return nil
}

// SetOwner sets the owner of a task
func (s *taskServiceImpl) SetOwner(ctx context.Context, taskID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Owner = userID
	return nil
}

// AddCandidateUser adds a candidate user to a task
func (s *taskServiceImpl) AddCandidateUser(ctx context.Context, taskID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Check if user already exists
	for _, u := range task.CandidateUsers {
		if u == userID {
			return nil // Already exists
		}
	}

	task.CandidateUsers = append(task.CandidateUsers, userID)
	return nil
}

// AddCandidateGroup adds a candidate group to a task
func (s *taskServiceImpl) AddCandidateGroup(ctx context.Context, taskID, groupID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Check if group already exists
	for _, g := range task.CandidateGroups {
		if g == groupID {
			return nil // Already exists
		}
	}

	task.CandidateGroups = append(task.CandidateGroups, groupID)
	return nil
}

// DeleteCandidateUser removes a candidate user from a task
func (s *taskServiceImpl) DeleteCandidateUser(ctx context.Context, taskID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	for i, u := range task.CandidateUsers {
		if u == userID {
			task.CandidateUsers = append(task.CandidateUsers[:i], task.CandidateUsers[i+1:]...)
			break
		}
	}

	return nil
}

// DeleteCandidateGroup removes a candidate group from a task
func (s *taskServiceImpl) DeleteCandidateGroup(ctx context.Context, taskID, groupID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	for i, g := range task.CandidateGroups {
		if g == groupID {
			task.CandidateGroups = append(task.CandidateGroups[:i], task.CandidateGroups[i+1:]...)
			break
		}
	}

	return nil
}

// SetPriority sets the priority of a task
func (s *taskServiceImpl) SetPriority(ctx context.Context, taskID string, priority int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Priority = priority
	return nil
}

// SetDueDate sets the due date of a task
func (s *taskServiceImpl) SetDueDate(ctx context.Context, taskID string, dueDate time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.DueDate = &dueDate
	return nil
}

// GetTaskVariables gets all variables of a task
func (s *taskServiceImpl) GetTaskVariables(ctx context.Context, taskID string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.tasks[taskID]; !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	// Return a copy
	result := make(map[string]interface{})
	if s.variables[taskID] != nil {
		for k, v := range s.variables[taskID] {
			result[k] = v
		}
	}
	return result, nil
}

// GetTaskVariable gets a specific variable of a task
func (s *taskServiceImpl) GetTaskVariable(ctx context.Context, taskID, variableName string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.tasks[taskID]; !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	if s.variables[taskID] == nil {
		return nil, nil
	}

	return s.variables[taskID][variableName], nil
}

// SetTaskVariable sets a variable on a task
func (s *taskServiceImpl) SetTaskVariable(ctx context.Context, taskID, variableName string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskID]; !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if s.variables[taskID] == nil {
		s.variables[taskID] = make(map[string]interface{})
	}

	s.variables[taskID][variableName] = value
	return nil
}

// SetTaskVariables sets multiple variables on a task
func (s *taskServiceImpl) SetTaskVariables(ctx context.Context, taskID string, variables map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskID]; !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if s.variables[taskID] == nil {
		s.variables[taskID] = make(map[string]interface{})
	}

	for k, v := range variables {
		s.variables[taskID][k] = v
	}
	return nil
}

// RemoveTaskVariable removes a variable from a task
func (s *taskServiceImpl) RemoveTaskVariable(ctx context.Context, taskID, variableName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskID]; !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if s.variables[taskID] != nil {
		delete(s.variables[taskID], variableName)
	}
	return nil
}

// AddComment adds a comment to a task
func (s *taskServiceImpl) AddComment(ctx context.Context, taskID, message string) (*Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskID]; !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	comment := &Comment{
		ID:      uuid.New().String(),
		TaskID:  taskID,
		Message: message,
		Time:    time.Now(),
	}

	s.comments[taskID] = append(s.comments[taskID], comment)
	return comment, nil
}

// GetTaskComments gets all comments for a task
func (s *taskServiceImpl) GetTaskComments(ctx context.Context, taskID string) ([]*Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.tasks[taskID]; !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	return s.comments[taskID], nil
}

// CreateAttachment creates an attachment for a task
func (s *taskServiceImpl) CreateAttachment(ctx context.Context, taskID, attachmentType, attachmentName, attachmentDescription string, content []byte) (*Attachment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	attachment := &Attachment{
		ID:                uuid.New().String(),
		Name:              attachmentName,
		Description:       attachmentDescription,
		Type:              attachmentType,
		TaskID:            taskID,
		ProcessInstanceID: task.ProcessInstanceID,
		Content:           content,
		Time:              time.Now(),
	}

	s.attachments[taskID] = append(s.attachments[taskID], attachment)
	return attachment, nil
}

// GetTaskAttachments gets all attachments for a task
func (s *taskServiceImpl) GetTaskAttachments(ctx context.Context, taskID string) ([]*Attachment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.tasks[taskID]; !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	return s.attachments[taskID], nil
}

// DeleteAttachment deletes an attachment
func (s *taskServiceImpl) DeleteAttachment(ctx context.Context, attachmentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find and delete the attachment
	for taskID, attachments := range s.attachments {
		for i, att := range attachments {
			if att.ID == attachmentID {
				s.attachments[taskID] = append(attachments[:i], attachments[i+1:]...)
				return nil
			}
		}
	}

	return fmt.Errorf("attachment not found: %s", attachmentID)
}
