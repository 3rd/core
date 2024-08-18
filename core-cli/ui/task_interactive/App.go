package taskinteractive

import (
	"core/ui/task_interactive/components"
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/3rd/core/core-lib/wiki"
	localWiki "github.com/3rd/core/core-lib/wiki/local"
	ui "github.com/3rd/go-futui"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/radovskyb/watcher"
)

const INDENT = "  "

func getIndentLevel(task *wiki.Task) int {
	indentLevel := 0
	lineText := task.LineText
	for strings.HasPrefix(lineText, INDENT) {
		indentLevel++
		lineText = lineText[len(INDENT):]
	}
	return indentLevel
}

type GetTasksResult struct {
	Nodes                      []wiki.Node
	Tasks                      []*wiki.Task
	ActiveTasks                []*wiki.Task
	LongestActiveProjectLength int
	LongestProjectLength       int
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
	app.state.LongestActiveProjectLength = getTasksResult.LongestActiveProjectLength
	app.state.LongestProjectLength = getTasksResult.LongestProjectLength

	// sort projects
	sortedNodes := getTasksResult.Nodes
	sort.Slice(sortedNodes, func(i, j int) bool {
		aMeta := sortedNodes[i].GetMeta()
		bMeta := sortedNodes[j].GetMeta()

		if aMeta != nil && bMeta != nil {
			aName := sortedNodes[i].GetName()
			bName := sortedNodes[j].GetName()

			// projects with no workable tasks to the bottom
			aTasks := []*wiki.Task{}
			for _, task := range sortedNodes[i].GetTasks() {
				if task.Status == wiki.TASK_STATUS_ACTIVE || task.Status == wiki.TASK_STATUS_DEFAULT {
					aTasks = append(aTasks, task)
				}
			}
			bTasks := []*wiki.Task{}
			for _, task := range sortedNodes[j].GetTasks() {
				if task.Status == wiki.TASK_STATUS_ACTIVE || task.Status == wiki.TASK_STATUS_DEFAULT {
					bTasks = append(bTasks, task)
				}
			}
			if len(aTasks) == 0 && len(bTasks) != 0 {
				return false
			}
			if len(aTasks) != 0 && len(bTasks) == 0 {
				return true
			}

			// by priority
			aPriority, err := strconv.Atoi(aMeta["priority"])
			if err != nil {
				aPriority = 0
			}
			bPriority, err := strconv.Atoi(bMeta["priority"])
			if err != nil {
				bPriority = 0
			}
			// by name if they have the same non-zero priority
			if aPriority == bPriority && aPriority != 0 {
				return strings.Compare(aName, bName) < 0
			}
			// by raw priority if it's different
			if aPriority != bPriority {
				return aPriority > bPriority
			}

			// by task count
			if len(aTasks) != len(bTasks) {
				return len(aTasks) > len(bTasks)
			}

			// by name otherwise
			return strings.Compare(aName, bName) < 0
		}

		// without meta
		aName := strings.TrimPrefix(strings.ToLower(sortedNodes[i].GetName()), "project-")
		bName := strings.TrimPrefix(strings.ToLower(sortedNodes[j].GetName()), "project-")
		return strings.Compare(aName, bName) < 0
	})
	app.state.Nodes = sortedNodes

	// guard out of bounds
	if app.state.ActiveSelectedIndex >= len(app.state.ActiveTasks) {
		app.state.ActiveSelectedIndex = len(app.state.ActiveTasks) - 1
	}
}

func (app *App) showNotification(message string) {
	app.state.Notification = &state.Notification{Message: message}
	app.Update()

	time.AfterFunc(state.NOTIFICATION_DURATION, func() {
		if app.state.Notification != nil && app.state.Notification.Message == message {
			app.state.Notification = nil
			app.Update()
		}
	})
}

func (app *App) handleActiveNavigateDown() {
	if app.state.ActiveSelectedIndex >= len(app.state.ActiveTasks)-1 {
		return
	}
	app.state.ActiveSelectedIndex++
	app.adjustActiveScroll()
	app.Update()
}

func (app *App) handleActiveNavigateUp() {
	if app.state.ActiveSelectedIndex <= 0 {
		return
	}
	app.state.ActiveSelectedIndex--
	app.adjustActiveScroll()
	app.Update()
}

