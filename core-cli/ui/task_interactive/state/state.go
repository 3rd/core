package state

import (
	"sort"
	"time"

	"github.com/3rd/core/core-lib/wiki"
)

type APP_TAB string
type APP_ACTIVE_MODE string

const (
	APP_TAB_ACTIVE  APP_TAB = ""
	APP_TAB_HISTORY APP_TAB = "history"

	APP_ACTIVE_MODE_DEFAULT APP_ACTIVE_MODE = ""
	APP_ACTIVE_MODE_EDITOR  APP_ACTIVE_MODE = "editor"
)

type HistoryEntry struct {
	Date  time.Time
	Tasks []*wiki.Task
}

type AppState struct {
	CurrentTab           APP_TAB
	Tasks                []*wiki.Task
	ActiveTasks          []*wiki.Task
	LongestProjectLength int
	ActiveMode           APP_ACTIVE_MODE
	ActiveSelectedIndex  int
	ActiveScrollOffset   int
	HistoryEntryOffset   int
}

func (app *AppState) GetLongestTaskLength() int {
	max := 0
	for _, task := range app.ActiveTasks {
		if len(task.Text) > max {
			max = len(task.Text)
		}
	}
	return max
}

func (app *AppState) GetDoneTasksCount() int {
	count := 0
	for _, task := range app.ActiveTasks {
		if task.IsDone() {
			count++
		}
	}
	return count
}

func (app *AppState) GetNotDoneTasksCount() int {
	count := 0
	for _, task := range app.ActiveTasks {
		if !task.IsDone() {
			count++
		}
	}
	return count
}

func (app *AppState) GetHistoryEntries() []HistoryEntry {
	entries := map[time.Time][]*wiki.Task{}

	for _, task := range app.Tasks {
		if !task.IsDone() {
			continue
		}

		lastSession := task.GetLastSession()
		if lastSession != nil {
			date := lastSession.Start
			lastSessionDayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
			entries[lastSessionDayStart] = append(entries[lastSessionDayStart], task)
		}
		// TODO: recurrent completions
	}

	historyEntries := []HistoryEntry{}
	for date, tasks := range entries {
		historyEntries = append(historyEntries, HistoryEntry{
			Date:  date,
			Tasks: tasks,
		})
	}

	sort.Slice(historyEntries, func(i, j int) bool {
		return historyEntries[j].Date.Before(historyEntries[i].Date)
	})

	return historyEntries
}
