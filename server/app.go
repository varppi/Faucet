package main

import (
	"FaucetServer/internal/config"
	"FaucetServer/internal/server"
	"context"
	"fmt"
	"strconv"

	"github.com/tawesoft/golib/v2/dialog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	err := config.Init()
	if err != nil {
		dialog.Error(fmt.Sprintf("Startup failed!\n\nError: %s", err.Error()))
		a.exit()
	}

	err = server.Listen()
	if err != nil {
		dialog.Error(fmt.Sprintf("Server listen failed!\n\nError: %s", err.Error()))
		a.exit()
	}

	a.ctx = ctx
}

func (a *App) exit() {
	runtime.Quit(a.ctx)
}

var cache = make(map[string]any)

func (a *App) GetListener() string {
	return config.Config["listen"]
}

func (a *App) GetConnected() bool {
	return server.Connected
}

func (a *App) GetFiles() map[string]any {
	cache = server.ListFiles()
	return cache
}

func (a *App) PressedItem(name string) any {
	if cache[name] == "d" || name == "../" {
		server.ChDir(name)
		return nil
	} else {
		return server.Preview(name)
	}
}

func (a *App) DownloadItem(name string) {
	if cache[name] == nil {
		return
	}
	server.Downloading = 1
	name = server.Download(name, cache[name].(string), config.Config["lootdir"])
	if name != "" {
		go dialog.Info(fmt.Sprintf("saved to %s", name))
	}
}

func (a *App) IsDownloading() string {
	if server.Downloading == -1 {
		return ""
	}
	return strconv.Itoa(server.Downloading)
}

func (a *App) KillAgent() {
	server.KillAgent()
}
