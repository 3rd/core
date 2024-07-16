package wiki

type Node interface {
	GetID() string
	GetName() string
	GetContent() (string, error)
	GetTasks() []*Task
}

type NodeFilter func(node Node) bool
