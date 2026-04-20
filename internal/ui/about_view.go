package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type AboutViewModel struct {
	width  int
	height int
}

func NewAboutViewModel() AboutViewModel {
	return AboutViewModel{}
}

func (m AboutViewModel) Init() {
}

func (m AboutViewModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(Purple)).
		Bold(true).
		Underline(true).
		MarginBottom(1)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(Lavender)).
		MarginBottom(1)

	linkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(Pink)).
		Italic(true)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("About RadioGo v2"),
		descStyle.Render("A modern terminal radio player built with Go and Bubble Tea."),
		"",
		// titleStyle.Render("Instructions"),
		// descStyle.Render("• Tab: Switch between Home and About"),
		// descStyle.Render("• Enter: Play selected station"),
		// descStyle.Render("• S: Stop playback"),
		// descStyle.Render("• A: Add a new station"),
		// descStyle.Render("• D: Delete selected station"),
		// descStyle.Render("• L: Like the current song"),
		// descStyle.Render("• ?: Toggle help menu"),
		// descStyle.Render("• Q / Ctrl+C: Quit"),
		// "",
		titleStyle.Render("Source Code"),
		linkStyle.Render("https://github.com/JbrownWFU/Radio-Go"),
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m *AboutViewModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
