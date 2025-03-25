package views

import (
	"fmt"
	"net/rpc"

	"github.com/XORbit01/retro/shared"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type searchResultItem struct {
	title    string
	desc     string
	ftype    shared.DResults
	duration string
}

func (i searchResultItem) Title() string {
	if i.ftype == shared.DCache {
		return i.title
	}
	return i.title
}

func (i searchResultItem) Description() string {
	return emojisType[i.ftype] + " " + string(i.ftype) + " " + i.duration
}

func (i searchResultItem) FilterValue() string { return "" }

type DetectCommand interface {
	Execute(query string, knowing shared.DResults, client *rpc.Client) ([]shared.SearchResult, error)
	QuitMessage() string
}

type RootModel struct {
	state  AppState
	client *rpc.Client
	query  string

	searchView     spinner.Model
	selectView     list.Model
	processingText string
	selectedItem   list.Item
	err            error

	commander DetectCommand
}

func NewRootModel(client *rpc.Client, query string, commander DetectCommand) RootModel {
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = GetTheme().SpinnerStyle

	return RootModel{
		state:      StateSearching,
		client:     client,
		query:      query,
		searchView: spin,
		commander:  commander,
	}
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(
		m.searchView.Tick,
		m.performSearch(m.query),
	)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateSearching:
		var cmd tea.Cmd
		m.searchView, cmd = m.searchView.Update(msg)

		switch msg := msg.(type) {
		case SearchFinishedMsg:
			if msg.Err != nil || len(msg.Results) == 0 {
				m.err = msg.Err
				m.state = StateQuitting
				return m, tea.Quit
			}
			m.selectView = list.New(msg.Results, GetTheme().ListDelegate, 50, 14)
			m.selectView.Title = "Select a song 👇"
			m.selectView.Styles.Title = GetTheme().SelectMusicsTitleStyle
			m.selectView.SetFilteringEnabled(false)
			m.selectView.SetShowHelp(false)
			m.state = StateSelecting
			return m, nil
		}
		return m, cmd

	case StateSelecting:
		var cmd tea.Cmd
		m.selectView, cmd = m.selectView.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				selected := m.selectView.SelectedItem()
				m.selectedItem = selected
				m.processingText = "Processing..."
				m.state = StateProcessing
				callback := func() tea.Msg {
					_, err := m.commander.Execute(selected.(searchResultItem).desc, selected.(searchResultItem).ftype, m.client)
					if err != nil {
						return CallbackFinishedMsg{err: err}
					}
					return CallbackFinishedMsg{err: nil}
				}
				return m, tea.Batch(m.searchView.Tick, callback)
			}
		}
		return m, cmd

	case StateProcessing:
		var cmd tea.Cmd
		m.searchView, cmd = m.searchView.Update(msg)
		switch msg := msg.(type) {
		case CallbackFinishedMsg:
			if msg.err != nil {
				m.err = msg.err
				return m, tea.Quit
			}
			m.state = StateQuitting
			return m, tea.Quit
		}
		return m, cmd
	}

	return m, nil
}

func (m RootModel) View() string {
	switch m.state {
	case StateSearching:
		return fmt.Sprintf("%s Searching for %q...", m.searchView.View(), m.query)
	case StateSelecting:
		return GetTheme().DocStyle.Render(m.selectView.View())
	case StateProcessing:
		return fmt.Sprintf("%s %s", m.searchView.View(), m.processingText)
	case StateQuitting:
		if m.err != nil {
			return GetTheme().FailStyle.Render(failedEmoji, m.err.Error())
		}
		if m.selectedItem != nil && m.commander.QuitMessage() != "" {
			return GetTheme().QuitTextStyle.Render(m.commander.QuitMessage())
		}
		return "Done!"

	default:
		return "Unknown state"
	}
}

func (m RootModel) performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		musics, err := m.commander.Execute(query, shared.DUnknown, m.client)
		if err != nil {
			return SearchFinishedMsg{Err: err}
		}
		if len(musics) == 0 {
			return SearchFinishedMsg{Results: nil} // triggers quit in Update
		}
		var items []list.Item
		for _, music := range musics {
			items = append(items, searchResultItem{
				title:    music.Title,
				desc:     music.Destination,
				ftype:    music.Type,
				duration: shared.DurationToString(music.Duration),
			})
		}
		return SearchFinishedMsg{Results: items}
	}
}
