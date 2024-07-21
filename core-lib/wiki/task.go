package wiki

import (
	"time"
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

type TaskSession struct {
	Start      time.Time
	End        *time.Time
	LineNumber uint32
}

func (session TaskSession) Duration() time.Duration {
	if session.End == nil {
		return time.Since(session.Start)
	}
	return session.End.Sub(session.Start)
}

func (session TaskSession) IsInProgress(atTime ...time.Time) bool {
	if len(atTime) > 1 {
		panic("IsInProgress takes at most one argument")
	}
	if len(atTime) == 1 {
		return session.Start.Before(atTime[0]) && (session.End == nil || session.End.After(atTime[0]))
	}
	return session.End == nil
}

type TaskSchedule struct {
	Start      time.Time
	End        *time.Time
	Repeat     string
	LineNumber uint32
}

func (schedule TaskSchedule) Duration() time.Duration {
	if schedule.End == nil {
		return 0
	}
	return schedule.End.Sub(schedule.Start)
}

func (schedule TaskSchedule) IsInProgress(atTime ...time.Time) bool {
	if len(atTime) != 1 {
		panic("TaskSchedule.IsInProgress requires the atTime argument")
	}
	// between start and end
	if schedule.End != nil && schedule.Start.Before(atTime[0]) && schedule.End.After(atTime[0]) {
		return true
	}
	// same day, after start, no end
	if schedule.End == nil {
		targetDayStart := time.Date(atTime[0].Year(), atTime[0].Month(), atTime[0].Day(), 0, 0, 0, 0, atTime[0].Location())
		scheduleDayStart := time.Date(schedule.Start.Year(), schedule.Start.Month(), schedule.Start.Day(), 0, 0, 0, 0, schedule.Start.Location())
		if targetDayStart.Equal(scheduleDayStart) {
			return true
		}
	}
	return false
}

type TaskCompletion struct {
	Timestamp  time.Time
	LineNumber uint32
}

type Task struct {
	Node        Node
	Parent      *Task
	Children    []*Task
	Sessions    []TaskSession
	Schedule    *TaskSchedule
	Completions []TaskCompletion
	Text        string
	LineNumber  uint32
	LineText    string
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
	if t.Schedule == nil || t.Schedule.Repeat != "" {
		return false
	}
	return t.Schedule.IsInProgress(time.Now())
}

func (t *Task) GetTotalSessionTime() time.Duration {
	duration := time.Duration(0)
	for _, session := range t.Sessions {
		duration += session.Duration()
	}
	return duration
}

func (t *Task) GetTotalSessionTimeDeep() time.Duration {
	duration := t.GetTotalSessionTime()
	for _, child := range t.Children {
		duration += child.GetTotalSessionTimeDeep()
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

func (t *Task) GetLastSession() *TaskSession {
	if len(t.Sessions) == 0 {
		return nil
	}
	last := t.Sessions[len(t.Sessions)-1]
	return &last
}

func (t *Task) GetCompletionForDate(date time.Time) *TaskCompletion {
	for _, completion := range t.Completions {
		if completion.Timestamp.Year() == date.Year() && completion.Timestamp.Month() == date.Month() && completion.Timestamp.Day() == date.Day() {
			return &completion
		}
	}
	return nil
}

func (t *Task) HasCompletionForDate(date time.Time) bool {
	return t.GetCompletionForDate(date) != nil
}

func (t *Task) GetLastCompletion() *TaskCompletion {
	if len(t.Completions) == 0 {
		return nil
	}
	last := t.Completions[len(t.Completions)-1]
	return &last
}

func (t *Task) GetIcon() rune {
	if t.Status == TASK_STATUS_DONE {
		return '☑'
	}
	return '☐'
}
