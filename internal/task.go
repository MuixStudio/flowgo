package internal

type EmailTask struct {
	To string
}

func (t *EmailTask) Execute() error {
	return nil
}

type UserTask struct {
}

func (t *UserTask) Execute() error {
	return nil
}
