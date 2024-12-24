package flowgo

type ProcessDefinition struct {
	ID           string
	Key          string
	Name         string
	Version      int
	Category     string
	DeploymentID string
	Suspended    bool
}

func NewProcessDefinition() *ProcessDefinition {
	return &ProcessDefinition{}
}
