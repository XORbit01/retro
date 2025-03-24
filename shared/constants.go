package shared

const (
	Downloading = iota
	Searching
)

type PState uint

const (
	Playing PState = iota
	Paused
	Stopped
)

const (
	HashPrefixLength = 5
)

type DResults string

const (
	DUnknown  DResults = "unknown"
	DDir      DResults = "dir"
	DFile     DResults = "file"
	DQueue    DResults = "queue"
	DPlaylist DResults = "playlist"
	DYoutube  DResults = "youtube"
	DCache    DResults = "cache"
)
