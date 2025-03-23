package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

type PlaySongView struct {
	Query string
}

func (v PlaySongView) Run(ctx UIContext) error {
	p := tea.NewProgram(NewRootModel(ctx.Client, v.Query))
	_, err := p.Run()
	return err
}
