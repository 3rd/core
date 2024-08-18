package components

import (
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"
	"strings"

	"github.com/3rd/core/core-lib/wiki"
	ui "github.com/3rd/go-futui"
)

type ProjectSidebar struct {
	ui.Component
	AppState     *state.AppState
	Width        int
	Height       int
	ScrollOffset int
}

func (c *ProjectSidebar) Render() ui.Buffer {
	b := ui.Buffer{}
	b.Resize(c.Width, c.Height)

	var entries []string
	longestTextLength := 0
	for i := c.ScrollOffset; i < len(c.AppState.Nodes) && i-c.ScrollOffset < c.Height; i++ {
		project := c.AppState.Nodes[i]
		tasks := []*wiki.Task{}
		for _, task := range project.GetTasks() {
			if task.Status == wiki.TASK_STATUS_DEFAULT || task.Status == wiki.TASK_STATUS_ACTIVE {
				tasks = append(tasks, task)
			}
		}

		projectName := project.GetName()
		projectName = strings.TrimPrefix(projectName, "project-")
		entry := fmt.Sprintf("%s (%d)", projectName, len(tasks))

		if len(entry) > longestTextLength {
			longestTextLength = len(entry)
		}

		entries = append(entries, entry)
	}

	width := longestTextLength + 2
	b.Resize(width, c.Height)
	b.FillStyle(ui.Style{Background: theme.PROJECT_SIDEBAR_BG, Foreground: theme.PROJECT_SIDEBAR_FG})

	for i := 0; i < c.Height && i < len(entries); i++ {
		entry := entries[i]
		style := ui.Style{Background: theme.PROJECT_SIDEBAR_BG, Foreground: theme.PROJECT_SIDEBAR_FG}
		if i+c.ScrollOffset == c.AppState.ProjectSelectedIndex {
			style.Background = theme.PROJECT_SIDEBAR_SELECTED_BG
			style.Foreground = theme.PROJECT_SIDEBAR_SELECTED_FG
		}

		lineBuffer := ui.Buffer{}
		lineBuffer.Resize(width, 1)
		lineBuffer.FillStyle(style)
		b.DrawBuffer(0, i, lineBuffer)

		text := entry
		b.Text(1, i, text, style)
	}

	c.Width = b.Width()
	c.Height = b.Height()

	return b
}