func (app *App) adjustActiveScroll() {
	maxVisibleTasks := app.Height() - app.state.HeaderHeight
	if maxVisibleTasks <= 0 {
		return // Prevent division by zero or negative values
	}
	maxScrollOffset := len(app.state.ActiveTasks) - maxVisibleTasks
	if maxScrollOffset < 0 {
		maxScrollOffset = 0
	}

	if app.state.ActiveSelectedIndex < app.state.ActiveScrollOffset {
		app.state.ActiveScrollOffset = app.state.ActiveSelectedIndex
	} else if app.state.ActiveSelectedIndex >= app.state.ActiveScrollOffset+maxVisibleTasks {
		app.state.ActiveScrollOffset = app.state.ActiveSelectedIndex - maxVisibleTasks + 1
	}

	if app.state.ActiveScrollOffset > maxScrollOffset {
		app.state.ActiveScrollOffset = maxScrollOffset
	}
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

func (app *App) handleActiveEdit() {
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

func (app *App) handleActiveToggleInProgress() {
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

		deletePreviousSession := false
		var previousSession *wiki.TaskSession
		for _, session := range task.Sessions {
			if session.End == nil {
				break
			}
			previousSession = &session
		}
		if previousSession != nil && previousSession.End != nil &&
			previousSession.Start.Year() == now.Year() &&
			previousSession.Start.Month() == now.Month() &&
			previousSession.Start.Day() == now.Day() &&
			previousSession.Start.Hour() == now.Hour() &&
			previousSession.Start.Minute() == now.Minute() {
			deletePreviousSession = true
		}

		lines[lastWorkSession.LineNumber] = st
		if deletePreviousSession {
			lines = append(lines[:previousSession.LineNumber], lines[previousSession.LineNumber+1:]...)
		}
	} else {
		// create new session
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

func (app *App) handleActiveToggleDone() {
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
				lines[task.LineNumber] = strings.Replace(lines[task.LineNumber], "[x]", "[-]", 1)
			} else {
				lines[task.LineNumber] = strings.Replace(lines[task.LineNumber], "[-]", "[x]", 1)
			}
		}

		// marker (scheduled): [x | -] <-> [ ]
		if task.Schedule != nil {
			switch task.Status {
			case wiki.TASK_STATUS_DONE:
				lines[task.LineNumber] = strings.Replace(lines[task.LineNumber], "[x]", "[ ]", 1)
			case wiki.TASK_STATUS_DEFAULT:
				lines[task.LineNumber] = strings.Replace(lines[task.LineNumber], "[ ]", "[x]", 1)
			case wiki.TASK_STATUS_ACTIVE:
				lines[task.LineNumber] = strings.Replace(lines[task.LineNumber], "[-]", "[x]", 1)
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

func (app *App) handleActiveDeactivateTask() {
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

func (app *App) handleProjectsNavigation(forward bool) {
	if len(app.state.Nodes) == 0 {
		return
	}
	if forward {
		if app.state.ProjectSelectedIndex < len(app.state.Nodes)-1 {
			app.state.ProjectSelectedIndex++
		}
	} else {
		if app.state.ProjectSelectedIndex > 0 {
			app.state.ProjectSelectedIndex--
		}
	}
	app.adjustProjectsSidebarScroll()
	app.state.ProjectsTaskSelectedIndex = 0
	app.state.ProjectsTaskScrollOffset = 0
	app.Update()
}

func (app *App) adjustProjectsSidebarScroll() {
	maxVisibleProjects := app.Height() - app.state.HeaderHeight
	if app.state.ProjectSelectedIndex < app.state.ProjectScrollOffset {
		app.state.ProjectScrollOffset = app.state.ProjectSelectedIndex
	} else if app.state.ProjectSelectedIndex >= app.state.ProjectScrollOffset+maxVisibleProjects {
		app.state.ProjectScrollOffset = app.state.ProjectSelectedIndex - maxVisibleProjects + 1
	}
}

func (app *App) handleProjectsTaskNavigation(down bool) {
	tasks := app.state.GetCurrentProjectTasks()
	if len(tasks) == 0 {
		return
	}
	if down {
		if app.state.ProjectsTaskSelectedIndex < len(tasks)-1 {
			app.state.ProjectsTaskSelectedIndex++
		}
	} else {
		if app.state.ProjectsTaskSelectedIndex > 0 {
			app.state.ProjectsTaskSelectedIndex--
		}
	}
	app.projectsAdjustTaskScroll()
}

func (app *App) handleProjectsEdit() {
	project := app.state.Nodes[app.state.ProjectSelectedIndex]
	if project == nil {
		return
	}

	editorArgs := []string{}

	tasks := app.state.GetCurrentProjectTasks()
	if len(tasks) > 0 {
		task := app.state.GetCurrentProjectTasks()[app.state.ProjectsTaskSelectedIndex]
		if task == nil {
			return
		}
		editorArgs = append(editorArgs, fmt.Sprintf("+%d", task.LineNumber+1))
	}

	initialMode := app.state.ActiveMode
	app.state.ActiveMode = state.APP_ACTIVE_MODE_EDITOR

	app.Screen.Suspend()
	editorArgs = append(editorArgs, project.(*localWiki.LocalNode).GetPath())
	editorArgs = append(editorArgs, "+norm zz")
	cmd := exec.Command("nvim", editorArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	app.Screen.Resume()
	app.state.ActiveMode = initialMode
	app.loadTasks()
	app.Update()
}

func (app *App) projectsAdjustSidebarScroll() {
	maxVisibleProjects := app.Height() - app.state.HeaderHeight - 1
	maxScrollOffset := len(app.state.Nodes) - maxVisibleProjects
	if maxScrollOffset < 0 {
		maxScrollOffset = 0
	}

	if app.state.ProjectSelectedIndex < app.state.ProjectScrollOffset {
		app.state.ProjectScrollOffset = app.state.ProjectSelectedIndex
	} else if app.state.ProjectSelectedIndex >= app.state.ProjectScrollOffset+maxVisibleProjects {
		app.state.ProjectScrollOffset = app.state.ProjectSelectedIndex - maxVisibleProjects + 1
	}

	if app.state.ProjectScrollOffset > maxScrollOffset {
		app.state.ProjectScrollOffset = maxScrollOffset
	}
}

func (app *App) projectsAdjustTaskScroll() {
	maxVisibleTasks := app.Height() - app.state.HeaderHeight - 1
	tasks := app.state.GetCurrentProjectTasks()
	maxScrollOffset := len(tasks) - maxVisibleTasks
	if maxScrollOffset < 0 {
		maxScrollOffset = 0
	}

	if app.state.ProjectsTaskSelectedIndex < app.state.ProjectsTaskScrollOffset {
		app.state.ProjectsTaskScrollOffset = app.state.ProjectsTaskSelectedIndex
	} else if app.state.ProjectsTaskSelectedIndex >= app.state.ProjectsTaskScrollOffset+maxVisibleTasks {
		app.state.ProjectsTaskScrollOffset = app.state.ProjectsTaskSelectedIndex - maxVisibleTasks + 1
	}

	if app.state.ProjectsTaskScrollOffset > maxScrollOffset {
		app.state.ProjectsTaskScrollOffset = maxScrollOffset
	}
}

func (app *App) handleProjectsToggleTask() {
	tasks := app.state.GetCurrentProjectTasks()
	if app.state.ProjectsTaskSelectedIndex < 0 || app.state.ProjectsTaskSelectedIndex > len(tasks) {
		return
	}

	task := tasks[app.state.ProjectsTaskSelectedIndex]
	if task == nil {
		return
	}
	node := task.Node.(*localWiki.LocalNode)
	text, _ := node.Text()
	lines := strings.Split(string(text), "\n")

	if task.Status == wiki.TASK_STATUS_ACTIVE {
		lines[task.LineNumber] = strings.Replace(lines[task.LineNumber], "[-]", "[ ]", 1)
	} else {
		lines[task.LineNumber] = strings.Replace(lines[task.LineNumber], "[ ]", "[-]", 1)
	}

	out, err := os.Create(node.GetPath())
	if err != nil {
		return
	}
	defer out.Close()
	out.WriteString(strings.Join(lines, "\n"))
}

func (app *App) handleNavigateTop() {
	switch app.state.CurrentTab {
	case state.APP_TAB_ACTIVE:
		app.state.ActiveScrollOffset = 0
		app.state.ActiveSelectedIndex = 0
	case state.APP_TAB_HISTORY:
		app.state.HistoryEntryOffset = 0
	case state.APP_TAB_PROJECTS:
		app.state.ProjectsTaskSelectedIndex = 0
		app.state.ProjectsTaskScrollOffset = 0
	}
	app.Update()
}

func (app *App) handleNavigateBottom() {
	switch app.state.CurrentTab {
	case state.APP_TAB_ACTIVE:
		if len(app.state.ActiveTasks) > 0 {
			app.state.ActiveSelectedIndex = len(app.state.ActiveTasks) - 1
			app.adjustActiveScroll()
		}
	case state.APP_TAB_HISTORY:
		app.state.HistoryEntryOffset = -1
	case state.APP_TAB_PROJECTS:
		tasks := app.state.GetCurrentProjectTasks()
		if len(tasks) > 0 {
			app.state.ProjectsTaskSelectedIndex = len(tasks) - 1
			app.projectsAdjustTaskScroll()
		}
	}
	app.Update()
}

func (app *App) handleYank() {
	var taskText string
	switch app.state.CurrentTab {
	case state.APP_TAB_ACTIVE:
		if app.state.ActiveSelectedIndex < len(app.state.ActiveTasks) {
			taskText = app.state.ActiveTasks[app.state.ActiveSelectedIndex].Text
		}
	case state.APP_TAB_PROJECTS:
		tasks := app.state.GetCurrentProjectTasks()
		if app.state.ProjectsTaskSelectedIndex < len(tasks) {
			taskText = tasks[app.state.ProjectsTaskSelectedIndex].Text
		}
	}

	if taskText != "" {
		err := clipboard.WriteAll(taskText)
		if err == nil {
			app.showNotification("Copied task to clipboard")
		} else {
			app.showNotification("Failed to copy task to clipboard")
		}
	}
}

func (app *App) OnKeypress(ev tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'Q':
			app.Quit()
		case '1':
			app.state.CurrentTab = state.APP_TAB_ACTIVE
		case '2':
			app.state.CurrentTab = state.APP_TAB_PROJECTS
		case '3':
			app.state.CurrentTab = state.APP_TAB_HISTORY
		case 'q':
			switch app.state.CurrentTab {
			case state.APP_TAB_PROJECTS:
				app.state.CurrentTab = state.APP_TAB_ACTIVE
			case state.APP_TAB_HISTORY:
				app.state.CurrentTab = state.APP_TAB_PROJECTS
			}
		case 'w':
			switch app.state.CurrentTab {
			case state.APP_TAB_ACTIVE:
				app.state.CurrentTab = state.APP_TAB_PROJECTS
			case state.APP_TAB_PROJECTS:
				app.state.CurrentTab = state.APP_TAB_HISTORY
			}
		case 'g':
			app.handleNavigateTop()
		case 'G':
			app.handleNavigateBottom()
		case 'y':
			app.handleYank()
		}
	}

	// active
	if app.state.CurrentTab == state.APP_TAB_ACTIVE {
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'j':
				app.handleActiveNavigateDown()
			case 'k':
				app.handleActiveNavigateUp()
			case ' ':
				app.handleActiveToggleInProgress()
			}
		case tcell.KeyCtrlC:
			app.handleActiveDeactivateTask()
		case tcell.KeyEnter:
			app.handleActiveEdit()
		case tcell.KeyCtrlSpace:
			app.handleActiveToggleDone()
		}
	}

	// history
	if app.state.CurrentTab == state.APP_TAB_HISTORY {
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'j':
				app.handleHistoryScrollDown()
			case 'd':
				app.handleHistoryScrollDown()
			case 'k':
				app.handleHistoryScrollUp()
			case 'u':
				app.handleHistoryScrollUp()
			case 's':
				app.handleHistoryScrollUp()
			}
		}
	}

	// projects
	if app.state.CurrentTab == state.APP_TAB_PROJECTS {
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'j':
				app.handleProjectsTaskNavigation(true)
			case 'k':
				app.handleProjectsTaskNavigation(false)
			case 'J':
				app.handleProjectsNavigation(true)
			case 'K':
				app.handleProjectsNavigation(false)
			case ' ':
				app.handleProjectsToggleTask()
			}
		case tcell.KeyTab:
			app.handleProjectsNavigation(true)
		case tcell.KeyBacktab:
			app.handleProjectsNavigation(false)
		case tcell.KeyEnter:
			app.handleProjectsEdit()
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
	app.state.HeaderHeight = headerBuffer.Height()

	// guard min height
	if app.Height() <= app.state.HeaderHeight {
		return b
	}

	// active tab
	if app.state.CurrentTab == state.APP_TAB_ACTIVE {
		taskList := components.TaskList{
			Tasks:                app.state.ActiveTasks,
			Width:                app.Width(),
			LongestProjectLength: app.state.LongestActiveProjectLength,
			SelectedIndex:        app.state.ActiveSelectedIndex,
			ScrollOffset:         app.state.ActiveScrollOffset,
			MaxHeight:            app.Height() - app.state.HeaderHeight,
		}
		b.DrawComponent(0, headerBuffer.Height(), &taskList)
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

	// projects tab
	if app.state.CurrentTab == state.APP_TAB_PROJECTS {
		projectSidebar := components.ProjectSidebar{
			AppState:     &app.state,
			Width:        app.state.LongestProjectLength + 2,
			Height:       app.Height() - headerBuffer.Height(),
			ScrollOffset: app.state.ProjectScrollOffset,
		}
		b.DrawComponent(0, headerBuffer.Height(), &projectSidebar)

		projectTaskList := components.ProjectTaskList{
			AppState: &app.state,
			Width:    app.Width() - projectSidebar.Width,
			Height:   app.Height() - headerBuffer.Height(),
		}
		b.DrawComponent(projectSidebar.Width, headerBuffer.Height(), &projectTaskList)
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
