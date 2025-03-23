package views

import (
	"github.com/XORbit01/retro/client/controller"
	"github.com/XORbit01/retro/shared"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
)

type PlaySongView struct {
	Query string
}

func (v PlaySongView) Run(ctx UIContext) error {
	model := NewModel(ctx.Client, v.Query)
	model.callback = playCallback
	model.quitMessage = PlayQuitMessage
	model.initCmd = model.PlaySearch

	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}

func playCallback(m model) error {
	i := m.selectList.Index()
	_, err := controller.DetectAndPlay(m.selectList.Items()[i].(searchResultItem).desc, m.client)
	return err
}

func PlayQuitMessage(m model) string {
	randEmoji := playingEmojis[rand.Intn(len(playingEmojis))]
	return GetTheme().QuitTextStyle.Render(
		randEmoji + " Playing song " + m.selectList.Items()[m.selectList.Index()].(searchResultItem).title + ", this may take a while if download needed",
	)
}

func (m model) PlaySearch() tea.Msg {
	var results []list.Item
	musics, err := controller.DetectAndPlay(m.query, m.client)
	if err != nil {
		return searchDone{nil, err}
	}
	for _, music := range musics {
		results = append(results, searchResultItem{
			title:    music.Title,
			desc:     music.Destination,
			ftype:    music.Type,
			duration: shared.DurationToString(music.Duration),
		})
	}
	return searchDone{
		results: results,
	}
}
