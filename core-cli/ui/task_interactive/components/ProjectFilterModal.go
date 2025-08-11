package components

import (
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"
	"strings"

	ui "github.com/3rd/go-futui"
)

type ProjectFilterModal struct {
	ui.Component
	AppState *state.AppState
	Width    int
	Height   int
}

func (c *ProjectFilterModal) Render() ui.Buffer {
	b := ui.Buffer{}
	modalState := &c.AppState.ProjectFilterModal

	if !modalState.IsVisible {
		return b
	}

	modalWidth := 60
	modalHeight := min(len(modalState.Projects)+5, c.Height-4)
	if modalWidth > c.Width-4 {
		modalWidth = c.Width - 4
	}

	b.Resize(modalWidth, modalHeight)
	b.FillStyle(ui.Style{
		Background: theme.HEADER_BG,
		Foreground: theme.FG,
	})

	// borders
	borderStyle := ui.Style{Foreground: theme.TASK_LABEL_FG}
	b.DrawCell(0, 0, '┌', borderStyle)
	for x := 1; x < modalWidth-1; x++ {
		b.DrawCell(x, 0, '─', borderStyle)
	}
	b.DrawCell(modalWidth-1, 0, '┐', borderStyle)
	for y := 1; y < modalHeight-1; y++ {
		b.DrawCell(0, y, '│', borderStyle)
		b.DrawCell(modalWidth-1, y, '│', borderStyle)
	}
	b.DrawCell(0, modalHeight-1, '└', borderStyle)
	for x := 1; x < modalWidth-1; x++ {
		b.DrawCell(x, modalHeight-1, '─', borderStyle)
	}
	b.DrawCell(modalWidth-1, modalHeight-1, '┘', borderStyle)

	// title
	title := " Filter Projects "
	titleX := (modalWidth - len(title)) / 2
	b.Text(titleX, 0, title, ui.Style{
		Foreground: theme.HEADER_FG,
		Bold:       true,
	})

	// project list
	listStartY := 2
	maxVisibleProjects := modalHeight - 5 // account for borders and help text inside

	// scroll offset
	scrollOffset := 0
	if modalState.CursorIndex >= maxVisibleProjects {
		scrollOffset = modalState.CursorIndex - maxVisibleProjects + 1
	}

	for i, project := range modalState.Projects {
		if i < scrollOffset {
			continue
		}
		if i-scrollOffset >= maxVisibleProjects {
			break
		}

		y := listStartY + i - scrollOffset

		// selection
		isSelected := i == modalState.CursorIndex
		checkbox := ""
		if project.IsEnabled {
			checkbox = ""
		}

		// project name
		projectName := project.ProjectName
		if projectName == "" {
			projectName = project.ProjectID
		}
		projectName = strings.TrimPrefix(projectName, "project-")

		// task count
		taskCountStr := fmt.Sprintf(" (%d)", project.TaskCount)

		// calculate available width for project name
		availableWidth := modalWidth - 2 - 2 - 1 - len(taskCountStr) - 2
		if len(projectName) > availableWidth {
			projectName = projectName[:availableWidth-3] + "..."
		}

		// line style
		lineStyle := ui.Style{Foreground: theme.FG}
		if !project.IsEnabled {
			lineStyle.Foreground = theme.TASK_DONE_FG
		}
		if isSelected {
			lineStyle.Foreground = theme.PROJECT_FG
			lineStyle.Bold = true
		}

		// item line
		var line string
		if checkbox != "" {
			line = fmt.Sprintf("%s %s%s", checkbox, projectName, taskCountStr)
		} else {
			line = fmt.Sprintf("  %s%s", projectName, taskCountStr)
		}
		b.Text(2, y, line, lineStyle)
	}

	// help text
	helpText := "j/k: move | space: toggle | q/esc: close"
	helpY := modalHeight - 2
	helpX := max((modalWidth-len(helpText))/2, 2)
	maxTextWidth := modalWidth - 4
	if len(helpText) > maxTextWidth {
		helpText = helpText[:maxTextWidth]
	}
	b.Text(helpX, helpY, helpText, ui.Style{
		Foreground: theme.TASK_LABEL_FG,
	})

	return b
}
