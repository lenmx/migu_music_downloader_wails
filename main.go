package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"migu_music_downloader_wails/app"
	"migu_music_downloader_wails/app/kuwo"
	"runtime"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	configPath := "/usr/local/etc/migu_music_downloader_wails"
	if runtime.GOOS == "windows" {
		configPath = ""
	}

	application, err := app.NewApp(configPath)
	if err != nil {
		println("Error:", err)
		return
	}

	kuwo := kuwo.NewAppKuwo()

	err = wails.Run(&options.App{
		Title:            "migu music downloader wails",
		Width:            1000,
		Height:           800,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			application.Startup(ctx)
			kuwo.Startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			application.Stop(ctx)
			kuwo.Stop(ctx)
		},
		Bind: []interface{}{
			application,
			kuwo,
		},
	})

	if err != nil {
		println("Error:", err)
	}
}
