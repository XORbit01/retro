package views

import (
	"github.com/XORbit01/retro/client/controller"
	"github.com/XORbit01/retro/shared"
	tea "github.com/charmbracelet/bubbletea"
	"net/rpc"
)

type PlayCommand struct {
}

func (c PlayCommand) QuitMessage() string {
	return "🎵 Playing song this may take a while if download is needed"
}

func (c PlayCommand) Execute(query string, knowing shared.DResults, client *rpc.Client) ([]shared.SearchResult, error) {
	return controller.DetectAndPlay(shared.DetectQuery{
		Query:   query,
		Knowing: knowing,
	}, client)
}

type PlaySongView struct {
	Query string
}

func NewPlaySongView(query string) *PlaySongView {
	return &PlaySongView{Query: query}
}

func (v PlaySongView) Run(ctx UIContext) error {
	p := tea.NewProgram(NewRootModel(ctx.Client, v.Query, PlayCommand{}))
	_, err := p.Run()
	return err
}
