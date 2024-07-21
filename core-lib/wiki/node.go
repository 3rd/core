package wiki

type Node interface {
	GetID() string
	GetName() string
	GetMeta() map[string]string
	GetContent() (string, error)
	GetTasks() []*Task
}

type NodeFilter func(node Node) bool
