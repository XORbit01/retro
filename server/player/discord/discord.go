package discord

// rich presence
import (
	"github.com/hugolgst/rich-go/client"
)

const DISCORD_CLIENT_ID = "1208868446926807040"

type State string

const (
	Playing State = "Playing"
	Paused  State = "Paused"
)

var Activity = client.Activity{
	State:      "Playing",
	Details:    "",
	LargeImage: "retro",
	LargeText:  "retro",
	Buttons: []*client.Button{
		{
			Label: "Download Retro",
			Url:   "https://github.com/XORbit01/retro",
		},
	},
}

func initClient(music string) error {
	if err := client.Login(
		DISCORD_CLIENT_ID,
	); err != nil {
		return err
	}
	Activity.Details = music
	Activity.State = string(Playing)
	Activity.SmallImage = "play"
	Activity.SmallText = "Playing"
	if err := client.SetActivity(
		Activity,
	); err != nil {
		return err
	}
	return nil
}

func Stop() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	client.Logout()
	return err
}

func Start(music string) error {
	return initClient(music)
}

func Pause() error {
	Activity.State = string(Paused)
	Activity.SmallImage = "pause"
	Activity.SmallText = "Paused"
	return client.SetActivity(
		Activity,
	)
}
