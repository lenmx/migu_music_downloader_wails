package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"migu_music_downloader_wails/app/consts"
	"migu_music_downloader_wails/app/model"
	"migu_music_downloader_wails/app/util"
	"os"
	"strconv"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	searchUrl   string
	downloadUrl string
	configPath  string

	downloader *util.Downloader
}

func NewApp(searchUrl, downloadUrl, configPath string) *App {
	app := &App{
		searchUrl:   searchUrl,
		downloadUrl: downloadUrl,
		configPath:  configPath,
	}

	d := util.NewDownloader(app.OnDownloadResult, 10)
	app.downloader = d
	go app.downloader.Run()

	return app
}

func (a *App) Startup(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}

func (a *App) OnSearch(keyword string, pageIndex, pageSize int) model.BaseResponse {
	res, err := a.search(keyword, pageIndex, pageSize)
	if err != nil {
		return a.genError(err.Error())
	}

	if res.Code != consts.SearchSuccess {
		return a.genError(res.Info)
	}

	total, _ := strconv.Atoi(res.SongResultData.TotalCount)
	resp := &model.PageRes{
		Total: total,
		Items: res.SongResultData.Result,
	}

	return a.genOk(resp)
}

func (a *App) OnDownload(sourceType string, downloadItemsJson string) model.BaseResponse {
	var items []model.DownloadItem
	err := json.Unmarshal([]byte(downloadItemsJson), &items)
	if err != nil {
		return a.genError(err.Error())
	}

	path := ""
	downloadLrc := false
	downloadCover := false
	setting, _ := a.getSetting()
	if setting != nil && len(setting.SavePath) > 0 {
		path = setting.SavePath
		downloadLrc = setting.DownloadLrc
		downloadCover = setting.DownloadCover
	} else {
		path, err = runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			DefaultDirectory:           "",
			DefaultFilename:            "",
			Title:                      "选择保存路径",
			Filters:                    nil,
			ShowHiddenFiles:            false,
			CanCreateDirectories:       false,
			ResolvesAliases:            false,
			TreatPackagesAsDirectories: false,
		})
		if err != nil {
			return a.genError("选择文件夹失败")
		}
		if len(path) <= 0 {
			return a.genError("取消下载")
		}
	}

	for _, item := range items {
		a.download(sourceType, path, item, downloadLrc, downloadCover)
		a.log(fmt.Sprintf("[%s]添加成功 ", item.Name))
	}

	return a.genOk(nil)
}

func (a *App) OnSelectSavePath() model.BaseResponse {
	existPath := ""
	setting, err := a.getSetting()
	if err != nil {
		a.log("getsetting err: "+err.Error())
		return a.genError("设置失败")
	}
	if setting != nil {
		existPath = setting.SavePath
	}

	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:           existPath,
		DefaultFilename:            "",
		Title:                      "选择保存路径",
		Filters:                    nil,
		ShowHiddenFiles:            false,
		CanCreateDirectories:       false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	})
	if err != nil {
		a.log("getsetting err: "+err.Error())
		return a.genError("设置失败")
	}

	if len(path) <= 0 {
		return a.genError("未选择路径")
	}

	return a.genOk(path)
}

func (a *App) getSetting() (*model.Setting, error) {
	if _, err := os.Stat(a.configPath); os.IsNotExist(err) {
		os.MkdirAll(a.configPath, os.ModePerm)
	}

	filename := a.configPath + "/conf.yaml"
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filename)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			setting := model.Setting{
				SavePath:      "",
				DownloadLrc:   false,
				DownloadCover: false,
			}
			settingJson, _ := yaml.Marshal(setting)
			_, err = file.Write(settingJson)
			if err != nil {
				return nil, err
			}

			return &setting, nil
		} else {
			return nil, err
		}
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var setting model.Setting
	err = yaml.Unmarshal(content, &setting)
	if err != nil {
		return nil, err
	}

	return &setting, nil
}

func (a *App) setSetting(setting model.Setting) error {
	if _, err := os.Stat(a.configPath); os.IsNotExist(err) {
		os.MkdirAll(a.configPath, os.ModePerm)
	}

	filename := a.configPath + "/conf.yaml"
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer file.Close()

			setting := model.Setting{
				SavePath:      "",
				DownloadLrc:   false,
				DownloadCover: false,
			}
			settingJson, _ := yaml.Marshal(setting)
			_, err = file.Write(settingJson)
			if err != nil {
				return err
			}

			return nil
		} else {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, consts.DefaultPermOpen)
	if err != nil {
		return err
	}
	defer file.Close()

	settingJson, _ := yaml.Marshal(setting)
	_, err = file.Write(settingJson)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) OnGetSetting() model.BaseResponse {
	setting, err := a.getSetting()
	if err != nil {
		a.log("读取配置文件错误: " + err.Error())
		return a.genError("读取配置文件错误")
	}

	return a.genOk(setting)
}

func (a *App) OnSetSetting(settingStr string) model.BaseResponse {
	var setting model.Setting
	err := json.Unmarshal([]byte(settingStr), &setting)
	if err != nil {
		a.log("保存配置失败: " + err.Error())
		return a.genError("保存配置失败")
	}

	err = a.setSetting(setting)
	if err != nil {
		a.log("保存配置失败: " + err.Error())
		return a.genError("保存配置失败")
	}

	return a.genOk(setting)
}

func (a *App) OnDownloadResult(res model.BaseResponse) {
	item := res.Data.(model.DownloadQueueItem)
	a.log(fmt.Sprintf("[%s]%s", item.Name, res.Message))
	a.pushDownloadResult(res)
}

func (a *App) search(keyword string, pageIndex, pageSize int) (*model.SearchRes, error) {
	// http://pd.musicapp.migu.cn/MIGUM2.0/v1.0/content/search_all.do?ua=Android_migu&version=5.0.1&pageNo=1&pageSize=10&text=周杰伦&searchSwitch=
	url := fmt.Sprintf(a.searchUrl, pageIndex, pageSize, keyword)

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

func (a *App) download(sourceType, path string, item model.DownloadItem, downloadLrc, downloadCover bool) {
	_sourceType := consts.SourceType_HQ
	if sourceType == "SQ" {
		_sourceType = consts.SourceType_SQ
	}

	if path[len(path)-1] != '/' {
		path += "/"
	}

	path += item.Name + consts.SourceType2FileExt[_sourceType]
	url := fmt.Sprintf(a.downloadUrl, string(_sourceType), item.ContentId)

	a.downloader.Push(a.ctx, model.DownloadQueueItem{
		DownloadItem:  item,
		Path:          path,
		Url:           url,
		DownloadLrc:   downloadLrc,
		DownloadCover: downloadCover,
	})
}

func (a *App) genOk(data interface{}) model.BaseResponse {
	return model.BaseResponse{
		Code:    0,
		Message: "",
		Data:    data,
	}
}

func (a *App) genError(message string) model.BaseResponse {
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

func (a *App) Stop(ctx context.Context) {
	a.cancel()
	a.downloader.Stop()
}
