package components

import (
	state "core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"
	"time"

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
	if c.AppState.Mode == state.APP_MODE_FOCUS {
		bgStyle = ui.Style{Background: theme.HEADER_BG_FOCUSED, Foreground: theme.HEADER_FG_FOCUSED}
	}
	leftStyle := bgStyle
	leftStyle.Background = leftStyle.Background.Darken(0.1)
	leftStyle.Foreground = leftStyle.Foreground.Darken(0.1)
	rightStyle := leftStyle

	// bg
	b.Resize(c.Width, 4)
	b.FillStyle(bgStyle)

	// left
	left := ui.Buffer{}
	left.Resize(1, 4)

	// left label
	leftLabel := ui.Buffer{}
	text := ""
	leftLabel.Text(0, 0, text, leftStyle)
	text = fmt.Sprintf("%d", c.AppState.GetNotDoneTasksCount())
	leftLabel.Text(leftLabel.Width()+1, 0, text, leftStyle)
	text = ""
	leftLabel.Text(leftLabel.Width()+1, 0, text, leftStyle)
	text = fmt.Sprintf("%d", c.AppState.GetDoneTasksCount())
	leftLabel.Text(leftLabel.Width()+1, 0, text, leftStyle)

	// left bar
	leftBar := ui.Buffer{}
	barWidth := leftLabel.Width()
	min := c.AppState.LongestProjectLength + 3

	if barWidth < min {
		barWidth = min
	}

	var point = float64(barWidth) * (float64(c.AppState.GetDoneTasksCount()) / (float64(c.AppState.GetDoneTasksCount()) + float64(c.AppState.GetNotDoneTasksCount())))

	for i := 0; i < barWidth; i++ {
		ch := "▭"
		if float64(i) < point {
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

	dayTime := time.Duration(0)
	for _, t := range c.AppState.Tasks {
		dayTime += t.GetWorkTime()
	}

	rightTime := ui.Buffer{}
	rightText := dayTime.Round(time.Second).String()
	rightTime.Text(0, 0, rightText, rightStyle)

	right.DrawBuffer(1, 1, rightTime)
	right.Resize(right.Width()+1, right.Height())
	right.ApplyStyle(rightStyle)

	b.DrawBuffer(0, 0, left)
	rightX := c.Width - right.Width()
	if rightX < 0 {
		rightX = 0
	}
	b.DrawBuffer(rightX, 0, right)

	// center
	if c.AppState.Mode == state.APP_MODE_FOCUS {
		center := ui.Buffer{}
		center.Resize(c.Width-left.Width()-right.Width(), 4)
		center.FillStyle(bgStyle)
		text := "FOCUS"
		center.Text(center.Width()/2-len(text)/2, 1, text, bgStyle)
		b.DrawBuffer(left.Width(), 0, center)
	}

	return b
}