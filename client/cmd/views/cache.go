package views

import (
	"github.com/XORbit01/retro/shared"
)

type CacheView struct {
	Songs []shared.CacheItem
}

func (v CacheView) Render(ctx UIContext) error {
	tree := TreeListView{
		Title: "Cache",
		Emoji: "📁",
		Items: v.Songs,
	}
	return tree.Render(ctx)
}
