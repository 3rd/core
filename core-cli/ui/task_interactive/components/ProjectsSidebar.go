package components

import (
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"
	"strings"

	ui "github.com/3rd/go-futui"
)

type ProjectSidebar struct {
	ui.Component
	AppState *state.AppState
	Width    int
	Height   int
}

func (c *ProjectSidebar) Render() ui.Buffer {
	b := ui.Buffer{}
	b.Resize(c.Width, c.Height)

	var entries []string
	longestTextLength := 0
	for i := 0; i < c.Height; i++ {
		projectIndex := i + c.AppState.ProjectScrollOffset
		if projectIndex >= len(c.AppState.Nodes) {
			break
		}
		project := c.AppState.Nodes[projectIndex]

		projectName := project.GetName()
		projectName = strings.TrimPrefix(projectName, "project-")
		entry := fmt.Sprintf("%s (%d)", projectName, len(project.GetTasks()))

		if len(entry) > longestTextLength {
			longestTextLength = len(entry)
		}

		entries = append(entries, entry)
	}

	width := longestTextLength + 2
	b.Resize(width, c.Height)
	b.FillStyle(ui.Style{Background: theme.PROJECT_SIDEBAR_BG, Foreground: theme.PROJECT_SIDEBAR_FG})

	for i := 0; i < c.Height; i++ {
		if i >= len(entries) {
			break
		}
		entry := entries[i]
		style := ui.Style{Background: theme.PROJECT_SIDEBAR_BG, Foreground: theme.PROJECT_SIDEBAR_FG}
		if i == c.AppState.ProjectSelectedIndex {
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
