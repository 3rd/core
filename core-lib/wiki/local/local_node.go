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

func (n *LocalNode) GetTasks() []*wiki.Task {
	syslangTasks := n.document.GetTasks()
	tasks := []*wiki.Task{}
	for _, syslangTask := range syslangTasks {
		task := &wiki.Task{
			Parent:     nil,
			Children:   []*wiki.Task{},
			Sessions:   syslangTask.Sessions,
			Schedule:   syslangTask.Schedule,
			Text:       syslangTask.Title,
			LineNumber: syslangTask.Line,
			Status:     wiki.TASK_STATUS_DEFAULT,
		}
		if syslangTask.Status == syslang.TaskStatusActive {
			task.Status = wiki.TASK_STATUS_ACTIVE
		}
		if syslangTask.Status == syslang.TaskStatusDone {
			task.Status = wiki.TASK_STATUS_DONE
		}
		if syslangTask.Status == syslang.TaskStatusCancelled {
			task.Status = wiki.TASK_STATUS_CANCELLED
		}
		tasks = append(tasks, task)
	}
	return tasks
}
