package wiki

import (
	"github.com/3rd/core/core-lib/fs"
	"github.com/3rd/syslang/go-syslang/pkg/syslang"
)

type WikiNode struct {
	fs.File
	document *syslang.Document
}

func NewWikiNode(path string) (*WikiNode, error) {
	file, err := fs.NewFile(path)
	if err != nil {
		return nil, err
	}

	node := WikiNode{*file, nil}
	return &node, nil
}

func (w *WikiNode) GetID() string {
	if w.document != nil {
		title := w.document.GetTitle()
		if title != "" {
			return title
		}
	}
	return w.File.GetName()
}

func (w *WikiNode) IsParsed() bool {
	return w.document != nil
}

func (w *WikiNode) Parse() error {
	text, err := w.Text()
	if err != nil {
		return err
	}
	w.document, err = syslang.NewDocument(text)
	if err != nil {
		return err
	}
	return nil
}

func (w *WikiNode) Refresh() error {
	if !w.IsParsed() {
		return nil
	}
	return w.Parse()
}
