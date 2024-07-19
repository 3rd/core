package components

import (
	"core/ui/task_interactive/theme"
	"time"

	"github.com/3rd/core/core-lib/wiki"
	ui "github.com/3rd/go-futui"
)

type TaskItem struct {
	ui.Component
	Task                 *wiki.Task
	Width                int
	LongestProjectLength int
	Selected             bool
}

func (c *TaskItem) Render() ui.Buffer {
	b := ui.Buffer{}
	hoffset := 0

	// styles
	taskStyle := ui.Style{Background: theme.TASK_BG, Foreground: theme.TASK_FG}
	projectStyle := taskStyle
	rewardStyle := ui.Style{Foreground: "#ffaa00"}
	projectStyle.Background = taskStyle.Background.Darken(0.05)
	projectStyle.Foreground = theme.PROJECT_FG

	if c.Task.Status == wiki.TASK_STATUS_DONE {
		taskStyle.Background = taskStyle.Background.Darken(0.05)
		taskStyle.Foreground = taskStyle.Foreground.Darken(0.5)
		projectStyle.Background = taskStyle.Background.Darken(0.05)
		projectStyle.Foreground = projectStyle.Foreground.Desaturate(1).Darken(0.1)
		rewardStyle.Foreground = rewardStyle.Foreground.Desaturate(1).Darken(0.1)
	}

	if c.Task.IsInProgress() {
		taskStyle.Background = theme.TASK_ACTIVE_BG
		taskStyle.Foreground = theme.TASK_ACTIVE_FG
		projectStyle.Background = taskStyle.Background.Darken(0.05)
		projectStyle.Foreground = projectStyle.Foreground.Lighten(0.1)
		rewardStyle.Foreground = taskStyle.Foreground.Desaturate(0.5).Lighten(0.1)
	}

	if c.Selected {
		taskStyle.Background = taskStyle.Background.Lighten(0.15)
		taskStyle.Foreground = taskStyle.Background.OptimalForeground()
		projectStyle.Background = projectStyle.Background.Lighten(0.1)
		projectStyle.Foreground = projectStyle.Foreground.Lighten(0.1)
	}

	checkmarkStyle := taskStyle
	textStyle := taskStyle

	// draw bg
	bg := ui.Buffer{}
	bg.Resize(c.Width, 1)
	bg.FillStyle(taskStyle)
	b.DrawBuffer(0, 0, bg)

	// checkmark
	checkmark := ui.Buffer{}
	checkmarkText := ""
	if c.Task.IsDone() {
		checkmarkText = ""
	}
	checkmark.Text(1, 0, checkmarkText, checkmarkStyle)
	b.DrawBuffer(0, 0, checkmark)
	hoffset = hoffset + checkmark.Width() + 1

	// project
	projectText := ""
	if c.Task.Node != nil {
		projectText = c.Task.Node.GetName()
	}
	project := ui.Buffer{}
	project.Resize(c.LongestProjectLength+2, 1)
	project.FillStyle(projectStyle)
	project.Text(0, 0, " "+projectText, projectStyle)
	b.DrawBuffer(hoffset, 0, project)
	hoffset = hoffset + project.Width() + 1

	// text
	text := ui.Buffer{}
	text.Text(0, 0, c.Task.Text, textStyle)
	b.DrawBuffer(hoffset, 0, text)

	// duration
	workTime := c.Task.GetWorkTime()
	if workTime > 0 {
		duration := ui.Buffer{}
		durationText := c.Task.GetWorkTime().Round(time.Second).String()
		duration.Text(0, 0, durationText, taskStyle)
		b.DrawBuffer(c.Width-duration.Width()-2, 0, duration)
	}

	return b
}
