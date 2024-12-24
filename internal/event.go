package internal

type StartEvent struct {
}

func (e *StartEvent) Execute() error {
	return nil
}

type EndEvent struct {
}

func (e *EndEvent) Execute() error {
	return nil
}
