package wiki

import (
	"fmt"

	"github.com/3rd/core/core-lib/fs"
)

type WikiConfig struct {
	Root  string
	Parse bool
}

type Wiki struct {
	Config WikiConfig
	Nodes  []WikiNode
}

func NewWiki(config WikiConfig) *Wiki {
	wiki := &Wiki{
		Config: config,
		Nodes:  make([]WikiNode, 0),
	}
	wiki.Refresh()
	return wiki
}

func (w *Wiki) FindNodeByID(id string) *WikiNode {
	for _, node := range w.Nodes {
		if node.GetID() == id {
			return &node
		}
	}
	return nil
}

func (w *Wiki) FindNodeByPath(path string) *WikiNode {
	for _, node := range w.Nodes {
		if node.GetPath() == path {
			return &node
		}
	}
	return nil
}

func (w *Wiki) Refresh() error {
	w.Nodes = make([]WikiNode, 0)

	// walk root
	files, err := fs.WalkFiles(w.Config.Root, nil)
	if err != nil {
		return err
	}

	// collect nodes, fail on the first collision
	nodes := []WikiNode{}
	nodemap := map[string]WikiNode{}
	for _, file := range files {
		node, err := NewWikiNode(file.GetPath())
		if err != nil {
			return err
		}
		if w.Config.Parse {
			err = node.Parse()
			if err != nil {
				return err
			}
		}
		nodes = append(nodes, *node)
		if previousNode, ok := nodemap[node.GetID()]; ok {
			return fmt.Errorf("failed to open wiki, found colliding nodes with id: %v\nA: %v\nB: %v", node.GetName(), previousNode, node)
		}
		nodemap[node.GetID()] = *node
	}

	w.Nodes = nodes
	return nil
}

func (w *Wiki) RefreshPath(path string) error {
	// refresh node with path if it exists
	for _, node := range w.Nodes {
		if node.GetPath() == path {
			err := node.Refresh()
			return err
		}
	}

	// create a new node for the new path
	node, err := NewWikiNode(path)
	if err != nil {
		return err
	}
	if w.Config.Parse {
		err = node.Parse()
		if err != nil {
			return err
		}
	}
	w.Nodes = append(w.Nodes, *node)
	return nil
}
