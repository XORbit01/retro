package player

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/XORbit01/retro/config"
	"github.com/XORbit01/retro/logger"
	"github.com/XORbit01/retro/server/player/db"
	"github.com/XORbit01/retro/shared"
)

func (p *Player) CheckWhatIsThis(unknown string) shared.DResults {
	// check if its a dir or file
	// TODO: check if its local path
	if fi, err := os.Stat(unknown); err == nil {
		if fi.IsDir() {
			// check if there is music files in the dir
			files, err := os.Open(unknown)
			if err != nil {
				return shared.DUnknown
			}

			entries, err := files.Readdir(0)
			if err != nil {
				return shared.DUnknown
			}

			for _, entry := range entries {
				if !entry.IsDir() {
					data, err := os.ReadFile(
						filepath.Join(
							unknown,
							entry.Name(),
						),
					)
					if err != nil {
						logger.LogWarn(
							"error reading file",
							err,
						)
					}
					logger.LogInfo(
						"Checking if",
						filepath.Join(unknown, entry.Name()),
						"is mp3",
					)
					isMp3, err := p.Director.Converter.IsMp3(
						data,
					)
					if err != nil {
						logger.LogWarn(
							"Failed to check if",
							filepath.Join(unknown, entry.Name()),
							"is mp3",
							err,
						)
					}
					if isMp3 {
						return shared.DDir
					}
				}
			}
			return shared.DUnknown
		} else {
			data, err := os.ReadFile(
				unknown,
			)
			if err != nil {
				return shared.DUnknown
			}
			isMp3, err := p.Director.Converter.IsMp3(
				data,
			)
			if err != nil {
				logger.LogWarn(
					"Failed to check if",
					unknown,
					"is mp3",
					err,
				)
			}
			if isMp3 {
				return shared.DFile
			}
		}
	}

	// check if its cache hash prefix
	_, err := p.Director.Db.GetMusicByHashPrefix(unknown)
	if err == nil {
		return shared.DCache
	}
	// check if its play list name
	_, err = p.Director.Db.GetPlaylist(
		unknown,
	)
	if err == nil {
		return shared.DPlaylist
	}
	i, err := strconv.Atoi(unknown)
	if err != nil {
		ok := p.Queue.GetMusicByName(
			unknown,
		)
		if ok != nil {
			return shared.DQueue
		}
	} else {
		if i < p.Queue.Size() && i >= 0 {
			return shared.DQueue
		}
	}

	engines := p.Director.GetEngines()
	for _, engine := range engines {
		ok, _ := engine.Exists(unknown)
		if ok {
			return engine.Name()
		}
	}
	return shared.DUnknown
}

func (p *Player) searchWorker(
	engine shared.DResults,
	unknown string,
	musicChan chan shared.SearchResult,
	wg *sync.WaitGroup,
) {
	defer func() {
		logger.LogInfo(
			"Search worker done for",
			engine,
			unknown,
		)
		wg.Done()
	}()

	searchRes, err := p.Director.Search(
		engine,
		unknown,
	)
	if err != nil {
		logger.LogWarn("Failed to search for", unknown, ":", err)
	}

	for _, music := range searchRes {
		musicChan <- music
	}
}

func (p *Player) GetAvailableMusicOptions(unknown string) []shared.SearchResult {
	// add task : this task displayed in the status: if the task is done, it will be removed
	p.addTask(
		unknown,
		shared.Searching,
	)
	wg := &sync.WaitGroup{}
	musicChan := make(
		chan shared.SearchResult,
	)
	var musics []shared.SearchResult
	ctx, cancel := context.WithTimeout(context.Background(), config.GetConfig().SearchTimeout)
	defer cancel()
	for name := range p.Director.GetEngines() {
		wg.Add(1)
		go p.searchWorker(
			name,
			unknown,
			musicChan,
			wg,
		)
	}

	wg.Add(1)
	go func() {
		// Get cached music
		defer wg.Done()
		ms, err := p.Director.Db.FilterMusic(
			unknown,
		)
		if err != nil {
			return
		}

		for _, m := range ms {
			music, err := NewMusic(
				m.Name,
				m.Data,
			)
			if err != nil {
				logger.LogWarn(
					"skipping music",
					err,
				)
				continue
			}
			var dur time.Duration
			p.concernSpeakerLock(
				func() {
					dur = music.DurationD()
				},
			)
			musicChan <- shared.SearchResult{
				Title:       m.Name,
				Destination: m.Key,
				Duration:    dur,
				Type:        "cache",
			}
		}
	}()

	go func() {
		wg.Wait()
		close(musicChan)
	}()

	for {
		select {
		case music, ok := <-musicChan:
			if !ok {
				return musics
			}
			musics = append(musics, music)
			p.removeTask(
				unknown,
			)
		case <-ctx.Done():
			p.errorTask(
				unknown,
				fmt.Errorf("search timed out"),
			)
			logger.LogWarn(
				"Search timed out",
				unknown,
			)
			return musics
		}
	}
}

