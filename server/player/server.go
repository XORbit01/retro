package player

import (
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/XORbit01/retro/logger"
	"github.com/XORbit01/retro/shared"
)

func (p *Player) RPCPlay(_ int, reply *int) error {
	logger.LogInfo("RPCPlay called")
	err := p.Play()
	logger.LogInfo("RPCPlay done")
	return err
}

func (p *Player) RPCNext(_ int, reply *int) error {
	logger.LogInfo("RPCNext called")
	err := p.Next()
	*reply = 1
	logger.LogInfo("RPCNext done")
	return err
}

func (p *Player) RPCPrev(_ int, reply *int) error {
	logger.LogInfo("RPCPrev called")
	err := p.Prev()
	*reply = 1
	logger.LogInfo("RPCPrev done")
	return err
}

func (p *Player) RPCPause(_ int, reply *int) error {
	logger.LogInfo("RPCPause called")
	err := p.Pause()
	*reply = 1
	logger.LogInfo("RPCPause done")
	return err
}

func (p *Player) RPCStop(_ int, reply *int) error {
	logger.LogInfo("RPCStop called")
	err := p.Stop()
	*reply = 1
	logger.LogInfo("RPCStop done")
	return err
}

func (p *Player) RPCSeek(d time.Duration, _ *int) error {
	logger.LogInfo("RPCSeek called with duration in seconds :", d*time.Second)
	err := p.Seek(d * time.Second)
	logger.LogInfo("RPCSeek done")
	return err
}

func (p *Player) RPCVolume(vp uint8 /*volume percentage*/, reply *int) error {
	logger.LogInfo("RPCVolume called with volume percentage :", vp)
	err := p.Volume(vp)
	*reply = 0
	logger.LogInfo("RPCVolume done")
	return err
}

func (p *Player) RPCResume(_ int, reply *int) error {
	logger.LogInfo("RPCResume called")
	err := p.Resume()
	*reply = 1
	logger.LogInfo("RPCResume done")
	return err
}

func (p *Player) RPCRemoveMusic(music shared.IntOrString, reply *int) error {
	logger.LogInfo("RPCRemoveMusic called with :", music)
	err := p.Remove(music)
	*reply = 1
	logger.LogInfo("RPCRemoveMusic done")
	return err
}

func (p *Player) RPCGetPlayerStatus(_ int, reply *shared.Status) error {
	logger.LogInfo("RPCGetPlayerStatus called")
	*reply = p.GetPlayerStatus()
	logger.LogInfo("RPCGetPlayerStatus done with reply :", *reply)
	return nil
}

func (p *Player) RPCDetectAndPlay(query shared.DetectQuery, reply *[]shared.SearchResult) error {
	logger.LogInfo("RPCDetectAndPlay called with query :", query)
	var err error
	*reply, err = p.DetectAndPlay(query)
	logger.LogInfo("RPCDetectAndPlay done with reply :", *reply)
	return err
}

func (p *Player) RPCPlayListsMeta(_ int, reply *[]shared.Playlist) error {
	logger.LogInfo("RPCPlayListMusics called")
	var err error
	*reply, err = p.PlayListsMeta()
	logger.LogInfo("RPCPlayListMusics done with reply :", *reply)
	return err
}

func (p *Player) RPCCreatePlayList(name string, reply *int) error {
	logger.LogInfo("RPCCreatePlaylist called with name :", name)
	err := p.CreatePlayList(name)
	*reply = 1
	logger.LogInfo("RPCCreatePlaylist done")
	return err
}

func (p *Player) RPCRemovePlayList(name string, reply *int) error {
	logger.LogInfo("RPCRemovePlaylist called with name :", name)
	err := p.RemovePlayList(name)
	*reply = 1
	logger.LogInfo("RPCRemovePlaylist done")
	return err
}

