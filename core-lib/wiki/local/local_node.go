package local

import (
	"github.com/3rd/core/core-lib/fs"
	"github.com/3rd/core/core-lib/wiki"
	"github.com/3rd/syslang/go-syslang/pkg/syslang"
)

type LocalNode struct {
	fs.File
	document *syslang.Document
}

func NewLocalNode(path string) (*LocalNode, error) {
	file, err := fs.NewFile(path)
	if err != nil {
		return nil, err
	}

	node := LocalNode{*file, nil}
	return &node, nil
}

func (n *LocalNode) GetID() wiki.NodeID {
	return wiki.NodeID(n.GetName())
}

func (n *LocalNode) GetName() string {
	if n.document != nil {
		title := n.document.GetTitle()
		if title != "" {
			return title
		}
	}
	return n.File.GetName()
}

func (n *LocalNode) GetContent() (string, error) {
	return n.Text()
}

func (n *LocalNode) IsParsed() bool {
	return n.document != nil
}

func (n *LocalNode) Parse() error {
	text, err := n.Text()
	if err != nil {
		return err
	}
	n.document, err = syslang.NewDocument(text)
	if err != nil {
		return err
	}
	return nil
}

func (n *LocalNode) Refresh() error {
	if !n.IsParsed() {
		return nil
	}
	return n.Parse()
}
