package views

import (
	"fmt"
	"net/rpc"

	"github.com/XORbit01/retro/shared"

	"github.com/XORbit01/retro/client/controller"
)

func CacheDisplay(client *rpc.Client) {
	songs := controller.GetCachedMusics(client)
	if len(songs) == 0 {
		fmt.Println("No music in cache")
		return
	}
	fmt.Print(GetTheme().PositionStyle.Render("📁 Cache\n"))
	fmt.Print("\n")

	for l, _ := range songs {
		if l == len(songs)-1 {
			fmt.Print(GetTheme().PositionStyle.Copy().Inherit(GetTheme().ColoredTextStyle).Render("└──["))
		} else {
			fmt.Print(GetTheme().PositionStyle.Copy().Inherit(GetTheme().ColoredTextStyle).Render("├──["))
		}
		fmt.Print(songs[l].Hash[:shared.HashPrefixLength])
		fmt.Print(GetTheme().ColoredTextStyle.Render("] "))
		fmt.Print(songs[l].Name)
	}

}
