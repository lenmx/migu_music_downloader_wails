package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"migu_music_downloader_wails/app/consts"
	"migu_music_downloader_wails/app/model"
	"migu_music_downloader_wails/app/util"
	"strconv"
)

type App struct {
	ctx         context.Context
	SearchUrl   string
	DownloadUrl string

	Downloader *util.Downloader
}

func NewApp(searchUrl, downloadUrl string) *App {
	d := util.NewDownloader(10)
	go d.Run()

	app := &App{
		SearchUrl:   searchUrl,
		DownloadUrl: downloadUrl,
		Downloader:  d,
	}

	go app.watchDownloadResult()
	return app
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OnSearch(keyword string, pageIndex, pageSize int) model.BaseResponse {
	res, err := a.search(keyword, pageIndex, pageSize)
	if err != nil {
		return a.Error(err.Error())
	}

	if res.Code != consts.SearchSuccess {
		return a.Error(res.Info)
	}

	total, _ := strconv.Atoi(res.SongResultData.TotalCount)
	resp := &model.PageRes{
		Total: total,
		Items: res.SongResultData.Result,
	}

	return a.Ok(resp)
}

func (a *App) OnDownload(sourceType string, downloadItemsJson string) model.BaseResponse {
	var items []model.DownloadItem
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
		a.log(fmt.Sprintf("[%s]添加成功 ", item.Name))
		a.download(sourceType, path, item)
	}

	return a.Ok(nil)
}

func (a *App) search(keyword string, pageIndex, pageSize int) (*model.SearchRes, error) {
	// http://pd.musicapp.migu.cn/MIGUM2.0/v1.0/content/search_all.do?ua=Android_migu&version=5.0.1&pageNo=1&pageSize=10&text=周杰伦&searchSwitch=
	url := fmt.Sprintf(a.SearchUrl, pageIndex, pageSize, keyword)

	res, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}

	var resp model.SearchRes
	err = json.Unmarshal(res.Body(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (a *App) download(sourceType, path string, item model.DownloadItem) {
	_sourceType := consts.SourceType_HQ
	if sourceType == "SQ" {
		_sourceType = consts.SourceType_SQ
	}

	if path[len(path)-1] != '/' {
		path += "/"
	}

	path += item.Name + consts.SourceType2FileExt[_sourceType]
	url := fmt.Sprintf(a.DownloadUrl, string(_sourceType), item.ContentId)

	a.Downloader.Push(model.DownloadQueueItem{
		DownloadItem: item,
		Path:         path,
		Url:          url,
	})
}

func (a *App) watchDownloadResult() {
	for {
		select {
		case res, ok := <-a.Downloader.ResultCh:
			if !ok {
				return
			}

			item := res.Data.(model.DownloadQueueItem)
			a.log(fmt.Sprintf("[%s]%s", item.Name, res.Message))
			a.pushDownloadResult(res)
		default:
		}
	}
}

func (a *App) Ok(data interface{}) model.BaseResponse {
	return model.BaseResponse{
		Code:    0,
		Message: "",
		Data:    data,
	}
}

func (a *App) Error(message string) model.BaseResponse {
	return model.BaseResponse{
		Code:    -1,
		Message: message,
		Data:    nil,
	}
}

func (a *App) log(message string) {
	util.Log(a.ctx, message)
}

func (a *App) pushDownloadResult(response model.BaseResponse) {
	util.PushDownloadResult(a.ctx, response)
}
