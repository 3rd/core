package components

import (
	"core/ui/task_interactive/theme"
	"core/utils"
	"regexp"
	"strconv"
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

var taskLabelRegex = regexp.MustCompile(`^([a-zA-Z0-9_-]+:)`)

func (c *TaskItem) Render() ui.Buffer {
	b := ui.Buffer{}
	hoffset := 0

	taskReward := utils.ComputeTaskReward(c.Task)

	// styles
	taskStyle := ui.Style{Background: theme.TASK_BG, Foreground: theme.TASK_FG}
	projectStyle := taskStyle
	rewardStyle := ui.Style{Foreground: theme.TASK_REWARD_FG}
	projectStyle.Background = taskStyle.Background.Darken(0.05)
	projectStyle.Foreground = theme.PROJECT_FG

	if taskReward >= 100 {
		taskStyle.Background = theme.TASK_STICKY_BG
		taskStyle.Foreground = theme.TASK_STICKY_FG
		projectStyle.Background = theme.TASK_STICKY_BG
		projectStyle.Foreground = theme.TASK_CURRENT_FG
	}

	if c.Task.IsInProgress() {
		taskStyle.Background = theme.TASK_CURRENT_BG
		taskStyle.Foreground = theme.TASK_CURRENT_FG
		projectStyle.Background = taskStyle.Background.Darken(0.05)
		projectStyle.Foreground = projectStyle.Foreground.Lighten(0.1)
		rewardStyle.Foreground = projectStyle.Foreground
	} else if c.Task.Status == wiki.TASK_STATUS_DONE {
		taskStyle.Background = theme.TASK_DONE_BG
		taskStyle.Foreground = theme.TASK_DONE_FG
		projectStyle.Background = taskStyle.Background.Darken(0.05)
		projectStyle.Foreground = theme.PROJECT_DONE_FG
		rewardStyle.Foreground = projectStyle.Foreground
	}

	if c.Selected {

		if c.Task.IsInProgress() {
			taskStyle.Background = theme.TASK_CURRENT_SELECTED_BG
			taskStyle.Foreground = theme.TASK_CURRENT_SELECTED_FG
			projectStyle.Background = taskStyle.Background.Darken(0.05)
		} else {
			taskStyle.Background = theme.TASK_SELECTED_BG
			taskStyle.Foreground = theme.TASK_SELECTED_FG
			projectStyle.Background = taskStyle.Background.Darken(0.05)
		}
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

	// label
	if taskLabelRegex.MatchString(c.Task.Text) {
		labelText := taskLabelRegex.FindStringSubmatch(c.Task.Text)[1]
		label := ui.Buffer{}
		label.Text(0, 0, labelText, ui.Style{Foreground: theme.TASK_LABEL_FG})
		b.DrawBuffer(hoffset, 0, label)
	}

	// reward
	reward := ui.Buffer{}
	rewardIcon := ui.Buffer{}
	rewardIcon.Text(0, 0, "", rewardStyle)
	reward.DrawBuffer(0, 0, rewardIcon)
	reward.Text(2, 0, strconv.Itoa(int(taskReward)), rewardStyle)
	b.DrawBuffer(c.Width-reward.Width()-1, 0, reward)

	// duration
	workTime := c.Task.GetTotalSessionTime()
	if workTime > 0 {
		duration := ui.Buffer{}
		durationText := workTime.Round(time.Second).String()
		duration.Text(0, 0, durationText, taskStyle)
		b.DrawBuffer(c.Width-duration.Width()-reward.Width()-2, 0, duration)
	}

	return b
}
