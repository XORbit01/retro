package controller

import (
	"fmt"
	"github.com/XORbit01/retro/config"
	"net/rpc"
	"os"

	"github.com/XORbit01/retro/shared"
)

func Next(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCNext", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Prev(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCPrev", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Pause(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCPause", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Resume(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCResume", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Stop(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCStop", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Seek(d int, client *rpc.Client) {
	args := d
	var reply int
	err := client.Call("Player.RPCSeek", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Volume(vp uint8, client *rpc.Client) {
	if vp > 100 {
		// health warning
		fmt.Print(" ⚠️ Volume greater than 100% may damage your ears, skip this warning? (y/n)")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil || response != "y" {
			return
		}
	}
	args := vp
	var reply int
	err := client.Call("Player.RPCVolume", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Remove(index interface{}, client *rpc.Client) {
	args := index
	var reply int
	err := client.Call("Player.RPCRemoveMusic", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetPlayerStatus(client *rpc.Client) shared.Status {
	var reply shared.Status
	err := client.Call("Player.RPCGetPlayerStatus", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func DetectAndPlay(query shared.DetectQuery, client *rpc.Client) ([]shared.SearchResult, error) {
	var reply []shared.SearchResult
	err := client.Call("Player.RPCDetectAndPlay", query, &reply)
	return reply, err
}

func GetTheme(client *rpc.Client) string {
	var reply string
	err := client.Call("Player.RPCGetTheme", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func SetTheme(theme string, client *rpc.Client) {
	args := theme
	var reply int
	err := client.Call("Player.RPCSetTheme", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetCachedMusics(client *rpc.Client) []shared.HashNamed {
	var reply []shared.HashNamed
	err := client.Call("Player.RPCGetCachedMusics", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func CleanCache(client *rpc.Client) {
	var reply int
	err := client.Call("Player.RPCCleanCache", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var client *rpc.Client

func GetClient() (*rpc.Client, error) {
	cfg := config.GetConfig()
	if client == nil {
		var err error
		client, err = rpc.Dial("tcp", "localhost:"+cfg.ServerPort)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}
