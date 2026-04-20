package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	Purple     = "#7e1dfb"
	DarkPurple = "#211338"
	Pink       = "#b28ecc"
	Lavender   = "#e6def7"
	Gray       = "#808080"
	White      = "#FFFFFF"

	// Basic Styles
	MainStyle = lipgloss.NewStyle()

	SubViewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(Purple)).
			Padding(1, 1)

	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(Purple)).
			Padding(1, 2).
			Align(lipgloss.Center)

	QuestionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(Pink))

	// Navigation
	NavBarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(Purple)).
			Align(lipgloss.Center)

	NavBarActiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(Pink)).
				Underline(true).
				MarginRight(2)

	NavBarInactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(Pink)).
				MarginRight(2)

	NavBarUnfocusedInactiveTabStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color(Gray)).
					MarginRight(2)

	NavBarUnfocusedActiveTabStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color(Gray)).
					Underline(true).
					MarginRight(2)

	// Table Styles
	TableSelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(DarkPurple)).
				Background(lipgloss.Color(Pink))

	TableUnfocusedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(DarkPurple)).
				Background(lipgloss.Color(Gray))

	// Footer Styles
	FooterStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(Purple)).
			Padding(0, 1)

	NowPlayingLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(White)).
				Background(lipgloss.Color(Purple)).
				Padding(0, 1).
				MarginRight(1).
				Bold(true)

	NowPlayingContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(Lavender))

	StationLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(DarkPurple)).
				Background(lipgloss.Color(Pink)).
				Padding(0, 1).
				MarginRight(1).
				Bold(true)

	StationContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(Pink)).
				MarginRight(3)
)
