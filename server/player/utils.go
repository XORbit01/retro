package player

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"

	"github.com/XORbit01/retro/config"
	"github.com/XORbit01/retro/logger"
	"github.com/XORbit01/retro/server/player/discord"
	"github.com/XORbit01/retro/shared"
)

type customReadCloser struct {
	io.Reader
	io.Seeker
}

func (crc *customReadCloser) Close() error {
	return nil
}

// MusicDecode decodes MP3 data from a byte slice and returns a StreamSeekCloser and Format.
func MusicDecode(data []byte) (beep.StreamSeekCloser, beep.Format, error) {
	reader := bytes.NewReader(data)
	readerCloser := &customReadCloser{Reader: reader, Seeker: reader}
	return mp3.Decode(readerCloser)
}

func copyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func createTmpFile(data []byte) (*os.File, error) {
	f, err := os.CreateTemp("", "retro_")
	if err != nil {
		return nil, err
	}
	if data != nil {
		f.Write(
			data,
		)
	}
	return f, nil
}

func hash(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

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

func adjustDiscordRPC(state shared.PState, music string) {
	if config.GetConfig().DiscordRPC {
		switch state {
		case shared.Stopped:
			if err := discord.Stop(); err != nil {
				logger.LogWarn(
					"error stop discord RPC",
					err,
				)
			}
		case shared.Playing:
			for {
				if err := discord.Start(music); err != nil {
					logger.LogWarn(
						"error start discord RPC trying again in 10 seconds",
						err,
					)
					discord.Stop()
					time.Sleep(10 * time.Second)
				} else {
					logger.LogInfo(
						"Discord RPC started",
						"music", music,
					)
					return
				}
			}
		case shared.Paused:
			if err := discord.Pause(); err != nil {
				logger.LogWarn(
					"error pause discord RPC",
					err,
				)
			} else {
				logger.LogInfo(
					"Discord RPC paused",
					"music", music,
				)
			}
		}
	}
}

// returns true if the player was locked
func (p *Player) unlockIfLocked() bool {
	if p.getPlayerState() == shared.Stopped {
		speaker.Unlock()
		return true
	}
	return false
}

func (p *Player) isSpeakerLocked() bool {
	return p.getPlayerState() == shared.Paused
}

func (p *Player) concernSpeakerLock(callback func()) {
	if p.isSpeakerLocked() {
		speaker.Unlock()
		callback()
		speaker.Lock()
	} else {
		callback()
	}
}
