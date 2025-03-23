package views

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"

	"github.com/XORbit01/retro/client/controller"
	"github.com/XORbit01/retro/shared"
)

func reformatDuration(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d", int(duration.Minutes()), int(duration.Seconds())%60)
}

func convertVolumeToEmoji(volume uint8) string {
	if volume == 0 {
		return volumeLevels[0]
	}
	if volume < 50 {
		return volumeLevels[1]
	}
	if volume < 85 {
		return volumeLevels[2]
	}
	return volumeLevels[3]
}

func DisplayStatus(client *rpc.Client) {
	status := controller.GetPlayerStatus(client)
	queue := status.MusicQueue
	if status.PlayerState == shared.Stopped {
		fmt.Println(GetTheme().StoppedStyle.Render(emojisStatus[shared.Stopped], " Stopped"))
	} else {

		currentMusicName := queue[status.CurrMusicIndex]

		currentPosition := status.CurrMusicPosition
		currentPositionStr := reformatDuration(currentPosition)

		totalDurationStr := reformatDuration(status.CurrMusicDuration)

		prog := progress.New(progress.WithSolidFill(GetTheme().MainColorStyle), progress.WithWidth(40))
		prog.SetPercent(0.5)
		prog.ShowPercentage = false
		fmt.Println(GetTheme().ProgressStyle.Render(prog.ViewAs(currentPosition.Seconds() / status.CurrMusicDuration.Seconds())))

		fmt.Println("   "+playingEmojis[rand.Intn(len(playingEmojis))], currentMusicName)
		fmt.Println(GetTheme().PositionStyle.Copy().Inherit(GetTheme().ColoredTextStyle).Render(currentPositionStr, " / ", totalDurationStr))

		switch status.PlayerState {
		case shared.Playing:
			fmt.Println(GetTheme().RunningStyle.Render(emojisStatus[shared.Playing], " Playing", convertVolumeToEmoji(status.Volume)))
		case shared.Paused:
			fmt.Println(GetTheme().PausedStyle.Render(emojisStatus[shared.Paused], " Paused"))
		}
		// display queue
		for i, music := range queue {
			if i == status.CurrMusicIndex {
				fmt.Println(GetTheme().SelectMusicStyle.Render("->", strconv.Itoa(i), ":", music))
			} else {
				fmt.Println("  ", i, ":", music)
			}
		}
	}
	// display tasks
	for target, task := range status.Tasks {
		if task.Error != "" {
			switch task.Type {
			case shared.Downloading:
				fmt.Println(
					GetTheme().FailStyle.Render(
						failedEmoji,
						"Failed to download ",
						target,
						":",
						task.Error,
					),
				)
			case shared.Searching:
				fmt.Println(
					GetTheme().FailStyle.Render(
						failedEmoji,
						"Failed to search ",
						target,
						":",
						task.Error,
					),
				)
			default:
				fmt.Println(
					GetTheme().FailStyle.Render(
						failedEmoji,
						"Failed to ",
						target,
						":",
						task.Error,
					),
				)
			}
			continue
		}
		switch task.Type {
		case shared.Downloading:
			fmt.Println(
				GetTheme().TaskStyle.Render(tasksEmojis[task.Type], "Downloading ", target),
			)
		case shared.Searching:
			fmt.Println(GetTheme().TaskStyle.Render(tasksEmojis[task.Type], "Searching ", target))
		}
	}
}
