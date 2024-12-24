package flowgo

type RuntimeService struct {
}

func NewRuntimeService() *RuntimeService {
	return &RuntimeService{}
}

func StartProcessInstanceByID(instanceID uint64) error {
	return nil
}
