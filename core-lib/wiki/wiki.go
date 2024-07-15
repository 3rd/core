package wiki

type Wiki interface {
	GetNodes() ([]*Node, error)
	FindNodes(filter NodeFilter) ([]*Node, error)
	GetNode(id NodeID) (*Node, error)
	FindNode(filter NodeFilter) (*Node, error)
	Refresh() error
}
