package taskinteractive

import (
	"core/ui/task_interactive/components"
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/3rd/core/core-lib/wiki"
	localWiki "github.com/3rd/core/core-lib/wiki/local"
	ui "github.com/3rd/go-futui"
	"github.com/gdamore/tcell/v2"
	"github.com/radovskyb/watcher"
)

type GetTasksResult struct {
	Tasks                []*wiki.Task
	LongestProjectLength int
}
type Providers struct {
	GetRoot  func() string
	GetTasks func() GetTasksResult
}

type App struct {
	ui.App
	state     state.AppState
	providers Providers
}

func (app *App) Setup() {
	app.loadTasks()

	// redraw ticker
	done := make(chan bool)
	ticker := time.NewTicker(time.Second / 2)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if app.state.Mode == state.APP_MODE_DEFAULT {
					app.Update()
				}
			}
		}
	}()

	// watcher
	w := watcher.New()
	w.FilterOps(watcher.Create, watcher.Move, watcher.Remove, watcher.Write)

	go func() {
		for {
			select {
			case <-w.Event:
				if app.state.Mode == state.APP_MODE_DEFAULT {
					app.loadTasks()
					app.Update()
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(app.providers.GetRoot()); err != nil {
		log.Fatalln(err)
	}
	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}
	go w.Start(time.Millisecond * 100)
}

func (app *App) loadTasks() {
	getTasksResult := app.providers.GetTasks()
	app.state.Tasks = getTasksResult.Tasks
	app.state.LongestProjectLength = getTasksResult.LongestProjectLength

	if app.state.SelectedIndex >= len(app.state.Tasks) {
		app.state.SelectedIndex = len(app.state.Tasks) - 1
	}
}

func (app *App) handleNavigateDown() {
	// select task
	i := app.state.SelectedIndex
	if i >= len(app.state.Tasks)-1 {
		return
	}
	i = i + 1
	app.state.SelectedIndex = i
	// scroll
	_, h := app.Screen.Size()
	if i >= h-2+app.state.ScrollOffset {
		app.state.ScrollOffset++
	}
	app.Update()
}

func (app *App) handleNavigateUp() {
	// select task
	i := app.state.SelectedIndex
	if i <= 0 {
		return
	}
	i = i - 1
	app.state.SelectedIndex = i
	// scroll
	if i < app.state.ScrollOffset {
		app.state.ScrollOffset--
	}
	app.Update()
}

func (app *App) handleEdit() {
	task := app.state.Tasks[app.state.SelectedIndex]
	node := task.Node.(*localWiki.LocalNode)
	if node == nil {
		return
	}

	initialMode := app.state.Mode
	app.state.Mode = state.APP_MODE_EDITOR

	app.Screen.Suspend()
	cmd := exec.Command("nvim", fmt.Sprintf("+%d", task.LineNumber+1), node.GetPath(), "+norm zz", "+norm zv")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	app.Screen.Resume()
	app.state.Mode = initialMode
	app.loadTasks()
	app.Update()
}

func (app *App) OnKeypress(ev tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'q':
			app.Quit()
		case 'j':
			app.handleNavigateDown()
		case 'k':
			app.handleNavigateUp()
		}
	case tcell.KeyCtrlC:
		app.Quit()
	case tcell.KeyEnter:
		app.handleEdit()
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
		SelectedIndex:        app.state.SelectedIndex,
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
