package wiki

type Wiki interface {
	GetNodes() ([]*Node, error)
	FindNodes(filter NodeFilter) ([]*Node, error)
	GetNode(id string) (*Node, error)
	FindNode(filter NodeFilter) (*Node, error)
	Refresh() error
}
