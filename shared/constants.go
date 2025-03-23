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
