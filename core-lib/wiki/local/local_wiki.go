package local

import (
	"fmt"

	"github.com/3rd/core/core-lib/fs"
	"github.com/3rd/core/core-lib/wiki"
)

type LocalWikiConfig struct {
	Root  string
	Parse bool
}

type LocalWiki struct {
	config LocalWikiConfig
	nodes  []*LocalNode
}

func NewLocalWiki(config LocalWikiConfig) (*LocalWiki, error) {
	wiki := &LocalWiki{
		config: config,
	}
	err := wiki.Refresh()
	if err != nil {
		return nil, err
	}
	return wiki, nil
}

func (w *LocalWiki) GetNodes() ([]*LocalNode, error) {
	var nodes []*LocalNode
	nodes = append(nodes, w.nodes...)
	return nodes, nil
}

func (w *LocalWiki) FindNodes(filter wiki.NodeFilter) ([]*LocalNode, error) {
	var nodes []*LocalNode
	for _, node := range w.nodes {
		if filter(node) {
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

func (w *LocalWiki) GetNode(id wiki.NodeID) (*LocalNode, error) {
	for _, node := range w.nodes {
		if node.GetID() == id {
			return node, nil
		}
	}
	return nil, nil
}

func (w *LocalWiki) FindNode(filter wiki.NodeFilter) (*LocalNode, error) {
	for _, node := range w.nodes {
		if filter(node) {
			return node, nil
		}
	}
	return nil, nil
}

func (w *LocalWiki) Refresh() error {
	// walk root
	files, err := fs.WalkFiles(w.config.Root, nil)
	if err != nil {
		return err
	}

	// collect nodes, fail on the first collision
	nodes := []*LocalNode{}
	nodemap := map[wiki.NodeID]*LocalNode{}
	for _, file := range files {
		node, err := NewLocalNode(file.GetPath())
		if err != nil {
			return err
		}
		if w.config.Parse {
			err = node.Parse()
			if err != nil {
				return err
			}
		}
		nodes = append(nodes, node)
		if previousNode, ok := nodemap[node.GetID()]; ok {
			return fmt.Errorf("failed to open wiki, found colliding nodes with id: %v\nA: %v\nB: %v", node.GetName(), previousNode, node)
		}
		nodemap[node.GetID()] = node
	}

	w.nodes = nodes
	return nil
}
