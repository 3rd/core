package local

import (
	"sort"
	"sync"

	"github.com/3rd/core/core-lib/fs"
	"github.com/3rd/core/core-lib/wiki"
)

type PARSE_MODE string

const (
	PARSE_MODE_NONE PARSE_MODE = "none"
	PARSE_MODE_FULL PARSE_MODE = "full"
	PARSE_MODE_META PARSE_MODE = "meta"
)

type LocalWikiConfig struct {
	Root            string
	Parse           PARSE_MODE
	SkipInitialLoad bool
}

type LocalWiki struct {
	config LocalWikiConfig
	nodes  []*LocalNode
}

func NewLocalWiki(config LocalWikiConfig) (*LocalWiki, error) {
	wiki := LocalWiki{
		config: config,
	}
	if !config.SkipInitialLoad {
		err := wiki.Reload()
		if err != nil {
			return nil, err
		}
	}
	return &wiki, nil
}

func (w *LocalWiki) GetNodes() ([]*LocalNode, error) {
	return w.nodes, nil
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

func (w *LocalWiki) GetNode(id string) (*LocalNode, error) {
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

func (w *LocalWiki) Reload() error {
	// walk root
	files, err := fs.WalkFiles(w.config.Root, nil)
	if err != nil {
		return err
	}

	// collect nodes, fail on the first collision
	nodes := []*LocalNode{}
	var wg sync.WaitGroup
	mutex := sync.Mutex{}
	wg.Add(len(files))
	for _, file := range files {
		go func(file fs.File) {
			defer wg.Done()
			node, err := NewLocalNode(file.GetPath())
			if err != nil {
				return
			}
			if w.config.Parse != PARSE_MODE_NONE {
				err = node.Parse(w.config.Parse)
				if err != nil {
					return
				}
			}
			mutex.Lock()
			nodes = append(nodes, node)
			mutex.Unlock()
		}(file)
	}
	wg.Wait()

	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].GetName() < nodes[j].GetName()
	})

	w.nodes = nodes
	return nil
}
