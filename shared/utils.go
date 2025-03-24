package shared

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	Type  int // download, search
	Error string
}

type Status struct {
	CurrMusicIndex    int
	CurrMusicPosition time.Duration
	CurrMusicDuration time.Duration
	PlayerState       PState
	MusicQueue        []string
	Volume            uint8
	Tasks             map[string]Task // key: target, value: task
}

func (s Status) String() string {
	var str string
	str += "CurrMusicIndex: " + fmt.Sprintf("%d", s.CurrMusicIndex) + "\n"
	str += "CurrMusicPosition: " + s.CurrMusicPosition.String() + "\n"
	str += "CurrMusicLength: " + s.CurrMusicDuration.String() + "\n"
	switch s.PlayerState {
	case Playing:
		str += "PlayerState: Playing\n"
	case Paused:
		str += "PlayerState: Paused\n"
	case Stopped:
		str += "PlayerState: Stopped\n"
	}

	str += "Volume: " + fmt.Sprintf("%d", s.Volume) + "\n"

	str += "MusicQueue " + "\n"
	str += "[\n"
	for _, music := range s.MusicQueue {
		str += "\t" + music + "\n"
	}
	str += "]"

	for target, task := range s.Tasks {
		str += fmt.Sprintf("Target: %s, Type: %d, Error: %v\n", target, task.Type, task.Error)
	}

	return str
}

type SearchResult struct {
	Title       string
	Destination string
	Type        DResults
	Duration    time.Duration
}

type RemoveMusicFromPlayListArgs struct {
	PlayListName string
	IndexOrName  IntOrString
}

type PlayListPlayMusicArgs struct {
	PlayListName string
	IndexOrName  IntOrString
}

// function converts 00:00:00 to time.Duration
func StringToDuration(s string) (time.Duration, error) {
	sp := strings.Split(s, ":")
	if len(sp) < 2 {
		return 0, fmt.Errorf("invalid duration: %s", s)
	}
	l := len(sp)
	sec := "0"
	min := "0"
	hour := "0"
	if l > 0 {
		sec = sp[l-1]
	}

	if l > 1 {
		min = sp[l-2]
	}
	if l > 2 {
		hour = sp[l-3]
	}
	return time.ParseDuration(hour + "h" + min + "m" + sec + "s")
}

func DurationToString(d time.Duration) string {
	// to format 00:00:00
	return fmt.Sprintf("%02d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
}

type IntOrString struct {
	IntVal int
	StrVal string
	IsInt  bool
}

type HashNamed struct {
	Name string
	Hash string
}

type DetectQuery struct {
	Query   string
	Knowing DResults
}
type AddToPlayListQuery struct {
	PlayListName string
	Query        string
	Knowing      DResults
}

type Playlist = HashNamed
type MusicMeta = HashNamed
type CacheItem = HashNamed
