package views

import (
	"net/rpc"

	"github.com/XORbit01/retro/client/controller"
)

// Handler builds a fully prepared view ready to be rendered.
type Handler interface {
	BuildView(client *rpc.Client) View
}

// ===== Concrete Handlers =====

// CacheHandler creates a TreeListView showing cached songs.
type CacheHandler struct{}

func (h CacheHandler) BuildView(client *rpc.Client) View {
	songs := controller.GetCachedMusics(client)
	return TreeListView{
		Title: "Cache",
		Emoji: "📁",
		Items: songs,
	}
}

// PlaylistsHandler creates a TreeListView showing playlist names.
type PlaylistsHandler struct{}

func (h PlaylistsHandler) BuildView(client *rpc.Client) View {
	names := controller.GetPlayListsMeta(client)
	return TreeListView{
		Title: "Playlists",
		Emoji: "📼",
		Items: names,
	}
}

// PlaylistSongsHandler creates a TreeListView showing songs in a playlist.
type PlaylistSongsHandler struct {
	Name string
}

func (h PlaylistSongsHandler) BuildView(client *rpc.Client) View {
	songs := controller.GetPlayListMusicsMeta(h.Name, client)
	return TreeListView{
		Title: "Playlist: " + h.Name,
		Emoji: "🎧",
		Items: songs,
	}
}

func NewCacheHandler() Handler {
	return CacheHandler{}
}

func NewPlaylistsHandler() Handler {
	return PlaylistsHandler{}
}

func NewPlaylistSongsHandler(name string) Handler {
	return PlaylistSongsHandler{Name: name}
}
