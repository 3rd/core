package theme

import (
	ui "github.com/3rd/go-futui"
)

var (
	// app
	BG ui.Color = "#2E2442"
	FG ui.Color = "#D9C7E6"

	// header
	HEADER_BG         ui.Color = "#35294D"
	HEADER_FG         ui.Color = HEADER_BG.OptimalForeground()
	HEADER_BG_FOCUSED ui.Color = "#633650"
	HEADER_FG_FOCUSED ui.Color = HEADER_BG_FOCUSED.OptimalForeground()

	// tasks
	TASK_BG        ui.Color = "#503E74" // default
	TASK_FG        ui.Color = TASK_BG.OptimalForeground()
	TASK_ACTIVE_BG ui.Color = "#731F9B" // active
	TASK_ACTIVE_FG ui.Color = TASK_ACTIVE_BG.OptimalForeground()
	TASK_STICKY_BG ui.Color = "#813160" // sticky
	TASK_STICKY_FG ui.Color = TASK_STICKY_BG.OptimalForeground()
	PROJECT_FG     ui.Color = "#E8AB0F" // project
)
