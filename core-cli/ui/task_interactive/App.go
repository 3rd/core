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

	// apply all filters
	app.applyAllFilters()

	// guard out of bounds
	if app.state.ActiveSelectedIndex >= len(app.state.FilteredTasks) {
		app.state.ActiveSelectedIndex = max(len(app.state.FilteredTasks)-1, 0)
	}
}

// applyAllFilters starting from ActiveTasks:
// 1. time filter (t, for done tasks)
// 2. focus filter (f)
// 3. project filters (p)
func (app *App) applyAllFilters() {
	filteredTasks := app.state.ActiveTasks

	// time filter
	if app.state.ActiveTimeFilter == state.TimeFilterToday {
		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		newFiltered := []*wiki.Task{}
		for _, task := range filteredTasks {
			// for done tasks, only include if done today
			if task.Status == wiki.TASK_STATUS_DONE {
				// check last session for regular tasks
				lastSession := task.GetLastSession()
				if lastSession != nil && !lastSession.Start.Before(todayStart) {
					newFiltered = append(newFiltered, task)
					continue
				}
				
				// for recurrent tasks, also check completion
				if task.Schedule != nil && task.Schedule.Repeat != "" {
					completion := task.GetCompletionForDate(now)
					if completion != nil {
						newFiltered = append(newFiltered, task)
					}
				}
			} else {
				// not done tasks always pass through
				newFiltered = append(newFiltered, task)
			}
		}
		filteredTasks = newFiltered
	}

	// focus filter
	if app.state.ActiveFocusedProjectID != "" {
		newFiltered := []*wiki.Task{}
		for _, task := range filteredTasks {
			if task.Node != nil && task.Node.GetID() == app.state.ActiveFocusedProjectID {
				newFiltered = append(newFiltered, task)
			}
		}
		filteredTasks = newFiltered
	} else {
		// project filters
		if len(app.state.ProjectFilterModal.FilteredProjects) > 0 {
			// check if any project is disabled
			hasDisabledProjects := false
			for _, enabled := range app.state.ProjectFilterModal.FilteredProjects {
				if !enabled {
					hasDisabledProjects = true
					break
				}
			}

			// only filter if some projects are disabled
			if hasDisabledProjects {
				newFiltered := []*wiki.Task{}
				for _, task := range filteredTasks {
					if task.Node != nil {
						projectID := task.Node.GetID()
						// Include task if project is enabled or not in filter list
						if enabled, exists := app.state.ProjectFilterModal.FilteredProjects[projectID]; !exists || enabled {
							newFiltered = append(newFiltered, task)
						}
					}
				}
				filteredTasks = newFiltered
			}
		}
	}

	app.state.FilteredTasks = filteredTasks
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
	if app.state.ActiveSelectedIndex >= len(app.state.FilteredTasks)-1 {
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

func (app *App) handleActiveToggleFocus() {
	if app.state.ActiveFocusedProjectID == "" {
		if len(app.state.FilteredTasks) > 0 && app.state.ActiveSelectedIndex < len(app.state.FilteredTasks) {
			focusedTask := app.state.FilteredTasks[app.state.ActiveSelectedIndex]
			if focusedTask.Node != nil {
				app.state.ActiveFocusedProjectID = focusedTask.Node.GetID()
			}
		}
	} else {
		app.state.ActiveFocusedProjectID = ""
	}

	// reapply filters
	app.applyAllFilters()

	// guard selected index
	if app.state.ActiveSelectedIndex >= len(app.state.FilteredTasks) {
		app.state.ActiveSelectedIndex = max(len(app.state.FilteredTasks)-1, 0)
	}

	app.Update()
}

func (app *App) adjustActiveScroll() {
	maxVisibleTasks := app.Height() - app.state.HeaderHeight
	if maxVisibleTasks <= 0 {
		return
	}
	maxScrollOffset := max(len(app.state.FilteredTasks)-maxVisibleTasks, 0)

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

func (app *App) handleToggleTimeFilter() {
	if app.state.ActiveTimeFilter == state.TimeFilterToday {
		app.state.ActiveTimeFilter = state.TimeFilter24Hours
	} else {
		app.state.ActiveTimeFilter = state.TimeFilterToday
	}
	app.applyAllFilters()

	// guard selected index
	if app.state.ActiveSelectedIndex >= len(app.state.FilteredTasks) {
		app.state.ActiveSelectedIndex = len(app.state.FilteredTasks) - 1
		if app.state.ActiveSelectedIndex < 0 {
			app.state.ActiveSelectedIndex = 0
		}
	}

	app.Update()
}

func (app *App) handleActiveEdit() {
	task := app.state.FilteredTasks[app.state.ActiveSelectedIndex]
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
	task := app.state.FilteredTasks[app.state.ActiveSelectedIndex]
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
	task := app.state.FilteredTasks[app.state.ActiveSelectedIndex]
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
	if len(app.state.FilteredTasks) == 0 {
		return
	}

	task := app.state.FilteredTasks[app.state.ActiveSelectedIndex]
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
		if len(app.state.FilteredTasks) > 0 {
			app.state.ActiveSelectedIndex = len(app.state.FilteredTasks) - 1
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
		if app.state.ActiveSelectedIndex < len(app.state.FilteredTasks) {
			taskText = app.state.FilteredTasks[app.state.ActiveSelectedIndex].Text
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

func (app *App) handleShowProjectFilterModal() {
	tasksToShow := app.state.ActiveTasks

	// if in focus mode, only show the focused project
	// TODO: should exit focus mode instead?
	if app.state.ActiveFocusedProjectID != "" {
		focusedTasks := []*wiki.Task{}
		for _, task := range tasksToShow {
			if task.Node != nil && task.Node.GetID() == app.state.ActiveFocusedProjectID {
				focusedTasks = append(focusedTasks, task)
			}
		}
		tasksToShow = focusedTasks
	}

	// collect unique projects from active tasks
	projectMap := make(map[string]*state.ProjectFilterItem)

	for _, task := range tasksToShow {
		if task.Node != nil {
			projectID := task.Node.GetID()
			if _, exists := projectMap[projectID]; !exists {
				projectName := task.Node.GetName()
				projectMap[projectID] = &state.ProjectFilterItem{
					ProjectID:   projectID,
					ProjectName: projectName,
					TaskCount:   0,
					IsEnabled:   true,
				}
			}
			projectMap[projectID].TaskCount++
		}
	}

	projects := []state.ProjectFilterItem{}
	totalTaskCount := 0
	for _, project := range projectMap {
		// check if project was previously filtered
		if app.state.ProjectFilterModal.FilteredProjects != nil {
			if enabled, exists := app.state.ProjectFilterModal.FilteredProjects[project.ProjectID]; exists {
				project.IsEnabled = enabled
			}
		}
		projects = append(projects, *project)
		totalTaskCount += project.TaskCount
	}

	// sort projects by name
	sort.Slice(projects, func(i, j int) bool {
		return strings.Compare(projects[i].ProjectName, projects[j].ProjectName) < 0
	})

	// add "All" item at the beginning
	allEnabled := true
	for _, project := range projects {
		if !project.IsEnabled {
			allEnabled = false
			break
		}
	}

	allItem := state.ProjectFilterItem{
		ProjectID:   "__all__",
		ProjectName: "All",
		TaskCount:   totalTaskCount,
		IsEnabled:   allEnabled,
	}
	projects = append([]state.ProjectFilterItem{allItem}, projects...)

	if app.state.ProjectFilterModal.FilteredProjects == nil {
		app.state.ProjectFilterModal.FilteredProjects = make(map[string]bool)
	}

	// update modal state
	app.state.ProjectFilterModal.IsVisible = true
	app.state.ProjectFilterModal.CursorIndex = 0
	app.state.ProjectFilterModal.Projects = projects

	app.Update()
}

func (app *App) handleProjectFilterModalKeypress(ev tcell.EventKey) {
	modalState := &app.state.ProjectFilterModal

	switch ev.Key() {
	case tcell.KeyEscape:
		modalState.IsVisible = false
		app.applyProjectFilters()
		app.Update()
		return

	case tcell.KeyRune:
		switch ev.Rune() {
		case 'q':
			modalState.IsVisible = false
			app.applyProjectFilters()
			app.Update()
			return
		case 'j':
			if modalState.CursorIndex < len(modalState.Projects)-1 {
				modalState.CursorIndex++
			}
		case 'k':
			if modalState.CursorIndex > 0 {
				modalState.CursorIndex--
			}
		case ' ':
			if modalState.CursorIndex < len(modalState.Projects) {
				// handle "All" item
				if modalState.CursorIndex == 0 && len(modalState.Projects) > 0 && modalState.Projects[0].ProjectID == "__all__" {
					// toggle all projects
					newState := !modalState.Projects[0].IsEnabled
					for i := range modalState.Projects {
						modalState.Projects[i].IsEnabled = newState
						if i > 0 { // skip "All"
							modalState.FilteredProjects[modalState.Projects[i].ProjectID] = newState
						}
					}
				} else {
					// regular project toggle
					project := &modalState.Projects[modalState.CursorIndex]
					project.IsEnabled = !project.IsEnabled
					modalState.FilteredProjects[project.ProjectID] = project.IsEnabled

					// update "All" item state if it exists
					if len(modalState.Projects) > 0 && modalState.Projects[0].ProjectID == "__all__" {
						allEnabled := true
						for i := 1; i < len(modalState.Projects); i++ {
							if !modalState.Projects[i].IsEnabled {
								allEnabled = false
								break
							}
						}
						modalState.Projects[0].IsEnabled = allEnabled
					}
				}
			}
		}

	case tcell.KeyDown:
		if modalState.CursorIndex < len(modalState.Projects)-1 {
			modalState.CursorIndex++
		}

	case tcell.KeyUp:
		if modalState.CursorIndex > 0 {
			modalState.CursorIndex--
		}
	}

	app.Update()
}

func (app *App) applyProjectFilters() {
	app.applyAllFilters()

	// guard selected index
	if app.state.ActiveSelectedIndex >= len(app.state.FilteredTasks) {
		app.state.ActiveSelectedIndex = max(len(app.state.FilteredTasks)-1, 0)
	}

	app.Update()
}

func (app *App) OnKeypress(ev tcell.EventKey) {
	if app.state.CurrentTab == state.APP_TAB_ACTIVE && app.state.ProjectFilterModal.IsVisible {
		app.handleProjectFilterModalKeypress(ev)
		return
	}

	switch ev.Key() {
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'q':
			app.Quit()
		case '1':
			app.state.CurrentTab = state.APP_TAB_ACTIVE
		case '2':
			app.state.CurrentTab = state.APP_TAB_PROJECTS
		case '3':
			app.state.CurrentTab = state.APP_TAB_HISTORY
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
			case 'f':
				app.handleActiveToggleFocus()
			case 'p':
				app.handleShowProjectFilterModal()
			case 't':
				app.handleToggleTimeFilter()
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
			Tasks:                app.state.FilteredTasks,
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

	// project filter modal
	if app.state.ProjectFilterModal.IsVisible {
		modal := components.ProjectFilterModal{
			AppState: &app.state,
			Width:    app.Width(),
			Height:   app.Height(),
		}
		modalBuffer := modal.Render()
		modalX := (app.Width() - modalBuffer.Width()) / 2
		modalY := (app.Height() - modalBuffer.Height()) / 2
		b.DrawBuffer(modalX, modalY, modalBuffer)
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

	app := App{
		providers: providers,
	}
	app.Run(&app)
}
