package views

import (
	"fmt"
	"strings"

	"github.com/XORbit01/retro/shared"
)

// TreeListView is a reusable component for rendering lists in tree format.
type TreeListView struct {
	Title string             // e.g. "Playlists", "Cache", "Playlist: Lofi"
	Emoji string             // e.g. 📁, 🎧, 📼
	Items []shared.HashNamed // the actual list of items to render
}

// Render outputs the styled tree list using the theme from UIContext.
func (v TreeListView) Render(ctx UIContext) error {
	// Print title
	fmt.Println(ctx.Theme.PositionStyle.Render(fmt.Sprintf("%s %s", v.Emoji, v.Title)))
	fmt.Println()

	// Print each item in tree structure
	for i, item := range v.Items {
		prefix := "├──["
		if i == len(v.Items)-1 {
			prefix = "└──["
		}
		fmt.Print(ctx.Theme.PositionStyle.Copy().Inherit(ctx.Theme.ColoredTextStyle).Render(prefix))
		fmt.Print(item.Hash[:shared.HashPrefixLength])

		fmt.Print(ctx.Theme.ColoredTextStyle.Render("] "))
		fmt.Println(strings.TrimSpace(item.Name))
	}

	return nil
}
