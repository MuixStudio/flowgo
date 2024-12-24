package flowgo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Deployment struct {
	ID        string
	Name      string
	Resources []string
	CreatedAt time.Time
}

type Model struct {
	ID       string
	Name     string
	Key      string
	Category string
	Version  int
	MetaInfo string
	Editor   []byte // 相当于 Flowable 的 JSON
}

type RepositoryService struct {
	deployments map[string]*Deployment
	processDefs map[string]*ProcessDefinition
	models      map[string]*Model
}

func NewRepositoryService() *RepositoryService {
	return &RepositoryService{
		deployments: make(map[string]*Deployment),
		processDefs: make(map[string]*ProcessDefinition),
		models:      make(map[string]*Model),
	}
}

func (r *RepositoryService) CreateDeployment(name string) *Deployment {
	d := &Deployment{
		ID:        uuid.New().String(),
		Name:      name,
		Resources: []string{},
		CreatedAt: time.Now(),
	}
	r.deployments[d.ID] = d
	return d
}

func (r *RepositoryService) GetDeployment(deploymentID string) (*Deployment, error) {
	if d, ok := r.deployments[deploymentID]; ok {
		return d, nil
	}
	return nil, errors.New("deployment not found")
}

func (r *RepositoryService) DeleteDeployment(deploymentID string, cascade bool) error {
	if _, ok := r.deployments[deploymentID]; !ok {
		return errors.New("deployment not found")
	}
	delete(r.deployments, deploymentID)

	// 如果 cascade=true，要删除相关流程定义
	if cascade {
		for id, def := range r.processDefs {
			if def.DeploymentID == deploymentID {
				delete(r.processDefs, id)
			}
		}
	}
	return nil
}

func (r *RepositoryService) ListDeployments() []*Deployment {
	res := []*Deployment{}
	for _, d := range r.deployments {
		res = append(res, d)
	}
	return res
}

func (r *RepositoryService) SaveProcessDefinition(key, name, depID string) *ProcessDefinition {
	def := &ProcessDefinition{
		ID:           uuid.New().String(),
		Key:          key,
		Name:         name,
		Version:      1, // 简单起见，自动设定
		DeploymentID: depID,
	}
	r.processDefs[def.ID] = def
	return def
}

func (r *RepositoryService) GetProcessDefinition(processDefinitionID string) (*ProcessDefinition, error) {
	if def, ok := r.processDefs[processDefinitionID]; ok {
		return def, nil
	}
	return nil, errors.New("process definition not found")
}

func (r *RepositoryService) GetProcessDefinitionsByKey(key string) []*ProcessDefinition {
	var defs []*ProcessDefinition
	for _, d := range r.processDefs {
		if d.Key == key {
			defs = append(defs, d)
		}
	}
	return defs
}

func (r *RepositoryService) SuspendProcessDefinitionById(id string) error {
	if def, ok := r.processDefs[id]; ok {
		def.Suspended = true
		return nil
	}
	return errors.New("process definition not found")
}

func (r *RepositoryService) ActivateProcessDefinitionById(id string) error {
	if def, ok := r.processDefs[id]; ok {
		def.Suspended = false
		return nil
	}
	return errors.New("process definition not found")
}

func (r *RepositoryService) SuspendProcessDefinitionByKey(key string) error {
	found := false
	for _, def := range r.processDefs {
		if def.Key == key {
			def.Suspended = true
			found = true
		}
	}
	if !found {
		return errors.New("process definition not found")
	}
	return nil
}

func (r *RepositoryService) ActivateProcessDefinitionByKey(key string) error {
	found := false
	for _, def := range r.processDefs {
		if def.Key == key {
			def.Suspended = false
			found = true
		}
	}
	if !found {
		return errors.New("process definition not found")
	}
	return nil
}

func (r *RepositoryService) SetProcessDefinitionCategory(processDefinitionID, category string) error {
	if def, ok := r.processDefs[processDefinitionID]; ok {
		def.Category = category
		return nil
	}
	return errors.New("process definition not found")
}

func (r *RepositoryService) GetDeploymentResourceNames(deploymentID string) ([]string, error) {
	if dep, ok := r.deployments[deploymentID]; ok {
		return dep.Resources, nil
	}
	return nil, errors.New("deployment not found")
}

func (r *RepositoryService) GetResourceAsBytes(deploymentID, resourceName string) ([]byte, error) {
	if dep, ok := r.deployments[deploymentID]; ok {
		// 这里做成模拟：假设 resourceName 就是字符串内容
		for _, res := range dep.Resources {
			if res == resourceName {
				return []byte("mock resource: " + res), nil
			}
		}
		return nil, errors.New("resource not found")
	}
	return nil, errors.New("deployment not found")
}

func (r *RepositoryService) NewModel(name, key string) *Model {
	m := &Model{
		ID:   uuid.New().String(),
		Name: name,
		Key:  key,
	}
	r.models[m.ID] = m
	return m
}

func (r *RepositoryService) SaveModel(m *Model) {
	r.models[m.ID] = m
}

func (r *RepositoryService) GetModel(modelID string) (*Model, error) {
	if m, ok := r.models[modelID]; ok {
		return m, nil
	}
	return nil, errors.New("model not found")
}

func (r *RepositoryService) DeleteModel(modelID string) error {
	if _, ok := r.models[modelID]; !ok {
		return errors.New("model not found")
	}
	delete(r.models, modelID)
	return nil
}

func (r *RepositoryService) ListModels() []*Model {
	var res []*Model
	for _, m := range r.models {
		res = append(res, m)
	}
	return res
}

func (r *RepositoryService) AddModelEditorSource(modelID string, bytes []byte) error {
	if m, ok := r.models[modelID]; ok {
		m.Editor = bytes
		return nil
	}
	return errors.New("model not found")
}

func (r *RepositoryService) GetModelEditorSource(modelID string) ([]byte, error) {
	if m, ok := r.models[modelID]; ok {
		return m.Editor, nil
	}
	return nil, errors.New("model not found")
}