func (p *Player) RPCDetectAndAddToPlayList(
	args shared.AddToPlayListQuery,
	reply *[]shared.SearchResult,
) error {
	logger.LogInfo(
		"RPCDetectAndAddToPlayList called with query :",
		args.Query,
		" and playlist name :",
		args.PlayListName,
	)
	var err error
	*reply, err = p.DetectAndAddToPlayList(args)
	logger.LogInfo("RPCDetectAndAddToPlayList done")
	return err
}

func (p *Player) RPCPlayListMusicsMeta(
	plname string,
	reply *[]shared.MusicMeta,
) error {
	logger.LogInfo("RPCPlayListMusics called with name :", plname)
	var err error
	*reply, err = p.GetPlayListMusicsMeta(plname)
	logger.LogInfo("RPCPlayListMusics done with reply :", *reply)
	return err
}

func (p *Player) RPCRemoveMusicFromPlayList(
	args shared.RemoveMusicFromPlayListArgs,
	reply *int,
) error {
	logger.LogInfo(
		"RPCRemoveMusicFromPlayList called with name :",
		args.PlayListName,
		"target",
		args.IndexOrName,
	)
	err := p.RemoveMusicFromPlayList(
		args.PlayListName,
		args.IndexOrName,
	)
	logger.LogInfo(
		"RPCRemoveMusicFromPlayList done",
	)
	return err
}

func (p *Player) RPCPlayListPlayMusic(args shared.PlayListPlayMusicArgs, reply *int) error {
	logger.LogInfo(
		"RPCPlayListPlayMusic called with name :",
		args.PlayListName,
		"target",
		args.IndexOrName,
	)
	err := p.PlayListPlayMusic(args.PlayListName, args.IndexOrName)
	*reply = 1
	logger.LogInfo("RPCPlayListPlayMusic done")
	return err
}

func (p *Player) RPCPlayListPlayAll(name string, reply *int) error {
	logger.LogInfo("RPCPlayListPlayAll called with name :", name)
	p.PlayListPlayAll(name)
	*reply = 1
	logger.LogInfo("RPCPlayListPlayAll done")
	return nil
}

func (p *Player) RPCGetTheme(_ int, reply *string) error {
	logger.LogInfo("RPCGetTheme called")
	*reply = p.GetTheme()
	logger.LogInfo("RPCGetTheme done with reply :", *reply)
	return nil
}

func (p *Player) RPCSetTheme(theme string, reply *int) error {
	logger.LogInfo("RPCSetTheme called with theme :", theme)
	p.SetTheme(theme)
	*reply = 1
	logger.LogInfo("RPCSetTheme done")
	return nil
}

func (p *Player) RPCGetLogs(_ int, reply *[]string) error {
	logger.LogInfo("GetLogs called")
	var err error
	*reply, err = logger.GetLogs()
	logger.LogInfo("GetLogs done")
	return err
}

func (p *Player) RPCCleanCache(_ int, reply *int) error {
	logger.LogInfo("RPCCleanCache called")
	p.CleanCache()
	*reply = 1
	logger.LogInfo("RPCCleanCache done")
	return nil
}

func (p *Player) RPCGetCachedMusics(_ int, reply *[]shared.HashNamed) error {
	logger.LogInfo("RPCGetCachedMusics called")
	var err error
	*reply, err = p.GetCachedMusics()
	logger.LogInfo("RPCGetCachedMusics done with reply :", *reply)
	return err
}

func StartIPCServer(port string) {

	// check update

	logger.LogInfo("Creating Player instance")
	player := NewPlayer()
	err := rpc.Register(player)
	if err != nil {
		logger.LogError(
			logger.GError(
				"Failed to register Player instance to RPC",
				err,
			),
		)
		return
	}
	logger.LogInfo("Player instance created and registered to RPC")
	lis, err := net.Listen("tcp", ":"+port)

	logger.LogInfo("Starting IPC server on ", lis.Addr().String())
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			logger.LogError(
				logger.GError(
					"Failed to accept connection",
					err,
				),
			)
		}
		go rpc.ServeConn(conn)
	}
}
