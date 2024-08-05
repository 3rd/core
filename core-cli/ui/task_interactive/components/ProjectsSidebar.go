package components

import (
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
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
	b.FillStyle(ui.Style{Background: theme.PROJECT_SIDEBAR_BG, Foreground: theme.PROJECT_SIDEBAR_FG})

	for i := 0; i < c.Height; i++ {
		projectIndex := i + c.AppState.ProjectScrollOffset
		if projectIndex >= len(c.AppState.Nodes) {
			break
		}
		project := c.AppState.Nodes[projectIndex]

		style := ui.Style{Background: theme.PROJECT_SIDEBAR_BG, Foreground: theme.PROJECT_SIDEBAR_FG}
		if projectIndex == c.AppState.ProjectSelectedIndex {
			style.Background = theme.PROJECT_SIDEBAR_SELECTED_BG
			style.Foreground = theme.PROJECT_SIDEBAR_SELECTED_FG
		}

		lineBuffer := ui.Buffer{}
		lineBuffer.Resize(c.Width, 1)
		lineBuffer.FillStyle(style)
		b.DrawBuffer(0, i, lineBuffer)

		projectText := project.GetName()
		projectText = strings.TrimPrefix(projectText, "project-")
		b.Text(1, i, projectText, style)
	}

	return b
}
