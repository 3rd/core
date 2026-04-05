package components

import (
	"core/ui/task_interactive/state"
	"core/ui/task_interactive/theme"
	"fmt"

	ui "github.com/3rd/go-futui"
)

type keybind struct {
	Key  string
	Desc string
}

type keybindSection struct {
	Title    string
	Bindings []keybind
}

var globalBindings = []keybind{
	{"?", "toggle help"},
	{"q", "quit"},
	{"1/2/3", "switch tab"},
	{"g/G", "top / bottom"},
	{"y", "yank task text"},
}

var activeBindings = []keybind{
	{"j/k", "navigate"},
	{"Space", "toggle in-progress"},
	{"f", "focus project"},
	{"p", "filter projects"},
	{"t", "time filter (today / 24h)"},
	{".", "show/hide done tasks"},
	{"Ctrl+X", "deactivate task"},
	{"Enter", "edit in nvim"},
	{"Ctrl+Space", "toggle done"},
}

var historyBindings = []keybind{
	{"j/d", "scroll down"},
	{"k/u/s", "scroll up"},
}

var projectsBindings = []keybind{
	{"j/k", "navigate tasks"},
	{"J/K", "navigate projects"},
	{"Space", "toggle active"},
	{"Tab/S-Tab", "next / prev project"},
	{"Enter", "edit in nvim"},
}

type HelpModal struct {
	ui.Component
	AppState *state.AppState
	Width    int
	Height   int
}

func (c *HelpModal) getSections() []keybindSection {
	sections := []keybindSection{}

	switch c.AppState.CurrentTab {
	case state.APP_TAB_ACTIVE:
		sections = append(sections, keybindSection{Title: "Active", Bindings: activeBindings})
	case state.APP_TAB_HISTORY:
		sections = append(sections, keybindSection{Title: "History", Bindings: historyBindings})
	case state.APP_TAB_PROJECTS:
		sections = append(sections, keybindSection{Title: "Projects", Bindings: projectsBindings})
	}

	sections = append(sections, keybindSection{Title: "Global", Bindings: globalBindings})
	return sections
}

func (c *HelpModal) Render() ui.Buffer {
	b := ui.Buffer{}
	modalState := &c.AppState.HelpModal

	if !modalState.IsVisible {
		return b
	}

	sections := c.getSections()

	// count total lines: section title + bindings + blank separator per section
	totalLines := 0
	for i, s := range sections {
		totalLines += 1 + len(s.Bindings)
		if i < len(sections)-1 {
			totalLines++
		}
	}

	modalWidth := 48
	modalHeight := min(totalLines+5, c.Height-4) // borders + title + help + padding
	if modalWidth > c.Width-4 {
		modalWidth = c.Width - 4
	}

	b.Resize(modalWidth, modalHeight)
	b.FillStyle(theme.MODAL_STYLE)

	// borders
	borderStyle := theme.MODAL_BORDER_STYLE
	b.DrawCell(0, 0, '┌', borderStyle)
	for x := 1; x < modalWidth-1; x++ {
		b.DrawCell(x, 0, '─', borderStyle)
	}
	b.DrawCell(modalWidth-1, 0, '┐', borderStyle)
	for y := 1; y < modalHeight-1; y++ {
		b.DrawCell(0, y, '│', borderStyle)
		b.DrawCell(modalWidth-1, y, '│', borderStyle)
	}
	b.DrawCell(0, modalHeight-1, '└', borderStyle)
	for x := 1; x < modalWidth-1; x++ {
		b.DrawCell(x, modalHeight-1, '─', borderStyle)
	}
	b.DrawCell(modalWidth-1, modalHeight-1, '┘', borderStyle)

	// title
	title := " Keybindings "
	titleX := (modalWidth - len(title)) / 2
	b.Text(titleX, 0, title, theme.MODAL_TITLE_STYLE)

	// content area
	contentStartY := 2
	maxVisibleLines := modalHeight - 4 // borders top/bottom + title row + help row

	// build flat line list
	type line struct {
		text      string
		style     ui.Style
		isBinding bool
	}
	var lines []line

	sectionTitleStyle := theme.HELP_SECTION_STYLE
	bindingStyle := theme.HELP_BINDING_STYLE
	keyStyle := theme.HELP_KEY_STYLE

	for i, s := range sections {
		lines = append(lines, line{text: s.Title, style: sectionTitleStyle})
		for _, bind := range s.Bindings {
			// pad key to fixed width for alignment
			keyPadded := fmt.Sprintf("%-14s", bind.Key)
			lines = append(lines, line{
				text:      keyPadded + bind.Desc,
				style:     bindingStyle,
				isBinding: true,
			})
		}
		if i < len(sections)-1 {
			lines = append(lines, line{text: "", style: bindingStyle})
		}
	}

	scrollOffset := modalState.ScrollOffset
	maxScroll := max(len(lines)-maxVisibleLines, 0)
	if scrollOffset > maxScroll {
		scrollOffset = maxScroll
		modalState.ScrollOffset = scrollOffset
	}

	for i := scrollOffset; i < len(lines); i++ {
		row := contentStartY + i - scrollOffset
		if row >= modalHeight-2 {
			break
		}
		l := lines[i]
		if l.text == "" {
			continue
		}

		// for binding lines, render key portion in accent color
		if l.isBinding && len(l.text) > 14 {
			b.Text(3, row, l.text[:14], keyStyle)
			b.Text(17, row, l.text[14:], bindingStyle)
		} else {
			b.Text(3, row, l.text, l.style)
		}
	}

	// help text
	helpText := "j/k: scroll | ?/q/esc: close"
	helpY := modalHeight - 2
	helpX := max((modalWidth-len(helpText))/2, 2)
	b.Text(helpX, helpY, helpText, theme.MODAL_HELP_STYLE)

	return b
}
