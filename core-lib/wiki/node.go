package wiki

type NodeID string

type Node interface {
	GetID() NodeID
	GetName() string
	GetContent() (string, error)
}

type NodeFilter func(node Node) bool
