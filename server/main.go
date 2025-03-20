package main

import (
	"github.com/XORbit01/retro/config"
	"github.com/XORbit01/retro/server/player"
)

func main() {
	// load config
	cfg := config.GetConfig()

	player.StartIPCServer(cfg.ServerPort)
}
