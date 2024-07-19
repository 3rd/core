package state

import "github.com/3rd/core/core-lib/wiki"

type APP_MODE string

const (
	APP_MODE_DEFAULT APP_MODE = ""
	APP_MODE_EDITOR  APP_MODE = "editor"
	APP_MODE_FOCUS   APP_MODE = "focus"
)

type AppState struct {
	Mode                 APP_MODE
	Tasks                []*wiki.Task
	LongestProjectLength int
}

func (app *AppState) GetLongestTaskLength() int {
	max := 0
	for _, task := range app.Tasks {
		if len(task.Text) > max {
			max = len(task.Text)
		}
	}
	return max
}

func (app *AppState) GetDoneTasksCount() int {
	count := 0
	for _, task := range app.Tasks {
		if task.IsDone() {
			count++
		}
	}
	return count
}

func (app *AppState) GetNotDoneTasksCount() int {
	count := 0
	for _, task := range app.Tasks {
		if !task.IsDone() {
			count++
		}
	}
	return count
}
