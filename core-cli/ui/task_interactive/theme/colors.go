package theme

import (
	ui "github.com/3rd/go-futui"
)

var (
	// app
	BG ui.Color = "#222138"
	FG ui.Color = "#c8c6ce"

	// header
	HEADER_BG        ui.Color = BG
	HEADER_FG        ui.Color = "#c8c6d2"
	HEADER_REWARD_FG ui.Color = "#0aaf50"
	TAB_ACTIVE_BG    ui.Color = "#35333e"
	TAB_ACTIVE_FG    ui.Color = FG
	TAB_INACTIVE_BG  ui.Color = "#212027"
	TAB_INACTIVE_FG  ui.Color = "#767481"

	// tasks
	TASK_BG                ui.Color = "#2b2a37"
	TASK_FG                ui.Color = FG
	TASK_DONE_BG           ui.Color = "#1d1c20"
	TASK_DONE_FG           ui.Color = "#7b7986"
	TASK_CURRENT_BG        ui.Color = "#9a2860"
	TASK_CURRENT_FG        ui.Color = TASK_CURRENT_BG.OptimalForeground()
	TASK_PROJECT_BG        ui.Color = "#272830"
	TASK_REWARD_DEFAULT_FG ui.Color = "#ffaa00"
	TASK_REWARD_MEDIUM_FG  ui.Color = "#f26c0d"
	TASK_REWARD_HIGH_FG    ui.Color = "#f2330d"
	TASK_PROJECT_FG        ui.Color = "#E8AB0F"
	TASK_PROJECT_DONE_FG   ui.Color = "#413E51"
	TASK_LABEL_FG          ui.Color = "#f069cb"

	// history
	HISTORY_DATE_FG    ui.Color = "#E8AB0F"
	HISTORY_TASK_FG    ui.Color = FG
	HISTORY_PROJECT_FG ui.Color = "#7b7986"

	// projects
	PROJECT_SIDEBAR_BG          ui.Color = "#383545"
	PROJECT_SIDEBAR_FG          ui.Color = FG
	PROJECT_SIDEBAR_SELECTED_BG ui.Color = "#5F5f58"
	PROJECT_SIDEBAR_SELECTED_FG ui.Color = FG
	PROJECTS_TASK_ACTIVE_FG     ui.Color = "#ffaa00"

	// modal
	MODAL_BG          ui.Color = "#2C2A37"
	MODAL_BORDER_FG   ui.Color = "#f069cb"
	MODAL_SELECTED_BG ui.Color = "#413E51"
	MODAL_SELECTED_FG ui.Color = FG
)
