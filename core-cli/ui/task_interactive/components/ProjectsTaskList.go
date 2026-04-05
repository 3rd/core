package components

import (
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"

	"github.com/3rd/core/core-lib/wiki"
	ui "github.com/3rd/go-futui"
)

type ProjectTaskList struct {
	ui.Component
	AppState *state.AppState
	Width    int
	Height   int
}

func (c *ProjectTaskList) Render() ui.Buffer {
	b := ui.Buffer{}
	b.Resize(c.Width, c.Height)
	b.FillStyle(theme.PROJECTS_TASK_STYLE)

	tasks := c.AppState.GetCurrentProjectTasks()
	for i := 0; i < c.Height; i++ {
		taskIndex := i + c.AppState.ProjectsTaskScrollOffset
		if taskIndex >= len(tasks) {
			break
		}
		task := tasks[taskIndex]

		style := theme.ProjectsTaskStyle(task.Status == wiki.TASK_STATUS_ACTIVE, taskIndex == c.AppState.ProjectsTaskSelectedIndex)

		// line
		lineBuffer := ui.Buffer{}
		lineBuffer.Resize(c.Width, 1)
		lineBuffer.FillStyle(style)
		b.DrawBuffer(0, i, lineBuffer)

		// marker
		statusMarker := " 󰄱"
		if task.Status == wiki.TASK_STATUS_ACTIVE {
			statusMarker = " ➡"
		}
		b.Text(0, i, statusMarker, style)

		// task text
		taskText := task.Text
		if task.Status == wiki.TASK_STATUS_DONE {
			taskText = task.Text
		}
		b.Text(3, i, taskText, style)
	}

	return b
}
