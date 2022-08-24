package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"migu_music_downloader_wails/app"
	"runtime"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	searchUrl := "http://pd.musicapp.migu.cn/MIGUM2.0/v1.0/content/search_all.do?ua=Android_migu&version=5.0.1&pageNo=%d&pageSize=%d&text=%s&searchSwitch="
	downloadUrl := "http://218.205.239.34/MIGUM2.0/v1.0/content/sub/listenSong.do?toneFlag=%s&netType=00&copyrightId=0&&contentId=%s&channel=0"
	configPath := "/usr/local/etc/migu_music_downloader_wails"
	if runtime.GOOS == "windows" {
		configPath = ""
	}

	application, err := app.NewAppQQ(searchUrl, downloadUrl, configPath)
	if err != nil {
		println("Error:", err)
		return
	}

	err = wails.Run(&options.App{
		Title:            "migu music downloader wails",
		Width:            1000,
		Height:           800,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        application.Startup,
		Bind: []interface{}{
			application,
		},
		OnShutdown: application.Stop,
	})

	if err != nil {
		println("Error:", err)
	}
}
