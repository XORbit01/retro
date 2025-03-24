package views

import (
	"net/rpc"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/XORbit01/retro/client/controller"
	"github.com/XORbit01/retro/shared"
)

type AddToPlaylistCommand struct {
	Name string
}

func (c AddToPlaylistCommand) QuitMessage() string {
	return "🔋 Adding music to playlist playlist"
}

func (c AddToPlaylistCommand) Execute(query string, knowing shared.DResults, client *rpc.Client) ([]shared.SearchResult, error) {
	return controller.DetectAndAddToPlayList(shared.AddToPlayListQuery{
		PlayListName: c.Name,
		Knowing:      knowing,
		Query:        query,
	}, client)
}

func SearchThenAddToPlayList(playlist, query string, client *rpc.Client) error {
	model := NewRootModel(client, query, AddToPlaylistCommand{playlist})
	model.processingText = "Adding song to playlist..."
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
