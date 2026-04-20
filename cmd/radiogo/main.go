package main

import (
	"fmt"
	"log"
	"os"
	"github.com/JbrownWFU/Radio-Go/internal/player"
	"github.com/JbrownWFU/Radio-Go/internal/storage"
	"github.com/JbrownWFU/Radio-Go/internal/ui"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateHome sessionState = iota
	stateLikes
	stateAbout
)

// KeyMap defines the global keybindings.
type keyMap struct {
	Tab   key.Binding
	Quit  key.Binding
	Help  key.Binding
	Add   key.Binding
	Del   key.Binding
	Like  key.Binding
	Play  key.Binding
	Stop  key.Binding
	Copy  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.Play, k.Stop},
		{k.Add, k.Del, k.Like, k.Copy},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch tab"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q", "esc"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add station"),
	),
	Del: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Like: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "like song"),
	),
	Play: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "play/select"),
	),
	Stop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy"),
	),
}

type mainModel struct {
	state          sessionState
	db             *storage.DBConn
	player         *player.Player
	homeView       ui.ManageViewModel
	likesView      ui.LikedSongsViewModel
	aboutView      ui.AboutViewModel
	help           help.Model
	keys           keyMap
	width          int
	height         int
	metadata       player.Metadata
	notification   string
	confirmingExit bool
}

type clearNotificationMsg struct{}

func initialModel() mainModel {
	dbPath := "./store.db"
	db, err := storage.InitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	p := player.New()

	m := mainModel{
		state:     stateHome,
		db:        db,
		player:    p,
		homeView:  ui.NewManageViewModel(db, p),
		likesView: ui.NewLikedSongsViewModel(db),
		aboutView: ui.NewAboutViewModel(),
		help:      help.New(),
		keys:      keys,
	}
	return m
}

// waitForMetadata is a command that waits for metadata from the player.
func (m mainModel) waitForMetadata() tea.Cmd {
	return func() tea.Msg {
		return ui.MetadataMsg(<-m.player.MetaChan)
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		m.homeView.Init(),
		m.likesView.Init(),
		m.waitForMetadata(),
	)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		contentH := m.height - 12
		contentW := m.width - 4

		m.homeView.SetSize(contentW, contentH)
		m.likesView.SetSize(contentW, contentH)
		m.aboutView.SetSize(contentW, contentH)

	case tea.KeyMsg:
		if m.confirmingExit {
			switch msg.String() {
			case "y", "enter":
				m.player.Stop()
				return m, tea.Quit
			case "n", "esc":
				m.confirmingExit = false
				return m, nil
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Quit) {
			m.confirmingExit = true
			return m, nil
		}

		if key.Matches(msg, m.keys.Tab) {
			m.state = (m.state + 1) % 3
			if m.state == stateLikes {
				m.likesView.RefreshTable()
			}
		}
		if key.Matches(msg, m.keys.Help) {
			m.help.ShowAll = !m.help.ShowAll
		}

		if key.Matches(msg, m.keys.Copy) && m.state == stateLikes {
			m.notification = "Copied to Clipboard!"
			cmds = append(cmds, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
				return clearNotificationMsg{}
			}))
		}

		if key.Matches(msg, m.keys.Like) && m.metadata.StreamTitle != "" {
			artist := m.metadata.Artist
			title := m.metadata.Title
			m.db.LikeSong(artist, title, "Radio")
			m.notification = "Song Liked!"
			cmds = append(cmds, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
				return clearNotificationMsg{}
			}))
		}

	case clearNotificationMsg:
		m.notification = ""

	case ui.MetadataMsg:
		m.metadata = player.Metadata(msg)
		cmds = append(cmds, m.waitForMetadata())
	}

	// Route updates to sub-views
	var cmd tea.Cmd
	switch m.state {
	case stateHome:
		m.homeView, cmd = m.homeView.Update(msg)
	case stateLikes:
		m.likesView, cmd = m.likesView.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	// Header
	tabs := []string{"Home", "Likes", "About"}
	var renderedTabs []string
	for i, t := range tabs {
		style := ui.NavBarInactiveTabStyle
		if (m.state == stateHome && i == 0) || (m.state == stateLikes && i == 1) || (m.state == stateAbout && i == 2) {
			style = ui.NavBarActiveTabStyle
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}
	header := ui.NavBarStyle.Width(m.width - 2).Render(lipgloss.JoinHorizontal(lipgloss.Left, renderedTabs...))

	// Content
	var content string
	switch m.state {
	case stateHome:
		content = m.homeView.View()
	case stateLikes:
		content = m.likesView.View()
	case stateAbout:
		content = m.aboutView.View()
	}
	mainView := ui.SubViewStyle.Width(m.width - 2).Height(m.height - 12).Render(content)

	// Footer: Station
	stationName := m.homeView.ActiveStationName()
	if stationName == "" {
		stationName = "None"
	}
	stationLabel := ui.StationLabelStyle.Render("STATION")
	stationContent := ui.StationContentStyle.Render(stationName)
	stationBlock := lipgloss.JoinHorizontal(lipgloss.Center, stationLabel, stationContent)

	// Footer: Now Playing
	nowPlaying := ui.NowPlayingLabelStyle.Render("NOW PLAYING")
	meta := m.metadata.StreamTitle
	if meta == "" {
		meta = "Stopped"
	}
	
	footerContent := ui.NowPlayingContentStyle.Render(meta)
	if m.notification != "" {
		footerContent = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.White)).
			Background(lipgloss.Color(ui.Pink)).
			Padding(0, 1).
			Render(m.notification)
	}

	nowPlayingBar := ui.FooterStyle.Width(m.width - 2).Render(
		lipgloss.JoinHorizontal(lipgloss.Center, stationBlock, nowPlaying, footerContent),
	)

	// Footer: Help
	helpView := m.help.View(m.keys)

	mainContent := ui.MainStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			mainView,
			nowPlayingBar,
			helpView,
		),
	)

	if m.confirmingExit {
		modalContent := lipgloss.JoinVertical(lipgloss.Center,
			ui.QuestionStyle.Render("Exit RadioGo?"),
			"",
			"Are you sure you want to quit?",
			"",
			"(Y/Enter to confirm, N/Esc to cancel)",
		)
		modal := ui.ModalStyle.Render(modalContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	return mainContent
}

func main() {
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
