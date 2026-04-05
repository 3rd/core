package components

import (
	"core/ui/task_interactive/theme"
	"core/utils"
	"regexp"
	"strconv"
	"strings"
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

func RenderProjectColumnText(projectName string) string {
	projectText := strings.TrimPrefix(projectName, "project-")
	return strings.Replace(projectText, "project:", "p:", 1)
}

func GetRenderedProjectColumnWidth(tasks []*wiki.Task) int {
	longest := 0
	for _, task := range tasks {
		if task.Node == nil {
			continue
		}

		projectText := RenderProjectColumnText(task.Node.GetName())
		if len(projectText) > longest {
			longest = len(projectText)
		}
	}

	return longest
}

func (c *TaskItem) Render() ui.Buffer {
	b := ui.Buffer{}
	hoffset := 0

	taskReward := utils.ComputeTaskReward(c.Task)
	isInProgress := c.Task.IsInProgress()
	isDone := c.Task.Status == wiki.TASK_STATUS_DONE

	// styles
	taskStyle := theme.TaskRowStyle(isInProgress, isDone, c.Selected)
	projectStyle := theme.TaskProjectStyle(isInProgress, isDone, c.Selected)
	rewardStyle := theme.TaskRewardStyle(taskReward, isInProgress, isDone)

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
		projectText = RenderProjectColumnText(c.Task.Node.GetName())
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
		label.Text(0, 0, labelText, theme.TASK_LABEL_STYLE)
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
	now := time.Now()
	workTime := c.Task.GetTotalSessionTimeForDate(now)
	if workTime > 0 {
		duration := ui.Buffer{}
		durationText := workTime.Round(time.Second).String()
		duration.Text(0, 0, durationText, taskStyle)
		b.DrawBuffer(c.Width-duration.Width()-reward.Width()-2, 0, duration)
	}

	return b
}
