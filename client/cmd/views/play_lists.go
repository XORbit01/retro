package views

import (
	"fmt"
	"github.com/XORbit01/retro/client/controller"
)

type PlaylistsView struct{}

func (v PlaylistsView) Render(ctx UIContext) error {
	playlists := controller.GetPlayListsMeta(ctx.Client)

	return TreeListView{
		Title: "Playlists",
		Emoji: "📼",
		Items: playlists, // []shared.Playlist (alias of HashNamed)
	}.Render(ctx)
}

type PlaylistSongsView struct {
	Name string
}

func (v PlaylistSongsView) Render(ctx UIContext) error {
	songs := controller.GetPlayListMusicsMeta(v.Name, ctx.Client)

	return TreeListView{
		Title: fmt.Sprintf("Playlist: %s", v.Name),
		Emoji: "🎧",
		Items: songs, // []shared.Song (alias of HashNamed)
	}.Render(ctx)
}
