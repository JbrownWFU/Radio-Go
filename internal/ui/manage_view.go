package ui

import (
	"fmt"
	"log"
	"github.com/JbrownWFU/Radio-Go/internal/player"
	"github.com/JbrownWFU/Radio-Go/internal/storage"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type MetadataMsg player.Metadata

type ManageViewModel struct {
	table            table.Model
	db               *storage.DBConn
	player           *player.Player
	activeStation    storage.Station
	form             *huh.Form
	addingNew        bool
	confirmingDelete bool
	width            int
	height           int
}

func NewManageViewModel(db *storage.DBConn, p *player.Player) ManageViewModel {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 20},
		{Title: "URL", Width: 40},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Selected = TableSelectedItemStyle
	t.SetStyles(s)

	m := ManageViewModel{
		table:  t,
		db:     db,
		player: p,
	}
	m.refreshTable()
	return m
}

func (m *ManageViewModel) refreshTable() {
	stations, err := m.db.GetStations()
	if err != nil {
		log.Println("Error fetching stations:", err)
		return
	}

	var rows []table.Row
	for _, s := range stations {
		rows = append(rows, table.Row{fmt.Sprintf("%d", s.ID), s.Name, s.Url})
	}
	m.table.SetRows(rows)
}

func (m *ManageViewModel) Init() tea.Cmd {
	return nil
}

func (m ManageViewModel) Update(msg tea.Msg) (ManageViewModel, tea.Cmd) {
	var cmds []tea.Cmd

	if m.addingNew {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		cmds = append(cmds, cmd)

		if m.form.State == huh.StateCompleted {
			name := m.form.GetString("name")
			url := m.form.GetString("url")
			if name != "" && url != "" {
				if err := m.db.InsertStation(name, url); err != nil {
					log.Println("Error inserting station:", err)
				}
				m.refreshTable()
			}
			m.addingNew = false
			m.table.Focus()
		}
		return m, tea.Batch(cmds...)
	}

	if m.confirmingDelete {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				row := m.table.SelectedRow()
				if len(row) > 0 {
					id := row[0]
					if err := m.db.DeleteStation(id); err != nil {
						log.Println("Error deleting station:", err)
					}
					m.refreshTable()
				}
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
				if len(row) < 3 {
					return m, nil
				}
				name := row[1]
				url := row[2]

				m.activeStation = storage.Station{Name: name, Url: url}
				m.player.Stop()

				return m, func() tea.Msg {
					info, err := m.player.Start(url)
					if err != nil {
						log.Println("Error starting player:", err)
						return nil
					}
					return MetadataMsg{StreamTitle: info.Name}
				}
			}
		case "s":
			m.player.Stop()
			m.activeStation = storage.Station{}
			return m, func() tea.Msg { return MetadataMsg{StreamTitle: ""} }
		case "a":
			m.addingNew = true
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Station Name").Key("name"),
					huh.NewInput().Title("URL").Key("url"),
				),
			)
			m.table.Blur()
			cmds = append(cmds, m.form.Init())
		case "d":
			if len(m.table.Rows()) > 0 {
				m.confirmingDelete = true
			}
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ManageViewModel) ActiveStationName() string {
	return m.activeStation.Name
}

func (m ManageViewModel) View() string {
	if m.addingNew {
		return m.form.View()
	}

	if m.confirmingDelete {
		row := m.table.SelectedRow()
		selectedName := "Unknown"
		if len(row) > 1 {
			selectedName = row[1]
		}
		content := lipgloss.JoinVertical(lipgloss.Center,
			QuestionStyle.Render("Delete Station?"),
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

func (m *ManageViewModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.table.SetHeight(h - 4)
	m.table.SetWidth(w - 4)
	
	// Resize columns
	idW := 4
	nameW := (w - 10 - idW) / 3
	urlW := w - 10 - idW - nameW
	m.table.SetColumns([]table.Column{
		{Title: "ID", Width: idW},
		{Title: "Name", Width: nameW},
		{Title: "URL", Width: urlW},
	})
}
