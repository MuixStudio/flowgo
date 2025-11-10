package history

import "context"

// Service provides operations for querying historical process data.
type Service interface {
	// Initialize initializes the history service
	Initialize(ctx context.Context) error

	// Shutdown gracefully shuts down the history service
	Shutdown(ctx context.Context) error

	// CreateHistoricProcessInstanceQuery creates a new historic process instance query
	CreateHistoricProcessInstanceQuery() *HistoricProcessInstanceQuery

	// CreateHistoricTaskInstanceQuery creates a new historic task instance query
	CreateHistoricTaskInstanceQuery() *HistoricTaskInstanceQuery

	// DeleteHistoricProcessInstance deletes a historic process instance
	DeleteHistoricProcessInstance(ctx context.Context, processInstanceID string) error
}
