package internal

type Flow struct {
	Name string
}

func NewFlow(name string) *Flow {
	return &Flow{
		Name: name,
	}
}

func (f *Flow) Next(node *Node) {}
