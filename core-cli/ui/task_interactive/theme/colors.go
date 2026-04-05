package theme

import (
	ui "github.com/3rd/go-futui"
)

func style(background ui.Color, foreground ui.Color) ui.Style {
	return ui.Style{Background: background, Foreground: foreground}
}

func textStyle(foreground ui.Color) ui.Style {
	return ui.Style{Foreground: foreground}
}

func boldTextStyle(foreground ui.Color) ui.Style {
	return ui.Style{Foreground: foreground, Bold: true}
}

var (
	// every rendered color state must be declared here explicitly so component
	// code can consume semantic styles without mutating colors at render time.

	// app
	BG ui.Color = "#201f25"
	FG ui.Color = "#c8c6ce"

	// header
	HEADER_BG        ui.Color = "#23222C"
	HEADER_PANEL_BG  ui.Color = "#32313E"
	HEADER_FG        ui.Color = "#c8c6d2"
	HEADER_REWARD_FG ui.Color = "#0aaf50"
	TAB_ACTIVE_BG    ui.Color = "#3C394A"
	TAB_ACTIVE_FG    ui.Color = FG
	TAB_INACTIVE_BG  ui.Color = "#32313E"
	TAB_INACTIVE_FG  ui.Color = "#767481"

	// tasks
	TASK_BG                        ui.Color = "#2b2a37"
	TASK_FG                        ui.Color = FG
	TASK_SELECTED_BG               ui.Color = "#424054"
	TASK_SELECTED_FG               ui.Color = "#FDFDFD"
	TASK_CURRENT_BG                ui.Color = "#9a2860"
	TASK_CURRENT_FG                ui.Color = "#FFFFFF"
	TASK_CURRENT_SELECTED_BG       ui.Color = "#C23278"
	TASK_CURRENT_SELECTED_FG       ui.Color = "#FFFFFF"
	TASK_DONE_BG                   ui.Color = "#2B2932"
	TASK_DONE_FG                   ui.Color = "#7b7986"
	TASK_DONE_SELECTED_BG          ui.Color = "#36343B"
	TASK_DONE_SELECTED_FG          ui.Color = "#B0AFB6"
	TASK_REWARD_DEFAULT_FG         ui.Color = "#ffaa00"
	TASK_REWARD_MEDIUM_FG          ui.Color = "#f26c0d"
	TASK_REWARD_HIGH_FG            ui.Color = "#f1052f"
	TASK_REWARD_CURRENT_DEFAULT_FG ui.Color = "#FFCC66"
	TASK_REWARD_CURRENT_MEDIUM_FG  ui.Color = "#F7A76E"
	TASK_REWARD_CURRENT_HIGH_FG    ui.Color = "#F7856E"
	TASK_LABEL_FG                  ui.Color = "#f069cb"

	// tasks: project
	TASK_PROJECT_BG                  ui.Color = "#272830"
	TASK_PROJECT_FG                  ui.Color = "#E8AB0F"
	TASK_PROJECT_DONE_FG             ui.Color = "#4C495E"
	TASK_PROJECT_SELECTED_BG         ui.Color = "#373545"
	TASK_DONE_PROJECT_BG             ui.Color = "#212027"
	TASK_DONE_SELECTED_PROJECT_BG    ui.Color = "#36343B"
	TASK_CURRENT_PROJECT_BG          ui.Color = "#862353"
	TASK_CURRENT_SELECTED_PROJECT_BG ui.Color = "#98476E"

	// history
	HISTORY_BG          ui.Color = BG
	HISTORY_DATE_FG     ui.Color = "#E8AB0F"
	HISTORY_TASK_FG     ui.Color = FG
	HISTORY_PROJECT_FG  ui.Color = "#7b7986"
	HISTORY_DURATION_FG ui.Color = TASK_DONE_FG

	// projects
	PROJECT_SIDEBAR_BG               ui.Color = "#383545"
	PROJECT_SIDEBAR_FG               ui.Color = FG
	PROJECT_SIDEBAR_SELECTED_BG      ui.Color = "#5F5f58"
	PROJECT_SIDEBAR_SELECTED_FG      ui.Color = FG
	PROJECTS_TASK_ACTIVE_FG          ui.Color = "#ffaa00"
	PROJECTS_TASK_SELECTED_BG        ui.Color = "#424054"
	PROJECTS_TASK_SELECTED_FG        ui.Color = "#FDFDFD"
	PROJECTS_TASK_ACTIVE_SELECTED_FG ui.Color = "#FFCC66"

	// modal
	MODAL_BG          ui.Color = "#2C2A37"
	MODAL_BORDER_FG   ui.Color = "#f069cb"
	MODAL_SELECTED_BG ui.Color = "#413E51"
	MODAL_SELECTED_FG ui.Color = FG
	MODAL_DISABLED_FG ui.Color = TASK_DONE_FG
)

