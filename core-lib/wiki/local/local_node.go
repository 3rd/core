package local

import (
	"strings"

	"github.com/3rd/core/core-lib/fs"
	"github.com/3rd/core/core-lib/wiki"
	"github.com/3rd/syslang/go-syslang/pkg/syslang"
)

type LocalNode struct {
	fs.File
	document   *syslang.Document
	parsedMode PARSE_MODE
}

func NewLocalNode(path string) (*LocalNode, error) {
	file, err := fs.NewFile(path)
	if err != nil {
		return nil, err
	}

	node := LocalNode{*file, nil, PARSE_MODE_NONE}
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

func (n *LocalNode) Parse(mode PARSE_MODE) error {
	if mode == PARSE_MODE_NONE {
		panic("cannot parse with PARSE_MODE_NONE, you have a bug")
	}
	n.parsedMode = mode

	text, err := n.Text()
	if err != nil {
		return err
	}

	if mode == PARSE_MODE_META {
		if !strings.HasPrefix(text, "@meta") {
			return nil
		}
		endIndex := strings.Index(text, "@end")
		if endIndex == -1 {
			return nil
		}
		text = text[:endIndex+len("@end")]
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
	return n.Parse(n.parsedMode)
}

func (n *LocalNode) GetTasks() []*wiki.Task {
	syslangTasks := n.document.GetTasks()
	tasks := []*wiki.Task{}
	for _, syslangTask := range syslangTasks {
		sessions := []wiki.TaskSession{}
		for _, session := range syslangTask.Sessions {
			sessions = append(sessions, wiki.TaskSession{
				Start: session.Start,
				End:   session.End,
			})
		}

		var schedule *wiki.TaskSchedule
		if syslangTask.Schedule != nil {
			schedule = &wiki.TaskSchedule{
				Start:  syslangTask.Schedule.Start,
				End:    syslangTask.Schedule.End,
				Repeat: syslangTask.Schedule.Repeat,
			}
		}

		completions := []wiki.TaskCompletion{}
		for _, completion := range syslangTask.Completions {
			completions = append(completions, wiki.TaskCompletion{
				Timestamp: completion.Start,
			})
		}

		task := &wiki.Task{
			Parent:      nil,
			Children:    []*wiki.Task{},
			Sessions:    sessions,
			Schedule:    schedule,
			Text:        syslangTask.Title,
			LineNumber:  syslangTask.Line,
			Status:      wiki.TASK_STATUS_DEFAULT,
			Completions: completions,
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
