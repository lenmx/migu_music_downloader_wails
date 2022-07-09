package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strconv"
)

type App struct {
	ctx         context.Context
	SearchUrl   string
	DownloadUrl string

	Downloader *Downloader
	//queue chan DownloadQueueItem
}

func NewApp(searchUrl, downloadUrl string) *App {
	d := NewDownloader(10)
	go d.Run()

	return &App{
		SearchUrl:   searchUrl,
		DownloadUrl: downloadUrl,
		Downloader:  d,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OnSearch(keyword string, pageIndex, pageSize int) BaseResponse {
	res, err := a.search(keyword, pageIndex, pageSize)
	if err != nil {
		return a.Error(err.Error())
	}

	if res.Code != searchSuccess {
		return a.Error(res.Info)
	}

	total, _ := strconv.Atoi(res.SongResultData.TotalCount)
	resp := &PageRes{
		Total: total,
		Items: res.SongResultData.Result,
	}

	return a.Ok(resp)
}

func (a *App) OnDownload(sourceType string, downloadItemsJson string) BaseResponse {
	var items []DownloadItem
	err := json.Unmarshal([]byte(downloadItemsJson), &items)
	if err != nil {
		return a.Error(err.Error())
	}

	path, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:           "/",
		DefaultFilename:            "",
		Title:                      "选择保存路径",
		Filters:                    nil,
		ShowHiddenFiles:            false,
		CanCreateDirectories:       false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	})
	if len(path) <= 0 {
		return a.Error("取消下载")
	}

	for _, item := range items {
		a.Log("添加成功 " + item.Name)
		a.download(sourceType, path, item)
	}

	return a.Ok(nil)
}

func (a *App) search(keyword string, pageIndex, pageSize int) (*SearchRes, error) {
	// http://pd.musicapp.migu.cn/MIGUM2.0/v1.0/content/search_all.do?ua=Android_migu&version=5.0.1&pageNo=1&pageSize=10&text=周杰伦&searchSwitch=
	url := fmt.Sprintf(a.SearchUrl, pageIndex, pageSize, keyword)

	res, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}

	var resp SearchRes
	err = json.Unmarshal(res.Body(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (a *App) download(sourceType, path string, item DownloadItem) {
	_sourceType := SourceType_HQ
	if sourceType == "SQ" {
		_sourceType = SourceType_SQ
	}

	if path[len(path)-1] != '/' {
		path += "/"
	}

	path += item.Name + SourceType2FileExt[_sourceType]
	url := fmt.Sprintf(a.DownloadUrl, string(sourceType), item.ContentId)

	a.Downloader.Push(DownloadQueueItem{
		DownloadItem: item,
		Path:         path,
		Url:          url,
		ctx:          a.ctx,
	})
}

func (a *App) Ok(data interface{}) BaseResponse {
	return BaseResponse{
		Code:    0,
		Message: "",
		Data:    data,
	}
}

func (a *App) Error(message string) BaseResponse {
	return BaseResponse{
		Code:    -1,
		Message: message,
		Data:    nil,
	}
}

func (a *App) Log(message string) {
	runtime.EventsEmit(a.ctx, "log", message)
}
