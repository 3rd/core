package components

import (
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"
	"strings"

	ui "github.com/3rd/go-futui"
)

type HistoryView struct {
	ui.Component
	AppState *state.AppState
	Width    int
	Height   int
}

func (c *HistoryView) Render() ui.Buffer {
	b := ui.Buffer{}
	b.Resize(c.Width, c.Height)
	b.FillStyle(ui.Style{Background: theme.BG, Foreground: theme.FG})

	yOffset := 0
	historyEntries := c.AppState.GetHistoryEntries()

	for i := c.AppState.HistoryEntryOffset; i < len(historyEntries) && yOffset < c.Height; i++ {
		entry := historyEntries[i]
		dateStr := entry.Date.Format("2006-01-02")
		b.Text(1, yOffset, dateStr, ui.Style{Foreground: theme.HISTORY_DATE_FG})
		yOffset++

		for _, task := range entry.Tasks {
			if yOffset >= c.Height {
				break
			}

			projectName := ""
			if task.Node != nil {
				projectName = task.Node.GetName()
				projectName = strings.TrimPrefix(projectName, "project-")
			}

			b.Text(0, yOffset, fmt.Sprintf("  ▕%s: ", projectName), ui.Style{Foreground: theme.HISTORY_PROJECT_FG})
			b.Text(7+len(projectName), yOffset, task.Text, ui.Style{Foreground: theme.HISTORY_TASK_FG})
			yOffset++
		}

		if yOffset < c.Height {
			yOffset++
		}
	}

	return b
}
