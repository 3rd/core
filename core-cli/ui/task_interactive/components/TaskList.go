package components

import (
	"github.com/3rd/core/core-lib/wiki"
	ui "github.com/3rd/go-futui"
)

type TaskList struct {
	ui.Component
	Tasks                []*wiki.Task
	Width                int
	SelectedIndex        int
	LongestProjectLength int
}

func (c *TaskList) Render() ui.Buffer {
	b := ui.Buffer{}
	voffset := 0

	for i, task := range c.Tasks {
		taskComponent := TaskItem{
			Task:                 task,
			Width:                c.Width,
			LongestProjectLength: c.LongestProjectLength,
			Selected:             i == c.SelectedIndex,
		}

		b.DrawComponent(0, voffset, &taskComponent)
		voffset = b.Height()
	}

	return b
}