func (p *Player) DetectAndAddToPlayList(
	query shared.AddToPlayListQuery,
) ([]shared.SearchResult, error) {
	// if client know what is this, lets save time without detecting
	whatIsThis := query.Knowing
	if whatIsThis == shared.DUnknown {
		whatIsThis = p.CheckWhatIsThis(query.Query)
	}
	pl, err := p.Director.Db.GetPlaylist(
		query.PlayListName,
	)
	if err != nil {
		return nil, logger.LogError(
			logger.GError(
				"Playlist does not exist",
			),
		)
	}
	switch whatIsThis {
	case shared.DDir:
		logger.LogInfo(
			"detected dir for",
			query.Query,
		)
		return nil, p.AddMusicsFromDir(
			query.Query,
			func(m db.Music) error {
				return p.Director.Db.AddMusicToPlaylist(
					m.Name,
					pl.Name,
				)
			},
		)
	case shared.DFile:
		logger.LogInfo(
			"Detected file",
			query.Query,
		)
		err := p.AddMusicFromFile(
			query.Query,
			func(m db.Music) error {
				return p.Director.Db.AddMusicToPlaylist(
					m.Name,
					pl.Name,
				)
			},
		)
		if err != nil {
			return nil, logger.LogError(
				logger.GError(
					"Failed to add to playlist",
					err,
				),
			)
		}
	case shared.DQueue:
		logger.LogInfo(
			"Detected queue",
			query.Query,
		)
		index, err := strconv.Atoi(query.Query)
		var m *Music
		if err == nil {
			m = p.Queue.GetMusicByIndex(
				index,
			)
		} else {
			m = p.Queue.GetMusicByName(
				query.Query,
			)
		}
		return nil, p.Director.Db.AddMusicToPlaylist(
			m.Name,
			pl.Name,
		)
	case shared.DCache:
		logger.LogInfo(
			"Detected cache",
			query.Query,
		)
		return nil, p.AddMusicFromHash(
			query.Query,
			func(m db.Music) error {
				return p.Director.Db.AddMusicToPlaylist(
					m.Name,
					query.Query,
				)
			},
		)
	case shared.DUnknown:
		logger.LogInfo(
			"Detected unknown",
			query.Query,
		)
		return p.GetAvailableMusicOptions(query.Query), nil
	default:
		logger.LogInfo(
			"Detected Engine",
			whatIsThis,
		)
		go p.AddMusicFromOnline(
			query.Query,
			whatIsThis,
			func(m db.Music) error {
				return p.Director.Db.AddMusicToPlaylist(
					m.Name,
					pl.Name,
				)
			},
		)
	}
	return []shared.SearchResult{}, nil
}

// DetectAndPlay if result is empty, it means it detects and plays the music if succeed other wise it returns the search results
func (p *Player) DetectAndPlay(query shared.DetectQuery) ([]shared.SearchResult, error) {
	logger.LogInfo("Checking what is this", query)
	// if client know what is this, lets save time without detecting
	whatIsThis := query.Knowing
	if whatIsThis == shared.DUnknown {
		whatIsThis = p.CheckWhatIsThis(query.Query)
	}

	switch whatIsThis {
	case shared.DDir:
		logger.LogInfo("Detected dir")
		return nil, p.AddMusicsFromDir(
			query.Query,
			func(m db.Music) error {
				pmusic, err := NewMusic(
					m.Name,
					m.Data,
				)
				if err != nil {
					return err
				}
				p.Queue.Enqueue(*pmusic)
				err = p.Play()
				return err
			},
		)
	case shared.DFile:
		logger.LogInfo("Detected file")
		return nil, p.AddMusicFromFile(
			query.Query,
			func(m db.Music) error {
				pmusic, err := NewMusic(
					m.Name,
					m.Data,
				)
				if err != nil {
					return err
				}
				p.Queue.Enqueue(*pmusic)
				p.Play()
				return nil
			},
		)
	case shared.DQueue:
		logger.LogInfo(
			"Detected queue",
			query,
		)
		index, err := strconv.Atoi(query.Query)
		var m *Music
		if err == nil {
			m = p.Queue.GetMusicByIndex(index)
		} else {
			m = p.Queue.GetMusicByName(query.Query)
		}
		p.Queue.SetCurrrMusic(m)
		return nil, p.Play()
	case shared.DPlaylist:
		return nil, p.PlayListPlayAll(
			query.Query,
		)
	case shared.DUnknown:
		logger.LogInfo("Detected unknown, searching for", query)
		return p.GetAvailableMusicOptions(query.Query), nil
	case shared.DCache:
		logger.LogInfo("Detected cache, searching for", query)
		return nil, p.AddMusicFromHash(
			query.Query,
			func(m db.Music) error {
				pmusic, err := NewMusic(
					m.Name,
					m.Data,
				)
				if err != nil {
					return err
				}
				logger.LogInfo(
					"Enqueue the music", m.Name,
				)
				p.Queue.Enqueue(*pmusic)
				return p.Play()
			},
		)
	default:
		logger.LogInfo("Detected Engine", whatIsThis)
		go p.AddMusicFromOnline(
			query.Query,
			whatIsThis,
			func(m db.Music) error {
				pmusic, err := NewMusic(
					m.Name,
					m.Data,
				)
				if err != nil {
					return err
				}
				p.Queue.Enqueue(*pmusic)
				p.Play()
				return nil
			},
		)
	}
	return []shared.SearchResult{}, nil
}
