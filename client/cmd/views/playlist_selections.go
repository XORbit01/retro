package views

import (
	"net/rpc"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/XORbit01/retro/client/controller"
	"github.com/XORbit01/retro/shared"
)

// SearchThenAddToPlayList runs a Bubble Tea program that
// searches for a song and adds it to a playlist upon selection.
func SearchThenAddToPlayList(playlist, query string, client *rpc.Client) error {
	musics, err := controller.DetectAndAddToPlayList(playlist, query, client)
	if err != nil {
		return err
	}
	if len(musics) == 0 {
		return nil // nothing to add
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

	model := NewRootModel(client, query)
	model.processingText = "Adding song to playlist..."

	model.runCallback = func(item list.Item, client *rpc.Client) tea.Cmd {
		return func() tea.Msg {
			desc := item.(searchResultItem).desc
			controller.DetectAndAddToPlayList(playlist, desc, client)
			return CallbackFinishedMsg{}
		}
	}

	model.quitMessage = func(item list.Item) string {
		return GetTheme().QuitTextStyle.Render(
			"🔋 Adding music " + item.(searchResultItem).title + " to playlist " + playlist,
		)
	}

	// inject pre-fetched items directly into the list
	model.state = StateSelecting
	model.selectView = list.New(items, GetTheme().ListDelegate, 50, 14)
	model.selectView.Title = "Select a song 👇"
	model.selectView.SetFilteringEnabled(false)
	model.selectView.SetShowHelp(false)

	p := tea.NewProgram(model)
	_, err = p.Run()
	return err
}
