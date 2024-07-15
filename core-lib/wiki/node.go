package wiki

type NodeID string

type Node interface {
	GetID() NodeID
	GetName() string
	GetContent() (string, error)
	GetTasks() []*Task
}

type NodeFilter func(node Node) bool
