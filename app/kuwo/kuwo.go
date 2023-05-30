package kuwo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"migu_music_downloader_wails/app"
	"migu_music_downloader_wails/app/model"
	"migu_music_downloader_wails/app/util"
	"path"
	"strings"
)

type AppKuwo struct {
	ctx    context.Context
	cancel context.CancelFunc

	core       *AppKuwoCore
	downloader *util.Downloader
}

func NewAppKuwo() *AppKuwo {
	app := &AppKuwo{
		core: NewAppKuwoCore(),
	}

	// downloader
	d := util.NewDownloader(app.onDownloadResult, 10)
	app.downloader = d
	go app.downloader.Run()

	return app
}

func (a *AppKuwo) OnSearch(keyword string, pageIndex, pageSize int) model.BaseResponse {
	res, err := a.core.Search(keyword, pageIndex, pageSize)
	if err != nil {
		return app.GApp.GenError(err.Error())
	}

	total := res.GetTotal()
	if total <= 0 || res.Abslist == nil || len(res.Abslist) <= 0 {
		return app.GApp.GenError(app.GApp.TR("SearchFail"))
	}

	items := []model.SearchItemRes{}
	for _, item := range res.Abslist {
		singers := []string{item.ARTIST}
		albums := []string{item.ALBUM}

		items = append(items, model.SearchItemRes{
			MusicId: strings.ReplaceAll(item.MUSICRID, "MUSIC_", ""),
			Name:    item.NAME,
			Title:   item.SONGNAME,
			Singers: singers,
			Albums:  albums,
			LrcUrl:  "",
			Cover:   "",
		})
	}

	resp := &model.PageRes{
		Total: total,
		Items: items,
	}

	return app.GApp.GenOk(resp)
}

func (a *AppKuwo) OnDownload(sourceType string, downloadItemJson string) model.BaseResponse {
	defer func() {
		if e := recover(); e != nil {
			app.GApp.Log("OnDownload panic: " + util.PanicTrace(e))
		}
	}()

	util.Log(a.ctx, fmt.Sprintf("%v, %v", sourceType, downloadItemJson))
	if sourceType == "SQ" {
		return app.GApp.GenError(app.GApp.TR("NoSQFile"))
	}

	var items []model.KuwoDownloadItem
	json.Unmarshal([]byte(downloadItemJson), &items)

	savePath := ""
	downloadLrc := false
	downloadCover := false
	setting, _ := app.GApp.GetSetting()
	if setting != nil && len(setting.SavePath) > 0 {
		savePath = setting.SavePath
		downloadLrc = setting.DownloadLrc
		downloadCover = setting.DownloadCover
	} else {
		savePath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			DefaultDirectory:           "",
			DefaultFilename:            "",
			Title:                      app.GApp.TR("ChooseFileSavePath"),
			Filters:                    nil,
			ShowHiddenFiles:            false,
			CanCreateDirectories:       false,
			ResolvesAliases:            false,
			TreatPackagesAsDirectories: false,
		})
		if err != nil {
			return app.GApp.GenError(app.GApp.TR("ChooseFileSavePathFail"))
		}
		if len(savePath) <= 0 {
			return app.GApp.GenError(app.GApp.TR("CancelDownload"))
		}
	}

	if savePath[len(savePath)-1] != '/' {
		savePath += "/"
	}

	for _, item := range items {
		downloadUrl := a.core.GetDownloadUrl1(item.MusicId)
		if downloadUrl == "" {
			continue
		}

		filename := path.Join(savePath, util.FixWindowsFileName2Normal(item.MusicName)+path.Ext(downloadUrl))
		queueItem := model.DownloadQueueItem{
			MusicId:    item.MusicId,
			Name:       item.MusicName,
			Filename:   filename,
			Url:        downloadUrl,
			Mp3Process: a.core.ProcessMp3,
		}
		if downloadLrc {
			queueItem.LrcUrl = a.core.GetLrcUrl(item.MusicId)
			queueItem.LrcProcess = a.core.ProcessLrc
		}
		if downloadCover {
			queueItem.PicUrl = a.core.GetPicUrl(item.MusicId)
			queueItem.PicProcess = a.core.ProcessPic
		}

		a.downloader.Push(a.ctx, queueItem)
		app.GApp.Log(fmt.Sprintf("[%s]%s", item.MusicName, app.GApp.TR("AddToDownloadCenterSuccess")))
	}

	return app.GApp.GenOk(nil)
}

func (a *AppKuwo) onDownloadResult(res model.BaseResponse) {
	//item := res.Data.(model.DownloadQueueItem)
	app.GApp.Log(fmt.Sprintf("%+v", res))
	a.pushDownloadResult(res)
}

func (a *AppKuwo) pushDownloadResult(response model.BaseResponse) {
	util.PushDownloadResult(a.ctx, response)
}

func (a *AppKuwo) Startup(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}

func (a *AppKuwo) Stop(ctx context.Context) {
	a.cancel()
	a.downloader.Stop()
}
