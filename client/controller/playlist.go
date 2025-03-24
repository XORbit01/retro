package controller

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/XORbit01/retro/shared"
)

// TODO:
func GetPlayListsMeta(client *rpc.Client) []shared.Playlist {
	var reply []shared.Playlist
	err := client.Call("Player.RPCPlayListsMeta", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func CreatePlayList(name string, client *rpc.Client) {
	args := name
	var reply int
	err := client.Call("Player.RPCCreatePlayList", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RemovePlayList(name string, client *rpc.Client) {
	args := name
	var reply int
	err := client.Call("Player.RPCRemovePlayList", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DetectAndAddToPlayList(
	query shared.AddToPlayListQuery,
	client *rpc.Client,
) ([]shared.SearchResult, error) {
	var reply []shared.SearchResult
	err := client.Call("Player.RPCDetectAndAddToPlayList", query, &reply)
	return reply, err
}

func GetPlayListMusicsMeta(name string, client *rpc.Client) []shared.MusicMeta {
	args := name
	var reply []shared.MusicMeta
	err := client.Call("Player.RPCPlayListMusicsMeta", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func RemoveMusicFromPlayList(name string, indexOrName shared.IntOrString, client *rpc.Client) {
	args := shared.RemoveMusicFromPlayListArgs{
		PlayListName: name,
		IndexOrName:  indexOrName,
	}
	var reply int
	err := client.Call("Player.RPCRemoveMusicFromPlayList", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PlayListPlayMusic(lname string, indexOrName shared.IntOrString, client *rpc.Client) {
	args := shared.PlayListPlayMusicArgs{
		PlayListName: lname,
		IndexOrName:  indexOrName,
	}
	var reply int
	err := client.Call("Player.RPCPlayListPlayMusicMeta", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PlayListPlayAll(name string, client *rpc.Client) {
	args := name
	var reply int
	err := client.Call("Player.RPCPlayListPlayAll", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
