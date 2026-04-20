package ui

import (
	"fmt"
	"log"
	"github.com/JbrownWFU/Radio-Go/internal/storage"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LikedSongsViewModel struct {
	table            table.Model
	db               *storage.DBConn
	songs            []storage.LikedSong
	showingDetails   bool
	selectedSong     storage.LikedSong
	confirmingDelete bool
	width            int
	height           int
}

func NewLikedSongsViewModel(db *storage.DBConn) LikedSongsViewModel {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Artist", Width: 20},
		{Title: "Title", Width: 20},
		{Title: "Station", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Selected = TableSelectedItemStyle
	t.SetStyles(s)

	m := LikedSongsViewModel{
		table: t,
		db:    db,
	}
	m.RefreshTable()
	return m
}

func (m *LikedSongsViewModel) RefreshTable() {
	songs, err := m.db.GetLikedSongs()
	if err != nil {
		log.Println("Error fetching liked songs:", err)
		return
	}

	m.songs = songs

	var rows []table.Row
	for _, s := range songs {
		rows = append(rows, table.Row{fmt.Sprintf("%d", s.ID), s.Artist, s.Title, s.StationName})
	}
	m.table.SetRows(rows)
}

func (m *LikedSongsViewModel) Init() tea.Cmd {
	return nil
}

func (m LikedSongsViewModel) Update(msg tea.Msg) (LikedSongsViewModel, tea.Cmd) {
	if m.showingDetails {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "esc":
				m.showingDetails = false
			}
		}
		return m, nil
	}

	if m.confirmingDelete {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				id := m.table.SelectedRow()[0]
				if err := m.db.DeleteLikedSong(id); err != nil {
					log.Println("Error deleting liked song:", err)
				}
				m.RefreshTable()
				m.confirmingDelete = false
			case "esc":
				m.confirmingDelete = false
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(m.table.Rows()) > 0 {
				row := m.table.SelectedRow()
				if len(row) > 0 {
					id := row[0]
					for _, s := range m.songs {
						if fmt.Sprintf("%d", s.ID) == id {
							m.selectedSong = s
							m.showingDetails = true
							break
						}
					}
				}
			}
		case "d":
			if len(m.table.Rows()) > 0 {
				m.confirmingDelete = true
			}
		case "c":
			if len(m.table.Rows()) > 0 {
				row := m.table.SelectedRow()
				artist := row[1]
				title := row[2]
				clipboard.WriteAll(fmt.Sprintf("%s - %s", artist, title))
			}
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m LikedSongsViewModel) View() string {
	if m.showingDetails {
		content := lipgloss.JoinVertical(lipgloss.Center,
			QuestionStyle.Render("Song Details"),
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color(White)).Bold(true).Render(m.selectedSong.Artist),
			lipgloss.NewStyle().Foreground(lipgloss.Color(Lavender)).Render(m.selectedSong.Title),
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color(Gray)).Render(fmt.Sprintf("Station: %s", m.selectedSong.StationName)),
			lipgloss.NewStyle().Foreground(lipgloss.Color(Gray)).Render(fmt.Sprintf("Liked: %s", m.selectedSong.LikedAt)),
			"",
			"(Press Enter/Esc to close)",
		)
		modal := ModalStyle.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	if m.confirmingDelete {
		selectedName := fmt.Sprintf("%s - %s", m.table.SelectedRow()[1], m.table.SelectedRow()[2])
		content := lipgloss.JoinVertical(lipgloss.Center,
			QuestionStyle.Render("Remove Liked Song?"),
			"",
			fmt.Sprintf("Are you sure you want to remove %q?", selectedName),
			"",
			"(Enter to confirm, Esc to cancel)",
		)
		modal := ModalStyle.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	return m.table.View()
}

func (m *LikedSongsViewModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.table.SetHeight(h - 4)
	m.table.SetWidth(w - 4)

	idW := 4
	totalW := w - 10 - idW
	artistW := totalW / 3
	titleW := totalW / 3
	stationW := totalW - artistW - titleW

	m.table.SetColumns([]table.Column{
		{Title: "ID", Width: idW},
		{Title: "Artist", Width: artistW},
		{Title: "Title", Width: titleW},
		{Title: "Station", Width: stationW},
	})
}
