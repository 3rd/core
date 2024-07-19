package taskinteractive

import (
	"core/ui/task_interactive/components"
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"log"
	"os"

	"github.com/3rd/core/core-lib/wiki"
	ui "github.com/3rd/go-futui"
	"github.com/gdamore/tcell/v2"
)

type GetTasksResult struct {
	Tasks                []*wiki.Task
	LongestProjectLength int
}
type Providers struct {
	GetTasks func() GetTasksResult
}

type App struct {
	ui.App
	state     state.AppState
	providers Providers
}

func (app *App) Setup() {
	getTasksResult := app.providers.GetTasks()
	app.state.Tasks = getTasksResult.Tasks
	app.state.LongestProjectLength = getTasksResult.LongestProjectLength
}

func (app *App) OnKeypress(ev tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyCtrlC:
		app.Quit()
	}
	app.Update()
}

func (app *App) OnResize() {
	app.Update()
}

func (app *App) Render() ui.Buffer {
	b := ui.Buffer{}
	b.Resize(app.Width(), app.Height())
	b.FillStyle(ui.Style{
		Background: theme.BG,
		Foreground: theme.FG,
	})

	header := components.Header{
		AppState: &app.state,
		Width:    app.Width(),
	}
	b.DrawComponent(0, 0, &header)

	taskList := components.TaskList{
		Tasks:                app.state.Tasks,
		Width:                app.Width(),
		LongestProjectLength: app.state.LongestProjectLength,
	}
	b.DrawComponent(0, 4, &taskList)

	return b
}

func Run(providers Providers) {
	// log to /tmp/core-task-interactive.log
	logFile, err := os.OpenFile("/tmp/core-task-interactive.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	app := App{providers: providers}
	app.Run(&app)
}
