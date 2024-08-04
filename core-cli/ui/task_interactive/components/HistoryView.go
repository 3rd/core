package components

import (
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"
	"strings"
	"time"

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
		dayWorkTime := time.Duration(0)

		// skip 1, will write the date and total work time at the end
		dateYOffset := yOffset
		yOffset++

		// tasks
		for _, task := range entry.Tasks {
			if yOffset >= c.Height {
				break
			}

			// project
			projectName := ""
			if task.Node != nil {
				projectName = task.Node.GetName()
				projectName = strings.TrimPrefix(projectName, "project-")
			}
			b.Text(0, yOffset, fmt.Sprintf("  ▕%s: ", projectName), ui.Style{Foreground: theme.HISTORY_PROJECT_FG})

			// task
			b.Text(7+len(projectName), yOffset, task.Text, ui.Style{Foreground: theme.HISTORY_TASK_FG})

			// work time
			taskWorkTime := task.GetTotalSessionTimeForDate(entry.Date)
			dayWorkTime += taskWorkTime
			b.Text(8+len(projectName)+len(task.Text), yOffset, fmt.Sprintf("(%s)", taskWorkTime), ui.Style{Foreground: theme.TASK_DONE_FG})

			yOffset++
		}

		// date & work time
		dateStr := entry.Date.Format("2006-01-02")
		b.Text(1, dateYOffset, dateStr, ui.Style{Foreground: theme.HISTORY_DATE_FG})
		b.Text(2+len(dateStr), dateYOffset, fmt.Sprintf("(%s)", dayWorkTime), ui.Style{Foreground: theme.TASK_DONE_FG})

		if yOffset < c.Height {
			yOffset++
		}
	}

	return b
}
