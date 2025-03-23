package cmd

import (
	"fmt"
	"github.com/XORbit01/retro/client/cmd/views/selections"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/XORbit01/retro/client/cmd/views"
	"github.com/XORbit01/retro/client/controller"
	"github.com/XORbit01/retro/shared"
)

var client, _ = controller.GetClient()

var playCmd = &cobra.Command{
	Use:   "play [query]",
	Short: "play a song <query>",
	Long: `play a song <query>
	play is smart enough to play the song from the query, you don't have to specify the type of the query
	if you want to explicitly specify the type of query, use the flags (TODO: add explicit flags)
		- if the query is a directory, it will play all the songs in the directory
		- if the query is a playlist, it will play all the songs in the playlist
		- if the query is a audio file, it will play the audio file
		- if the query is a youtube link, it will play the audio from the link
		- if the query is a search query, it will search and return the results to select from
	`,
	ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		if client == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var options []string
		list := controller.GetPlayListsMeta(client)
		for _, song := range list {
			options = append(options, song.Name)
		}
		// songs in the queue
		options = append(options, controller.GetPlayerStatus(client).MusicQueue...)
		return options, cobra.ShellCompDirectiveDefault
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			song := strings.Join(args, " ")
			view := selections.PlaySongView{Query: song}
			ctx := views.UIContext{
				Client: client,
				Theme:  views.GetTheme(),
			}
			_ = view.Run(ctx)
		} else {
			fmt.Println("no song specified")
		}
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "pause the current song",
	Long: `this command will pause the current song if it's playing
very easy to use, just type "pause" and hit enter`,
	Run: func(_ *cobra.Command, _ []string) {
		controller.Pause(client)
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "resume the current song",
	Long: `this command will resume the current song if it's paused
very easy to use, just type "resume" and hit enter`,
	Run: func(_ *cobra.Command, _ []string) {
		controller.Resume(client)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop the current song",
	Long: `this command will stop the current song 
it will also clear the queue and reset the player and remove the tasks if any
`,
	Run: func(_ *cobra.Command, _ []string) {
		controller.Stop(client)
	},
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "play the next song",
	Long: `this command will play the next song in the queue
if there is no next song, it will do nothing if the queue is empty 
it will play the first song in the queue if the queue is not empty 
`,
	Run: func(_ *cobra.Command, _ []string) {
		controller.Next(client)
	},
}

var prevCmd = &cobra.Command{
	Use:   "prev",
	Short: "play the previous song",
	Long: `play the previous song
if there is no previous song, it will do nothing if the queue is empty
it will play the last song in the queue if the queue is not empty
`,
	Run: func(_ *cobra.Command, _ []string) {
		controller.Prev(client)
	},
}

var seekCmd = &cobra.Command{
	Use:   "seek [seconds]",
	Short: "seek to a position in the current song",
	Long: `seek to a position in the current song 
if you are seeking forward, use positive seconds
if you are seeking backward, use negative seconds with -- seconds
	seek -- 10
	seek 10
you can use seekback to "seek" backward instead of using negative seconds

it will seek 5 seconds toward the end of the song if no seconds provided
	seek
	seek 5
	seekback
	seekback 5
if the seek seconds is greater than the song duration, it will play the next song in the queue
if the seek seconds is less than the song duration, it will rewind to the beginning of the song
`,
	Run: func(_ *cobra.Command, args []string) {
		var seekSeconds int
		if len(args) > 0 {
			var err error
			seekSeconds, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			seekSeconds = 5
		}
		controller.Seek(seekSeconds, client)
	},
}

var seekBackCmd = &cobra.Command{
	Use:   "seekback [seconds]",
	Short: "seek back by a number of seconds",
	Long: `seek back by a number of seconds 
this command will seek back by a number of seconds
it will seek 5 seconds back if no seconds provided
	seekback
	seekback 5
if the seekback seconds is less than the song duration, it will rewind to the beginning of the song
if you are seeking forward check "seek" command
`,
	Run: func(_ *cobra.Command, args []string) {
		var seekSeconds int
		if len(args) > 0 {
			var err error
			seekSeconds, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			seekSeconds = 5
		}
		controller.Seek(-seekSeconds, client)
	},
}

var volumeCmd = &cobra.Command{
	Use:   "vol [percentage]",
	Short: "set the volume to a percentage",
	Args:  cobra.MinimumNArgs(1),
	Long: `set the volume to a percentage
this command will set the volume to a percentage between 0 and 100
`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) > 0 {
			vol, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			controller.Volume(uint8(vol), client)
		} else {
			fmt.Println("no volume specified")
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove <index> | <song name>",
	Short: "remove a song from the queue by index or name",
	Long: `remove a song from the queue
this command will remove a song from the queue
it accepts the index of the song in the queue or the name of the song
`,
	ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		if client == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		playerStatus := controller.GetPlayerStatus(client)

		names := make([]string, 0, len(playerStatus.MusicQueue))
		names = append(names, playerStatus.MusicQueue...)
		return names, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(_ *cobra.Command, args []string) {
		if len(args) > 0 {
			if v, err := strconv.Atoi(args[0]); err == nil {
				controller.Remove(
					shared.IntOrString{
						IntVal: v,
						IsInt:  true,
					},
					client,
				)
			} else {
				controller.Remove(
					shared.IntOrString{
						StrVal: strings.Join(args, " "),
						IsInt:  false,
					},
					client,
				)
			}
		} else {
			fmt.Println("no song specified")
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "get the current status of the player queue",
	Long: `get the current status of the player queue
this command will display the current status of the player queue
including the current song, the queue, the current position, the tasks, volume, and the volume level
you can change the theme of the status display using the "theme" command 
`,
	Run: func(_ *cobra.Command, _ []string) {
		views.DisplayStatus(client)
	},
}

var playlistCmd = &cobra.Command{
	Use:   "list",
	Short: "list playlists | list songs in a playlist",
	Long: `list playlists 
This command will list all the playlists.
If no playlist name is provided, it will list all playlists.
If a playlist name is provided, it will list all the songs in that playlist.
`,
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if client == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		listMetas := controller.GetPlayListsMeta(client)
		listNames := make([]string, len(listMetas))
		for _, listM := range listMetas {
			listNames = append(listNames, listM.Name)
		}

		if len(args) == 0 {
			return listNames, cobra.ShellCompDirectiveDefault
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(_ *cobra.Command, args []string) {
		var handler views.Handler

		if len(args) > 0 {
			listname := strings.TrimSpace(strings.Join(args, " "))
			handler = views.NewPlaylistSongsHandler(listname)
		} else {
			handler = views.NewPlaylistsHandler()
		}

		view := handler.BuildView(client)
		ctx := views.UIContext{
			Client: client,
			Theme:  views.GetTheme(),
		}

		_ = view.Render(ctx)
	},
}
var playlistCreateCmd = &cobra.Command{
	Use:   "create <playlist name>",
	Short: "create a new playlist",
	Long: `create a new playlist
this command will create a new playlist with the provided name
playlist stores the songs in path provided in the config file
default: $HOME/.retro/playlists
`,
	Run: func(_ *cobra.Command, args []string) {
		lists := controller.GetPlayListsMeta(client)
		if len(args) > 0 {
			name := strings.Join(args, " ")
			for _, list := range lists {
				if list.Name == name {
					fmt.Println("playlist already exist")
					return
				}
			}
			controller.CreatePlayList(name, client)
		} else {
			fmt.Println("no playlist name specified")
		}
	},
}

// remove
var playlistRemoveCmd = &cobra.Command{
	Use:   "remove <playlist> | <playlist> <song index>",
	Short: "remove a playlist (and its songs) | remove a song from a playlist",
	Long: `this command will remove a playlist (and its songs) | remove a song from a playlist
if no song index is provided, it will remove the playlist and its songs
if a song index is provided, it will remove the song from the playlist
`,
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if client == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		listMetas := controller.GetPlayListsMeta(client)
		listNames := make([]string, len(listMetas))
		for _, listM := range listMetas {
			listNames = append(listNames, listM.Name)
		}
		if len(args) == 0 {
			return listNames, cobra.ShellCompDirectiveDefault
		}
		if len(args) == 1 {
			songs := controller.GetPlayListMusicsMeta(args[0], client)
			parsedMusics := make([]string, 0, len(songs))
			for _, song := range songs {
				parsedMusics = append(parsedMusics, song.Name)
			}
			return parsedMusics, cobra.ShellCompDirectiveDefault
		}

		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(_ *cobra.Command, args []string) {
		listname := strings.TrimSpace(args[0])
		songs := controller.GetPlayListMusicsMeta(listname, client)
		if len(args) == 1 {
			controller.RemovePlayList(listname, client)
		} else if len(args) == 2 {
			// check if number provided is valid
			songIndex, err := strconv.Atoi(strings.TrimSpace(args[1]))
			if err == nil && songIndex >= 0 && songIndex < len(songs) {
				controller.RemoveMusicFromPlayList(
					listname,
					shared.IntOrString{
						IntVal: songIndex,
						IsInt:  true,
					},
					client,
				)
			} else {
				songName := strings.TrimSpace(args[1])
				controller.RemoveMusicFromPlayList(
					listname,
					shared.IntOrString{
						StrVal: songName,
					},
					client,
				)
			}
		} else {
			fmt.Println("playlist name required or playlist name and song index required")
		}
	},
}

// add song to a playlist
var playlistAddCmd = &cobra.Command{
	Use:   "add <listname> <query>",
	Short: "add music(s) to a playlist",
	Long: `add music(s) to a playlist
this command is similar to the "play" command, but it will add the music to the playlist instead of adding it to the queue
you can check the "list <playlist>" command to see the songs in the playlist
and you can play it using the "list play" command
`,
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if client == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		listMetas := controller.GetPlayListsMeta(client)
		listNames := make([]string, len(listMetas))
		for _, listM := range listMetas {
			listNames = append(listNames, listM.Name)
		}
		if len(args) == 0 {
			return listNames, cobra.ShellCompDirectiveDefault
		}
		if len(args) == 1 {
			// get music in queue
			musics := controller.GetPlayerStatus(client).MusicQueue

			return musics, cobra.ShellCompDirectiveDefault
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(_ *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("playlist name and query required")
			return
		}
		listName := strings.TrimSpace(args[0])
		query := strings.Join(args[1:], " ")
		err := views.SearchThenAddToPlayList(listName, query, client)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var playlistPlayCmd = &cobra.Command{
	Use:   "play <playlist> | <playlist> <song_name|index>",
	Short: "play a playlist | play a song from a playlist",
	Long: `play a playlist | play a song from a playlist
this command will play a playlist | play a song from a playlist
if no song name is provided, it will add the all the songs in the playlist to the queue
if a song name is provided, it will add the song to the queue
`,
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if client == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		listMetas := controller.GetPlayListsMeta(client)
		listNames := make([]string, len(listMetas))
		for _, listM := range listMetas {
			listNames = append(listNames, listM.Name)
		}

		if len(args) == 0 {
			return listNames, cobra.ShellCompDirectiveDefault
		}
		if len(args) == 1 {
			songs := controller.GetPlayListMusicsMeta(args[0], client)
			songsNames := make([]string, len(songs))
			for _, song := range songs {
				songsNames = append(songsNames, song.Name)
			}
			return songsNames, cobra.ShellCompDirectiveDefault
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 2 {
			lisname := args[0]
			songIndex, err := strconv.Atoi(args[1])
			if err != nil {
				controller.PlayListPlayMusic(
					lisname,
					shared.IntOrString{
						StrVal: args[1],
					},
					client,
				)
			} else {
				controller.PlayListPlayMusic(
					lisname,
					shared.IntOrString{
						IntVal: songIndex,
						IsInt:  true,
					},
					client,
				)
			}
		} else if len(args) == 1 {
			controller.PlayListPlayAll(args[0], client)
		} else {
			fmt.Println("playlist name and music required")
		}
	},
}

var setThemeCmd = &cobra.Command{
	Use:   "theme",
	Short: "set the theme [purple|pink|blue]",
	Long: `set the theme [purple|pink|blue]
this command will set the theme of the retro client
the theme is stored in the config file 
`,
	ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"purple", "pink", "blue"}, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(_ *cobra.Command, args []string) {
		if len(args) > 0 {
			theme := strings.TrimSpace(args[0])
			controller.SetTheme(theme, client)
		} else {
			fmt.Println("no theme specified")
		}
	},
}

var logCmd = &cobra.Command{
	Use:   "logs",
	Short: "get logs",
	Long: `get logs
  this command will get the logs from the server
  `,
	Run: func(_ *cobra.Command, _ []string) {
		controller.PrintAllLogs(client)
	},
}

var logErrCmd = &cobra.Command{
	Use:   "error",
	Short: "get error logs",
	Long: `get error logs
  this command will get the error logs from the server
  `,
	Run: func(_ *cobra.Command, _ []string) {
		controller.PrintErrorLogs(client)
	},
}
var logInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get info logs",
	Long: `get info logs
  this command will get the info logs from the server
  `,
	Run: func(_ *cobra.Command, _ []string) {
		controller.PrintInfoLogs(client)
	},
}

var logWarnCmd = &cobra.Command{
	Use:   "warn",
	Short: "get warn logs",
	Long: `get warn logs
  this command will get the warn logs from the server
  `,
	Run: func(_ *cobra.Command, _ []string) {
		controller.PrintWarnLogs(client)
	},
}

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "get cached musics",
	Long: `get cached musics
  this command will get the cached musics from the server
  `,
	Run: func(_ *cobra.Command, _ []string) {
		ctx := views.UIContext{
			Client: client,
			Theme:  views.GetTheme(), // or a factory-loaded theme later
		}
		handler := views.NewCacheHandler()
		view := handler.BuildView(client)
		_ = view.Render(ctx)
	},
}

var cleanCacheCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean the cache",
	Long: `clean the cache
  this command will clean the cache from the server
  `,
	Run: func(_ *cobra.Command, _ []string) {
		controller.CleanCache(client)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update the retro",
	Long: `update the retro
	  this command will update the retro client
			  `,
	Run: func(_ *cobra.Command, _ []string) {
		if err := controller.Update(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("retro updated successfully")
	},
}
