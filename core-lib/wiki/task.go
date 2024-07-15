package wiki

import (
	"time"

	"github.com/3rd/syslang/go-syslang/pkg/syslang"
)

type TASK_STATUS string
type TASK_PATTERN string

const (
	TASK_STATUS_DEFAULT   TASK_STATUS = "default"
	TASK_STATUS_ACTIVE    TASK_STATUS = "active"
	TASK_STATUS_DONE      TASK_STATUS = "done"
	TASK_STATUS_CANCELLED TASK_STATUS = "cancelled"
)

const (
	TASK_RE_DEFAULT TASK_PATTERN = `^\s*\[ \] (?P<Text>.*)$`
	TASK_RE_ACTIVE  TASK_PATTERN = `^\s*\[-\] (?P<Text>.*)$`
	TASK_RE_DONE    TASK_PATTERN = `^\s*\[x\] (?P<Text>.*)$`
)

type Task struct {
	Parent      *Task
	Children    []*Task
	Sessions    []syslang.TaskSession
	Schedule    *syslang.TaskSchedule
	Text        string
	LineNumber  uint32
	IndentLevel uint32
	Status      TASK_STATUS
	Priority    uint32
	DetailLines []string
	Tags        []string
}

func (t *Task) IsDone() bool {
	return t.Status == TASK_STATUS_DONE
}

func (t *Task) IsInProgress() bool {
	// on-going working session
	for _, session := range t.Sessions {
		if session.End == nil {
			return true
		}
	}
	// scheduled and inside the scheduled interval now
	if t.Schedule == nil {
		return false
	}
	return t.Schedule.IsInProgress(time.Now())
}

func (t *Task) GetWorkTime() time.Duration {
	duration := time.Duration(0)
	for _, session := range t.Sessions {
		duration += session.Duration()
	}
	return duration
}

func (t *Task) GetTotalWorkTime() time.Duration {
	duration := time.Duration(0)
	for _, child := range t.Children {
		duration += child.GetTotalWorkTime()
	}
	return duration
}

func (t *Task) GetTotalPriority() uint32 {
	priority := t.Priority
	for _, child := range t.Children {
		priority += child.GetTotalPriority()
	}
	return priority
}

func (t *Task) GetIcon() rune {
	if t.Status == TASK_STATUS_DONE {
		return '☑'
	}
	return '☐'
}

func (t *Task) GetLastWorkSession() *syslang.TaskSession {
	if len(t.Sessions) == 0 {
		return nil
	}
	last := t.Sessions[len(t.Sessions)-1]
	return &last
}
