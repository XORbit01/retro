package views

import (
	"fmt"
	"github.com/XORbit01/retro/client/controller"
	"github.com/XORbit01/retro/shared"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"net/rpc"
)

type searchResultItem struct {
	title    string
	desc     string
	ftype    string
	duration string
}

func (i searchResultItem) Title() string {
	if i.ftype == "cache" {
		return i.title
	}
	return i.title
}

func (i searchResultItem) Description() string {
	return emojisType[i.ftype] + " " + i.ftype + " " + i.duration
}

func (i searchResultItem) FilterValue() string { return "" }

type RootModel struct {
	state  AppState
	client *rpc.Client
	query  string

	searchView     spinner.Model
	selectView     list.Model
	processingText string
	selectedItem   list.Item
	err            error

	runCallback func(item list.Item, client *rpc.Client) tea.Cmd
	quitMessage func(item list.Item) string
}

func NewRootModel(client *rpc.Client, query string) RootModel {
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = GetTheme().SpinnerStyle

	return RootModel{
		state:       StateSearching,
		client:      client,
		query:       query,
		searchView:  spin,
		runCallback: defaultCallback,
		quitMessage: defaultQuitMessage,
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
				m.err = msg.Err // if err is nil but results empty, it's still a no-op
				m.state = StateQuitting
				return m, tea.Quit
			}
			m.selectView = list.New(msg.Results, GetTheme().ListDelegate, 50, 14)
			m.selectView.Title = "Select a song 👇"
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
				return m, tea.Batch(m.searchView.Tick, m.runCallback(selected, m.client))
			}
		}
		return m, cmd

	case StateProcessing:
		var cmd tea.Cmd
		m.searchView, cmd = m.searchView.Update(msg)
		switch msg.(type) {
		case CallbackFinishedMsg:
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
			return fmt.Sprintf("An error occurred: %v", m.err)
		}
		if m.selectedItem != nil && m.quitMessage != nil {
			return m.quitMessage(m.selectedItem)
		}

		return "Done!"
	default:
		return "Unknown state"
	}
}

func (m RootModel) performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		musics, err := controller.DetectAndPlay(query, m.client)
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

func defaultCallback(item list.Item, client *rpc.Client) tea.Cmd {
	return func() tea.Msg {
		desc := item.(searchResultItem).desc
		controller.DetectAndPlay(desc, client)
		return CallbackFinishedMsg{}
	}
}

func defaultQuitMessage(item list.Item) string {
	title := item.(searchResultItem).title
	return GetTheme().QuitTextStyle.Render("🎵 Playing song " + title + ", this may take a while if download is needed")
}
