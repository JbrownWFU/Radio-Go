package main

import (
	"fmt"
	"os"
	"log"

	_ "modernc.org/sqlite"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

var keys = keyMap{
	Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "up")),
	Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "down")),
	Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}

type model struct {
	width int
	height int
	
	keys keyMap
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {

	m := model{
		keys: keys,
	}
	
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle window resizing
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Enter):
			log.Println("Enter")
			return m, nil
		}
	}
	
	return m, nil
	
}

func (m model) View() string {
	root := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("200"))
	
	body := lipgloss.Place(
		m.width - 2, m.height - 2,
		lipgloss.Center, lipgloss.Center,
		"Welcome!",
	)
	
	window := root.Render(body)
	
	return window
}

func main() {
	// Set up logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("Error creating log file: ", err)
		os.Exit(1)
	}
	defer f.Close()

	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}