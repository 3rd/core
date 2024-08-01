package taskinteractive

import (
	"core/ui/task_interactive/components"
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/3rd/core/core-lib/wiki"
	localWiki "github.com/3rd/core/core-lib/wiki/local"
	ui "github.com/3rd/go-futui"
	"github.com/gdamore/tcell/v2"
	"github.com/radovskyb/watcher"
)

const INDENT = "  "

type GetTasksResult struct {
	Tasks                      []*wiki.Task
	ActiveTasks                []*wiki.Task
	LongestActiveProjectLength int
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
				if app.state.ActiveMode == state.APP_ACTIVE_MODE_DEFAULT {
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
				if app.state.ActiveMode == state.APP_ACTIVE_MODE_DEFAULT {
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
	go w.Start(time.Millisecond * 100)
}

func (app *App) loadTasks() {
	getTasksResult := app.providers.GetTasks()
	app.state.Tasks = getTasksResult.Tasks
	app.state.ActiveTasks = getTasksResult.ActiveTasks
	app.state.LongestProjectLength = getTasksResult.LongestActiveProjectLength

	// guard out of bounds
	if app.state.ActiveSelectedIndex >= len(app.state.ActiveTasks) {
		app.state.ActiveSelectedIndex = len(app.state.ActiveTasks) - 1
	}
}

func (app *App) handleNavigateDown() {
	// select task
	i := app.state.ActiveSelectedIndex
	if i >= len(app.state.ActiveTasks)-1 {
		return
	}
	i = i + 1
	app.state.ActiveSelectedIndex = i
	// scroll
	_, h := app.Screen.Size()
	if i >= h-2+app.state.ActiveScrollOffset {
		app.state.ActiveScrollOffset++
	}
	app.Update()
}

func (app *App) handleNavigateUp() {
	// select task
	i := app.state.ActiveSelectedIndex
	if i <= 0 {
		return
	}
	i = i - 1
	app.state.ActiveSelectedIndex = i
	// scroll
	if i < app.state.ActiveScrollOffset {
		app.state.ActiveScrollOffset--
	}
	app.Update()
}

func (app *App) handleHistoryScrollDown() {
	historyEntries := app.state.GetHistoryEntries()
	if app.state.HistoryEntryOffset < len(historyEntries)-1 {
		app.state.HistoryEntryOffset++
		app.Update()
	}
}

func (app *App) handleHistoryScrollUp() {
	if app.state.HistoryEntryOffset > 0 {
		app.state.HistoryEntryOffset--
		app.Update()
	}
}

func (app *App) handleEdit() {
	task := app.state.ActiveTasks[app.state.ActiveSelectedIndex]
	node := task.Node.(*localWiki.LocalNode)
	if node == nil {
		return
	}

	initialMode := app.state.ActiveMode
	app.state.ActiveMode = state.APP_ACTIVE_MODE_EDITOR

	app.Screen.Suspend()
	cmd := exec.Command("nvim", fmt.Sprintf("+%d", task.LineNumber+1), node.GetPath(), "+norm zz", "+norm zv")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	app.Screen.Resume()
	app.state.ActiveMode = initialMode
	app.loadTasks()
	app.Update()
}

func getIndentLevel(task *wiki.Task) int {
	indentLevel := 0
	lineText := task.LineText
	for strings.HasPrefix(lineText, INDENT) {
		indentLevel++
		lineText = lineText[len(INDENT):]
	}
	return indentLevel
}

func (app *App) handleToggleInProgress() {
	task := app.state.ActiveTasks[app.state.ActiveSelectedIndex]
	node := task.Node.(*localWiki.LocalNode)
	now := time.Now()
	text, err := node.Text()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(text, "\n")

	if task.IsInProgress() {
		// end current session
		lastWorkSession := task.GetLastSession()
		st := fmt.Sprintf("Session: %04d.%02d.%02d %02d:%02d-%02d:%02d", lastWorkSession.Start.Year(), lastWorkSession.Start.Month(), lastWorkSession.Start.Day(), lastWorkSession.Start.Hour(), lastWorkSession.Start.Minute(), now.Hour(), now.Minute())
		for i := 0; i <= getIndentLevel(task); i++ {
			st = INDENT + st
		}
		lines[lastWorkSession.LineNumber] = st
	} else {
		// start new session
		st := fmt.Sprintf("Session: %04d.%02d.%02d %02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
		for i := 0; i <= getIndentLevel(task); i++ {
			st = INDENT + st
		}
		i := task.LineNumber + 1
		lastWorkSession := task.GetLastSession()
		if lastWorkSession != nil {
			i = lastWorkSession.LineNumber + 1
		} else if task.Schedule != nil {
			i = task.Schedule.LineNumber + 1
		}
		lines = append(lines, "")
		copy(lines[i+1:], lines[i:])
		lines[i] = st
	}

	out, err := os.Create(node.GetPath())
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = out.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		panic(err)
	}

	app.loadTasks()
	app.Update()
}

func (app *App) handleToggleDone() {
	task := app.state.ActiveTasks[app.state.ActiveSelectedIndex]
	node := task.Node.(*localWiki.LocalNode)
	now := time.Now()
	text, err := node.Text()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(text, "\n")

	// recurring tasks
	if task.Schedule != nil && task.Schedule.Repeat != "" {
		completion := task.GetCompletionForDate(time.Now())

		// remove completion
		if completion != nil {
			lines = append(lines[:completion.LineNumber], lines[completion.LineNumber+1:]...)
		} else {
			// add completion
			now := time.Now()
			st := fmt.Sprintf("Done: %04d.%02d.%02d %02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
			for i := 0; i <= getIndentLevel(task); i++ {
				st = INDENT + st
			}
			i := task.Schedule.LineNumber + 1
			lastCompletion := task.GetLastCompletion()
			if lastCompletion != nil {
				i = lastCompletion.LineNumber + 1
			} else {
				lastSession := task.GetLastSession()
				if lastSession != nil {
					i = lastSession.LineNumber + 1
				}
			}
			lines = append(lines, "")
			copy(lines[i+1:], lines[i:])
			lines[i] = st

			// session in-progress task -> end current session
			if task.IsInProgress() {
				now := time.Now()
				// end current session
				last := *task.GetLastSession()
				sessionText := fmt.Sprintf("Session: %04d.%02d.%02d %02d:%02d-%02d:%02d", last.Start.Year(), last.Start.Month(), last.Start.Day(), last.Start.Hour(), last.Start.Minute(), now.Hour(), now.Minute())
				for i := 0; i <= getIndentLevel(task); i++ {
					sessionText = INDENT + sessionText
				}
				lines[last.LineNumber] = sessionText
			}
		}
	} else {
		// non-recurring tasks

		// marker (no schedule): [x] -> [-] or [-] -> [x]
		if task.Schedule == nil {
			if task.Status == wiki.TASK_STATUS_DONE {
				lines[task.LineNumber] = strings.ReplaceAll(lines[task.LineNumber], "[x]", "[-]")
			} else {
				lines[task.LineNumber] = strings.ReplaceAll(lines[task.LineNumber], "[-]", "[x]")
			}
		}

		// marker (scheduled): [x | -] <-> [ ]
		if task.Schedule != nil {
			switch task.Status {
			case wiki.TASK_STATUS_DONE:
				lines[task.LineNumber] = strings.ReplaceAll(lines[task.LineNumber], "[x]", "[ ]")
			case wiki.TASK_STATUS_DEFAULT:
				lines[task.LineNumber] = strings.ReplaceAll(lines[task.LineNumber], "[ ]", "[x]")
			case wiki.TASK_STATUS_ACTIVE:
				lines[task.LineNumber] = strings.ReplaceAll(lines[task.LineNumber], "[-]", "[x]")
			}
		}

		// current task
		if task.IsInProgress() {
			// end current session
			lastWorkSession := task.GetLastSession()
			sessionText := fmt.Sprintf("Session: %04d.%02d.%02d %02d:%02d-%02d:%02d", lastWorkSession.Start.Year(), lastWorkSession.Start.Month(), lastWorkSession.Start.Day(), lastWorkSession.Start.Hour(), lastWorkSession.Start.Minute(), now.Hour(), now.Minute())
			for i := 0; i <= getIndentLevel(task); i++ {
				sessionText = INDENT + sessionText
			}
			lines[lastWorkSession.LineNumber] = sessionText
		}

		// inactive task -> insert empty work session
		if !task.IsInProgress() && len(task.Sessions) == 0 && task.Status != wiki.TASK_STATUS_DONE {
			now := time.Now()
			st := fmt.Sprintf("Session: %04d.%02d.%02d %02d:%02d-%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Hour(), now.Minute())
			for i := 0; i <= getIndentLevel(task); i++ {
				st = INDENT + st
			}
			i := task.LineNumber + 1
			lines = append(lines, "")
			copy(lines[i+1:], lines[i:])
			lines[i] = st
		}
	}

	out, err := os.Create(node.GetPath())
	if err != nil {
		panic(err)
	}
	defer out.Close()
	out.WriteString(strings.Join(lines, "\n"))

	app.loadTasks()
	app.Update()
}

func (app *App) handleDeactivateTask() {
	if len(app.state.ActiveTasks) == 0 {
		return
	}

	task := app.state.ActiveTasks[app.state.ActiveSelectedIndex]
	node := task.Node.(*localWiki.LocalNode)
	text, _ := node.Text()
	lines := strings.Split(string(text), "\n")

	updatedLineText := strings.Replace(task.LineText, "[-]", "[ ]", 1)
	lines[task.LineNumber] = updatedLineText

	out, err := os.Create(node.GetPath())
	if err != nil {
		panic(err)
	}
	defer out.Close()
	out.WriteString(strings.Join(lines, "\n"))

	app.Update()
}

func (app *App) OnKeypress(ev tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'q':
			app.Quit()
		case 'j':
			if app.state.CurrentTab == state.APP_TAB_ACTIVE {
				app.handleNavigateDown()
			}
			if app.state.CurrentTab == state.APP_TAB_HISTORY {
				app.handleHistoryScrollDown()
			}
		case 'k':
			if app.state.CurrentTab == state.APP_TAB_ACTIVE {
				app.handleNavigateUp()
			}
			if app.state.CurrentTab == state.APP_TAB_HISTORY {
				app.handleHistoryScrollUp()
			}
		case ' ':
			if app.state.CurrentTab == state.APP_TAB_ACTIVE {
				app.handleToggleInProgress()
			}
		case '1':
			app.state.CurrentTab = state.APP_TAB_ACTIVE
			app.Render()
		case '2':
			app.state.CurrentTab = state.APP_TAB_HISTORY
			app.Render()
		}
	case tcell.KeyCtrlC:
		if app.state.CurrentTab == state.APP_TAB_ACTIVE {
			app.handleDeactivateTask()
		}
	case tcell.KeyEnter:
		if app.state.CurrentTab == state.APP_TAB_ACTIVE {
			app.handleEdit()
		}
	case tcell.KeyCtrlSpace:
		if app.state.CurrentTab == state.APP_TAB_ACTIVE {
			app.handleToggleDone()
		}
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
	headerBuffer := b.DrawComponent(0, 0, &header)

	// active tab
	if app.state.CurrentTab == state.APP_TAB_ACTIVE {
		taskList := components.TaskList{
			Tasks:                app.state.ActiveTasks,
			Width:                app.Width(),
			LongestProjectLength: app.state.LongestProjectLength,
			SelectedIndex:        app.state.ActiveSelectedIndex,
		}
		b.DrawComponent(0, headerBuffer.Height(), &taskList)

		// guard max scroll
		maxScroll := len(app.state.ActiveTasks) - app.Height() + headerBuffer.Height()
		if maxScroll < 0 {
			maxScroll = 0
		}
		if app.state.ActiveScrollOffset > maxScroll {
			app.state.ActiveScrollOffset = maxScroll
		}
	}

	// history tab
	if app.state.CurrentTab == state.APP_TAB_HISTORY {
		historyView := components.HistoryView{
			AppState: &app.state,
			Width:    app.Width(),
			Height:   app.Height() - headerBuffer.Height(),
		}
		b.DrawComponent(0, headerBuffer.Height(), &historyView)
	}

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
