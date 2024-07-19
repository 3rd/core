package taskinteractive

import (
	"core/ui/task_interactive/components"
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"

	"github.com/3rd/core/core-lib/wiki"
	ui "github.com/3rd/go-futui"
	"github.com/gdamore/tcell/v2"
)

type App struct {
	ui.App
	state state.AppState
}

func (app *App) Setup() {
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
		State: &app.state,
		Width: app.Width(),
	}
	b.DrawComponent(0, 0, &header)

	return b
}

func Run() {
	state := state.AppState{
		Mode: state.APP_MODE_DEFAULT,
		Tasks: []*wiki.Task{
			{
				Text: "task 1",
			},
			{
				Text: "task 2",
			},
		},
	}
	app := App{state: state}

	app.Run(&app)
}