var (
	APP_STYLE = style(BG, FG)

	HEADER_STYLE        = style(HEADER_BG, HEADER_FG)
	HEADER_PANEL_STYLE  = style(HEADER_PANEL_BG, HEADER_FG)
	HEADER_REWARD_STYLE = textStyle(HEADER_REWARD_FG)
	TAB_ACTIVE_STYLE    = style(TAB_ACTIVE_BG, TAB_ACTIVE_FG)
	TAB_INACTIVE_STYLE  = style(TAB_INACTIVE_BG, TAB_INACTIVE_FG)

	TASK_STYLE                          = style(TASK_BG, TASK_FG)
	TASK_SELECTED_STYLE                 = style(TASK_SELECTED_BG, TASK_SELECTED_FG)
	TASK_CURRENT_STYLE                  = style(TASK_CURRENT_BG, TASK_CURRENT_FG)
	TASK_CURRENT_SELECTED_STYLE         = style(TASK_CURRENT_SELECTED_BG, TASK_CURRENT_SELECTED_FG)
	TASK_DONE_STYLE                     = style(TASK_DONE_BG, TASK_DONE_FG)
	TASK_DONE_SELECTED_STYLE            = style(TASK_DONE_SELECTED_BG, TASK_DONE_SELECTED_FG)
	TASK_PROJECT_STYLE                  = style(TASK_PROJECT_BG, TASK_PROJECT_FG)
	TASK_PROJECT_SELECTED_STYLE         = style(TASK_PROJECT_SELECTED_BG, TASK_PROJECT_FG)
	TASK_CURRENT_PROJECT_STYLE          = style(TASK_CURRENT_PROJECT_BG, TASK_CURRENT_FG)
	TASK_CURRENT_PROJECT_SELECTED_STYLE = style(TASK_CURRENT_SELECTED_PROJECT_BG, TASK_CURRENT_FG)
	TASK_DONE_PROJECT_STYLE             = style(TASK_DONE_PROJECT_BG, TASK_PROJECT_DONE_FG)
	TASK_DONE_PROJECT_SELECTED_STYLE    = style(TASK_DONE_SELECTED_PROJECT_BG, TASK_PROJECT_DONE_FG)
	TASK_LABEL_STYLE                    = textStyle(TASK_LABEL_FG)
	TASK_REWARD_DEFAULT_STYLE           = textStyle(TASK_REWARD_DEFAULT_FG)
	TASK_REWARD_MEDIUM_STYLE            = textStyle(TASK_REWARD_MEDIUM_FG)
	TASK_REWARD_HIGH_STYLE              = textStyle(TASK_REWARD_HIGH_FG)
	TASK_REWARD_CURRENT_DEFAULT_STYLE   = textStyle(TASK_REWARD_CURRENT_DEFAULT_FG)
	TASK_REWARD_CURRENT_MEDIUM_STYLE    = textStyle(TASK_REWARD_CURRENT_MEDIUM_FG)
	TASK_REWARD_CURRENT_HIGH_STYLE      = textStyle(TASK_REWARD_CURRENT_HIGH_FG)
	TASK_REWARD_DONE_STYLE              = textStyle(TASK_PROJECT_DONE_FG)

	HISTORY_STYLE          = style(HISTORY_BG, FG)
	HISTORY_DATE_STYLE     = style(HISTORY_BG, HISTORY_DATE_FG)
	HISTORY_TASK_STYLE     = style(HISTORY_BG, HISTORY_TASK_FG)
	HISTORY_PROJECT_STYLE  = style(HISTORY_BG, HISTORY_PROJECT_FG)
	HISTORY_DURATION_STYLE = style(HISTORY_BG, HISTORY_DURATION_FG)

	PROJECT_SIDEBAR_STYLE               = style(PROJECT_SIDEBAR_BG, PROJECT_SIDEBAR_FG)
	PROJECT_SIDEBAR_SELECTED_STYLE      = style(PROJECT_SIDEBAR_SELECTED_BG, PROJECT_SIDEBAR_SELECTED_FG)
	PROJECTS_TASK_STYLE                 = style(TASK_BG, TASK_FG)
	PROJECTS_TASK_ACTIVE_STYLE          = style(TASK_BG, PROJECTS_TASK_ACTIVE_FG)
	PROJECTS_TASK_SELECTED_STYLE        = style(PROJECTS_TASK_SELECTED_BG, PROJECTS_TASK_SELECTED_FG)
	PROJECTS_TASK_ACTIVE_SELECTED_STYLE = style(PROJECTS_TASK_SELECTED_BG, PROJECTS_TASK_ACTIVE_SELECTED_FG)

	MODAL_STYLE                          = style(MODAL_BG, FG)
	MODAL_BORDER_STYLE                   = textStyle(MODAL_BORDER_FG)
	MODAL_TITLE_STYLE                    = boldTextStyle(HEADER_FG)
	MODAL_HELP_STYLE                     = textStyle(TASK_LABEL_FG)
	MODAL_PROJECT_ENABLED_STYLE          = style(MODAL_BG, FG)
	MODAL_PROJECT_DISABLED_STYLE         = style(MODAL_BG, MODAL_DISABLED_FG)
	MODAL_PROJECT_SELECTED_ENABLED_STYLE = ui.Style{
		Background: MODAL_SELECTED_BG,
		Foreground: TASK_PROJECT_FG,
		Bold:       true,
	}
	MODAL_PROJECT_SELECTED_DISABLED_STYLE = ui.Style{
		Background: MODAL_SELECTED_BG,
		Foreground: MODAL_SELECTED_FG,
		Bold:       true,
	}

	HELP_SECTION_STYLE = boldTextStyle(TASK_PROJECT_FG)
	HELP_BINDING_STYLE = textStyle(FG)
	HELP_KEY_STYLE     = textStyle(TASK_LABEL_FG)
	NOTIFICATION_STYLE = textStyle(TASK_LABEL_FG)
)

