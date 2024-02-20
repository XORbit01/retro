package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Malwarize/goplay/client/controller"
)

var rootCmd = &cobra.Command{
	Use:   "goplay",
	Short: "goplay is a music player",
	Long: `goplay is a music player
goplay is client for goplayer server
you can controll goplayer server like any other systemd service
goplay [command] --help for more information about a command`,

	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("goplay is a music player")
		fmt.Println("use goplay --help to see available commands")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	_, err := controller.GetClient()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(resumeCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(prevCmd)
	rootCmd.AddCommand(seekCmd)
	rootCmd.AddCommand(seekBackCmd)
	rootCmd.AddCommand(volumeCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(setThemeCmd)
	rootCmd.AddCommand(playlistCmd)
	playlistCmd.AddCommand(playlistCreateCmd)
	playlistCmd.AddCommand(playlistRemoveCmd)
	playlistCmd.AddCommand(playlistAddCmd)
	playlistCmd.AddCommand(playlistPlayCmd)
}
