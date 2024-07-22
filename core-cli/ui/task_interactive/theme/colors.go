package theme

import (
	ui "github.com/3rd/go-futui"
)

var (
	// app
	BG ui.Color = "#23212B"
	FG ui.Color = "#C5C2D6"

	// header
	HEADER_BG         ui.Color = "#2C2A37"
	HEADER_FG         ui.Color = "#C5C2D6"
	HEADER_BG_FOCUSED ui.Color = "#633650"
	HEADER_FG_FOCUSED ui.Color = HEADER_BG_FOCUSED.OptimalForeground()
	HEADER_REWARD_FG  ui.Color = "#0aaf50"

	// tasks
	TASK_BG                  ui.Color = "#4B475C"
	TASK_FG                  ui.Color = "#C5C2D6"
	TASK_DONE_BG             ui.Color = "#383545"
	TASK_DONE_FG             ui.Color = "#827d98"
	TASK_CURRENT_BG          ui.Color = "#7a0891"
	TASK_CURRENT_FG          ui.Color = TASK_CURRENT_BG.OptimalForeground()
	TASK_STICKY_BG           ui.Color = "#813160"
	TASK_STICKY_FG           ui.Color = TASK_STICKY_BG.OptimalForeground()
	TASK_SELECTED_BG         ui.Color = "#40666d"
	TASK_SELECTED_FG         ui.Color = "#5ee5e5"
	TASK_CURRENT_SELECTED_BG ui.Color = "#932ac2"
	TASK_CURRENT_SELECTED_FG ui.Color = TASK_CURRENT_SELECTED_BG.OptimalForeground()
	TASK_REWARD_FG           ui.Color = "#ffaa00"
	PROJECT_FG               ui.Color = "#E8AB0F"
	PROJECT_DONE_FG          ui.Color = "#746f8c"
	TASK_LABEL_FG            ui.Color = "#f069cb"
)
