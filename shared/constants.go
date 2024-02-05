package shared

import "time"

const (
	NotStarted = iota
	Running
	Finished
)

const (
	Downloading = iota
	Searching
)

const (
	Playing = iota
	Paused
	Stopped
)

var Separator = "_#__#_"

//TODO: make this configurable + create it if not exists
var GoPlayPath = "./goplay_storage/"
var CachePath = GoPlayPath + "cache/"
var PlaylistPath = GoPlayPath + "playlists/"

var Pathytldpl = "yt-dlp"
var Pathffmpeg = "ffmpeg"
var Pathffprobe = "ffprobe"

var SearchTimeOut = 60 * time.Second
