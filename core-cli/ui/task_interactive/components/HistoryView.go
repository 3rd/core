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
	b.FillStyle(theme.HISTORY_STYLE)

	yOffset := 0
	historyEntries := c.AppState.GetHistoryEntries()

	// when offset is -1, set it to the bottom offset
	if c.AppState.HistoryEntryOffset == -1 {
		c.AppState.HistoryEntryOffset = len(historyEntries) - 1
	}

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
			b.Text(0, yOffset, fmt.Sprintf("  ▕%s: ", projectName), theme.HISTORY_PROJECT_STYLE)

			// task
			b.Text(7+len(projectName), yOffset, task.Text, theme.HISTORY_TASK_STYLE)

			// work time
			taskWorkTime := task.GetTotalSessionTimeForDate(entry.Date)
			dayWorkTime += taskWorkTime
			b.Text(8+len(projectName)+len(task.Text), yOffset, fmt.Sprintf("(%s)", taskWorkTime), theme.HISTORY_DURATION_STYLE)

			yOffset++
		}

		// date & work time
		dateStr := entry.Date.Format("2006-01-02")
		b.Text(1, dateYOffset, dateStr, theme.HISTORY_DATE_STYLE)
		b.Text(2+len(dateStr), dateYOffset, fmt.Sprintf("(%s)", dayWorkTime), theme.HISTORY_DURATION_STYLE)

		if yOffset < c.Height {
			yOffset++
		}
	}

	return b
}