func TaskRowStyle(isInProgress bool, isDone bool, isSelected bool) ui.Style {
	switch {
	case isInProgress && isSelected:
		return TASK_CURRENT_SELECTED_STYLE
	case isInProgress:
		return TASK_CURRENT_STYLE
	case isDone && isSelected:
		return TASK_DONE_SELECTED_STYLE
	case isDone:
		return TASK_DONE_STYLE
	case isSelected:
		return TASK_SELECTED_STYLE
	default:
		return TASK_STYLE
	}
}

func TaskProjectStyle(isInProgress bool, isDone bool, isSelected bool) ui.Style {
	switch {
	case isInProgress && isSelected:
		return TASK_CURRENT_PROJECT_SELECTED_STYLE
	case isInProgress:
		return TASK_CURRENT_PROJECT_STYLE
	case isDone && isSelected:
		return TASK_DONE_PROJECT_SELECTED_STYLE
	case isDone:
		return TASK_DONE_PROJECT_STYLE
	case isSelected:
		return TASK_PROJECT_SELECTED_STYLE
	default:
		return TASK_PROJECT_STYLE
	}
}

func TaskRewardStyle(points int, isInProgress bool, isDone bool) ui.Style {
	if isDone {
		return TASK_REWARD_DONE_STYLE
	}

	if isInProgress {
		if points >= 100 {
			return TASK_REWARD_CURRENT_HIGH_STYLE
		}
		if points > 10 {
			return TASK_REWARD_CURRENT_MEDIUM_STYLE
		}
		return TASK_REWARD_CURRENT_DEFAULT_STYLE
	}

	if points >= 100 {
		return TASK_REWARD_HIGH_STYLE
	}
	if points > 10 {
		return TASK_REWARD_MEDIUM_STYLE
	}
	return TASK_REWARD_DEFAULT_STYLE
}

func ProjectsTaskStyle(isActive bool, isSelected bool) ui.Style {
	switch {
	case isActive && isSelected:
		return PROJECTS_TASK_ACTIVE_SELECTED_STYLE
	case isActive:
		return PROJECTS_TASK_ACTIVE_STYLE
	case isSelected:
		return PROJECTS_TASK_SELECTED_STYLE
	default:
		return PROJECTS_TASK_STYLE
	}
}

func ProjectSidebarStyle(isSelected bool) ui.Style {
	if isSelected {
		return PROJECT_SIDEBAR_SELECTED_STYLE
	}
	return PROJECT_SIDEBAR_STYLE
}

func ModalProjectLineStyle(isEnabled bool, isSelected bool) ui.Style {
	switch {
	case isEnabled && isSelected:
		return MODAL_PROJECT_SELECTED_ENABLED_STYLE
	case isSelected:
		return MODAL_PROJECT_SELECTED_DISABLED_STYLE
	case isEnabled:
		return MODAL_PROJECT_ENABLED_STYLE
	default:
		return MODAL_PROJECT_DISABLED_STYLE
	}
}
