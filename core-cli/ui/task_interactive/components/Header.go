package components

import (
	state "core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"core/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/3rd/core/core-lib/wiki"
	ui "github.com/3rd/go-futui"
)

type Header struct {
	ui.Component
	AppState *state.AppState
	Width    int
}

func (c *Header) Render() ui.Buffer {
	b := ui.Buffer{}

	// styles
	bgStyle := ui.Style{Background: theme.HEADER_BG, Foreground: theme.HEADER_FG}
	leftStyle := bgStyle
	leftStyle.Background = leftStyle.Background.Darken(0.03)
	rightStyle := leftStyle

	// bg
	b.Resize(c.Width, 4)
	b.FillStyle(bgStyle)

	// left
	left := ui.Buffer{}
	left.Resize(1, 4)

	// left: label
	leftLabel := ui.Buffer{}
	text := ""
	leftLabel.Text(0, 0, text, leftStyle)
	text = fmt.Sprintf("%d", c.AppState.GetNotDoneTasksCount())
	leftLabel.Text(leftLabel.Width()+1, 0, text, leftStyle)
	text = ""
	leftLabel.Text(leftLabel.Width()+1, 0, text, leftStyle)
	text = fmt.Sprintf("%d", c.AppState.GetDoneTasksCount())
	leftLabel.Text(leftLabel.Width()+1, 0, text, leftStyle)

	// left: bar
	leftBar := ui.Buffer{}
	barWidth := leftLabel.Width()
	min := c.AppState.LongestActiveProjectLength + 3

	if barWidth < min {
		barWidth = min
	}
	var midPoint = float64(barWidth) * (float64(c.AppState.GetDoneTasksCount()) / (float64(c.AppState.GetDoneTasksCount()) + float64(c.AppState.GetNotDoneTasksCount())))

	for i := 0; i < barWidth; i++ {
		ch := "▭"
		if float64(i) < midPoint {
			ch = "▬"
		}
		leftBar.Text(i, 0, ch, leftStyle)
	}

	left.DrawBuffer(1, 1, leftLabel)
	left.DrawBuffer(1, 2, leftBar)
	left.Resize(left.Width()+1, left.Height())
	left.ApplyStyle(leftStyle)

	// right
	right := ui.Buffer{}
	right.Resize(1, 4)

	// compute total work time and reward points
	totalWorkTime := time.Duration(0)
	totalRewardPoints := 0
	now := time.Now()
	for _, t := range c.AppState.ActiveTasks {
		totalWorkTime += t.GetTotalSessionTimeForDate(now)
		if t.Status == wiki.TASK_STATUS_DONE {
			totalRewardPoints += utils.ComputeTaskReward(t)
		}
	}

	// right: work time
	rightWorkTime := ui.Buffer{}
	rightWorkTimeText := totalWorkTime.Round(time.Second).String()
	rightWorkTime.Text(0, 0, rightWorkTimeText, rightStyle)

	// right: points
	rightRewardPoints := ui.Buffer{}
	rightRewardPointsText := strconv.Itoa(totalRewardPoints)
	rightRewardPoints.Text(0, 0, "", ui.Style{Foreground: theme.HEADER_REWARD_FG})
	rightRewardPoints.Text(2, 0, rightRewardPointsText, rightStyle)

	// draw right
	rightWidth := rightWorkTime.Width() + rightRewardPoints.Width() + 2
	right.Resize(rightWidth, right.Height())
	right.DrawBuffer(rightWidth-rightWorkTime.Width()-1, 1, rightWorkTime)
	right.ApplyStyle(rightStyle)
	right.DrawBuffer(rightWidth-rightRewardPoints.Width()-1, 2, rightRewardPoints)

	// draw left/right
	b.DrawBuffer(0, 0, left)
	rightX := c.Width - right.Width()
	if rightX < 0 {
		rightX = 0
	}
	b.DrawBuffer(rightX, 0, right)

	// tabs
	tabsBuffer := ui.Buffer{}
	activeTabStyle := ui.Style{Background: theme.TAB_ACTIVE_BG, Foreground: theme.TAB_ACTIVE_FG}
	inactiveTabStyle := ui.Style{Background: theme.TAB_INACTIVE_BG, Foreground: theme.TAB_INACTIVE_FG}

	activeTab := " (1) Active "
	if c.AppState.CurrentTab == state.APP_TAB_ACTIVE {
		tabsBuffer.Text(0, 0, activeTab, activeTabStyle)
	} else {
		tabsBuffer.Text(0, 0, activeTab, inactiveTabStyle)
	}

	projectsTab := " (2) Projects "
	if c.AppState.CurrentTab == state.APP_TAB_PROJECTS {
		tabsBuffer.Text(len(activeTab), 0, projectsTab, activeTabStyle)
	} else {
		tabsBuffer.Text(len(activeTab), 0, projectsTab, inactiveTabStyle)
	}

	historyTab := " (3) History "
	if c.AppState.CurrentTab == state.APP_TAB_HISTORY {
		tabsBuffer.Text(len(activeTab)+len(projectsTab), 0, historyTab, activeTabStyle)
	} else {
		tabsBuffer.Text(len(activeTab)+len(projectsTab), 0, historyTab, inactiveTabStyle)
	}

	b.DrawBuffer(c.Width/2-tabsBuffer.Width()/2, 1, tabsBuffer)
	return b
}
