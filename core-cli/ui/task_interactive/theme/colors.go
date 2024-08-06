package theme

import (
	ui "github.com/3rd/go-futui"
)

var (
	// app
	BG ui.Color = "#23212B"
	FG ui.Color = "#C5C2D6"

	// header
	HEADER_BG        ui.Color = "#2C2A37"
	HEADER_FG        ui.Color = "#C5C2D6"
	HEADER_REWARD_FG ui.Color = "#0aaf50"
	TAB_ACTIVE_BG    ui.Color = "#413E51"
	TAB_ACTIVE_FG    ui.Color = "#C5C2D6"
	TAB_INACTIVE_BG  ui.Color = "#383545"
	TAB_INACTIVE_FG  ui.Color = "#A5A0BA"

	// tasks
	TASK_BG                ui.Color = "#4B475C"
	TASK_FG                ui.Color = "#C5C2D6"
	TASK_DONE_BG           ui.Color = "#383545"
	TASK_DONE_FG           ui.Color = "#827d98"
	TASK_CURRENT_BG        ui.Color = "#7a0891"
	TASK_CURRENT_FG        ui.Color = TASK_CURRENT_BG.OptimalForeground()
	TASK_REWARD_DEFAULT_FG ui.Color = "#ffaa00"
	TASK_REWARD_MEDIUM_FG  ui.Color = "#f26c0d"
	TASK_REWARD_HIGH_FG    ui.Color = "#f2330d"
	PROJECT_FG             ui.Color = "#E8AB0F"
	PROJECT_DONE_FG        ui.Color = "#746f8c"
	TASK_LABEL_FG          ui.Color = "#f069cb"

	// history
	HISTORY_DATE_FG    ui.Color = "#E8AB0F"
	HISTORY_TASK_FG    ui.Color = "#C5C2D6"
	HISTORY_PROJECT_FG ui.Color = "#9B96B0"

	// projects
	PROJECT_SIDEBAR_BG          ui.Color = "#383545"
	PROJECT_SIDEBAR_FG          ui.Color = "#B0ACC3"
	PROJECT_SIDEBAR_SELECTED_BG ui.Color = "#5F5B76"
	PROJECT_SIDEBAR_SELECTED_FG ui.Color = "#C5C2D6"
	PROJECTS_TASK_ACTIVE_FG     ui.Color = "#ffaa00"
)
