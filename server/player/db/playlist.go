package db

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/XORbit01/retro/logger"
	"github.com/XORbit01/retro/shared"
)

func (d *Db) GetPlaylist(name string) (shared.Playlist, error) {
	var playlist shared.Playlist
	err := d.db.QueryRow(
		`SELECT name FROM playlist WHERE name = ?`,
		name,
	).Scan(
		&playlist.Name,
	)
	return playlist, err
}

func (d *Db) GetPlaylists() ([]shared.Playlist, error) {
	rows, err := d.db.Query(
		`SELECT name,hash FROM playlist`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playlists []shared.Playlist
	for rows.Next() {
		var playlist shared.Playlist
		err := rows.Scan(
			&playlist.Name,
			&playlist.Hash,
		)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}
func hashString(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
func (d *Db) AddPlaylist(plname string) error {
	plname = strings.TrimSpace(plname)
	hash := hashString(plname) // or UUID or md5(plname)

	_, err := d.db.Exec(
		`INSERT OR IGNORE INTO playlist (name, hash) VALUES (?, ?)`,
		plname, hash,
	)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return fmt.Errorf("playlist %s already exists", plname)
	}
	return err
}
func (d *Db) RemovePlaylist(name string) error {
	_, err := d.db.Exec(
		`DELETE FROM playlist WHERE name = ?`,
		name,
	)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(
		`DELETE FROM music_playlist WHERE playlist_name = ?`,
		name,
	)
	return err
}

func (d *Db) AddMusicToPlaylist(musicName, playlistName string) error {
	_, err := d.db.Exec(
		`INSERT INTO music_playlist (music_name, playlist_name) VALUES (?, ?)`,
		musicName,
		playlistName,
	)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return logger.GError("music already in the playlist")
	}
	return err
}

func (d *Db) RemoveMusicFromPlaylist(
	playlistName string,
	musicName string,
) error {
	_, err := d.db.Exec(
		`DELETE FROM music_playlist WHERE music_name= ? AND playlist_name = ?`,
		musicName,
		playlistName,
	)
	return err
}

func (d *Db) GetMusicsFromPlaylist(playlistName string) ([]Music, error) {
	rows, err := d.db.Query(
		`SELECT m.name, m.source, m.key, m.data,m.hash
         FROM music m
        JOIN music_playlist mp ON m.name = mp.music_name
         WHERE mp.playlist_name = ?`,
		playlistName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var musics []Music
	for rows.Next() {
		var music Music
		err := rows.Scan(
			&music.Name,
			&music.Source,
			&music.Key,
			&music.Data,
			&music.Hash,
		)
		if err != nil {
			return nil, err
		}
		musics = append(musics, music)
	}

	return musics, nil
}

func (d *Db) GetPlaylistMusicsMeta(playlistName string) ([]shared.MusicMeta, error) {
	rows, err := d.db.Query(
		`SELECT m.name, m.hash
		 FROM music m
		 JOIN music_playlist mp ON m.name = mp.music_name
		 WHERE mp.playlist_name = ?`,
		playlistName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []shared.MusicMeta
	for rows.Next() {
		var song shared.MusicMeta
		if err := rows.Scan(&song.Name, &song.Hash); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}
func (d *Db) InitPlaylist() error {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS playlist (
      name TEXT PRIMARY KEY,
      hash TEXT UNIQUE NOT NULL
    )`,
	)
	return err
}

// relation ship
func (d *Db) InitMusicPlaylist() error {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS music_playlist (
      music_name TEXT,
      playlist_name TEXT,
      PRIMARY KEY (music_name, playlist_name),
      FOREIGN KEY (music_name) REFERENCES music (name),
      FOREIGN KEY (playlist_name) REFERENCES playlist (name)
    )`,
	)
	return err
}
