package internal

type Node struct {
	Name string
}

func NewNode(name string) *Node {
	return &Node{
		Name: name,
	}
}
