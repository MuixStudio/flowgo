package flowgo

type OptionFunc func(engine *Engine)
type Engine struct {
	name              string
	repositoryService *RepositoryService
	runtimeService    *RuntimeService
	taskService       *TaskService
	historyService    *HistoryService
	managementService *ManagementService
	identityService   *IdentityService
	formService       *FormService
}

func NewEngine(opts ...OptionFunc) *Engine {
	repo := NewRepositoryService()
	runtime := NewRuntimeService()
	task := NewTaskService()
	history := NewHistoryService()
	management := NewManagementService()
	identity := NewIdentityService()
	form := NewFormService()

	return &Engine{
		repositoryService: repo,
		runtimeService:    runtime,
		taskService:       task,
		historyService:    history,
		managementService: management,
		identityService:   identity,
		formService:       form,
	}
}

func (e *Engine) GetRepositoryService() *RepositoryService {
	return e.repositoryService
}

func (e *Engine) GetRuntimeService() *RuntimeService {
	return e.runtimeService
}

func (e *Engine) GetTaskService() *TaskService {
	return e.taskService
}

func (e *Engine) GetHistoryService() *HistoryService {
	return e.historyService
}

func (e *Engine) GetManagementService() *ManagementService {
	return e.managementService
}

func (e *Engine) GetIdentityService() *IdentityService {
	return e.identityService
}

func (e *Engine) GetFormService() *FormService {
	return e.formService
}
