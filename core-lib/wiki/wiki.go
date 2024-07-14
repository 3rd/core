package wiki

type Wiki interface {
	GetNodes() ([]Node, error)
	FindNode(filter NodeFilter) (*Node, error)
	FindNodeByID(id NodeID) (*Node, error)
}
