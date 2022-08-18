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
	"strings"
)

type AppQQ struct {
	ctx    context.Context
	cancel context.CancelFunc

	searchUrl   string
	downloadUrl string
	configPath  string

	core       *AppQQCore
	downloader *util.Downloader
}

func NewAppQQ(searchUrl, downloadUrl, configPath string) *AppQQ {
	app := &AppQQ{
		searchUrl:   searchUrl,
		downloadUrl: downloadUrl,
		configPath:  configPath,
	}

	// app core
	core, err := NewAppQQCore()
	if err != nil {
		fmt.Println("初始化失败", err)
	}
	app.core = core

	// downloader
	d := util.NewDownloader(app.onDownloadResult, 10)
	app.downloader = d
	go app.downloader.Run()

	return app
}

func (a *AppQQ) Startup(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}

func (a *AppQQ) OnSearch(keyword string, pageIndex, pageSize int) model.BaseResponse {
	res, err := a.search(keyword, pageIndex, pageSize)
	if err != nil {
		return a.genError(err.Error())
	}

	if res.Code != 0 {
		return a.genError("搜索失败")
	}

	items := []model.SearchQQItem{}
	for _, item := range res.MusicSearchSearchCgiService.Data.Body.Song.List {
		singers := []string{}
		for _, singer := range item.Singer {
			singers = append(singers, singer.Title)
		}

		albums := []string{item.Album.Title}

		fileInfos := map[string]model.FileInfo{}
		//if item.File.SizeHires > 0 {
		//	fileInfos["hires"] = model.FileInfo{
		//		Code:   "RS01",
		//		Format: "flac",
		//		Tips:   "高解析无损 Hi-Res",
		//		Size:   item.File.SizeHires,
		//	}
		//}
		if item.File.SizeFlac > 0 {
			fileInfos["flac"] = model.FileInfo{
				Code:   "F000",
				Format: "flac",
				Tips:   "无损品质 FLAC",
				Size:   item.File.SizeFlac,
			}
		}
		//if item.File.Size320Mp3 > 0 {
		//	fileInfos["320mp3"] = model.FileInfo{
		//		Code:   "M800",
		//		Format: "mp3",
		//		Tips:   "超高品质 320kbps",
		//		Size:   item.File.Size320Mp3,
		//	}
		//}
		if item.File.Size128Mp3 > 0 {
			fileInfos["mp3"] = model.FileInfo{
				Code:   "M500",
				Format: "mp3",
				Tips:   "标准品质 128kbps",
				Size:   item.File.SizeHires,
			}
		}

		items = append(items, model.SearchQQItem{
			ContentId:  item.Id,
			DocId:      item.Docid,
			Mid:        item.Mid,
			Name:       item.Name,
			Title:      item.Title,
			TimePublic: item.TimePublic,
			Singers:    singers,
			Albums:     albums,
			LrcUrl:     "",
			Cover:      "",
			File:       item.File,
			FileInfos:  fileInfos,
		})
	}

	resp := &model.PageRes{
		Total: res.MusicSearchSearchCgiService.Data.Meta.Sum,
		Items: items,
	}

	return a.genOk(resp)
}

func (a *AppQQ) OnDownload(sourceType string, downloadItemsJson string) model.BaseResponse {
	a.log("receive item: " + downloadItemsJson)
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
	if path[len(path)-1] != '/' {
		path += "/"
	}

	for _, item := range items {
		fileInfo := item.FileInfos["mp3"]
		if sourceType == "SQ" {
			fileInfo = item.FileInfos["flac"]
		}
		filename := getMusicFileName(fileInfo.Code, item.File.MediaMid, fileInfo.Format)
		url, err := a.core.getMusicUrl(filename)
		if err != nil || (!strings.Contains(url, "qqmusic.qq.com") && strings.Contains(url, `"title":"Not Found"`)) {
			url, err = a.core.reGetMusicUrl(filename, item.Mid)
		}
		if err != nil {
			a.log(fmt.Sprintf("[%s]下载地址解析失败", item.Name))
			continue
		}

		item.Url = url
		filePath := path + item.Name + "." + fileInfo.Format
		if downloadLrc {
			item.LrcContent, _ = a.core.getLyric(item.Mid)
		}
		a.download(filePath, item, downloadLrc, downloadCover)
		a.log(fmt.Sprintf("[%s]添加成功", item.Name))
	}

	return a.genOk(nil)
}

func (a *AppQQ) OnSelectSavePath() model.BaseResponse {
	existPath := ""
	setting, err := a.getSetting()
	if err != nil {
		a.log("getsetting err: " + err.Error())
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
		a.log("getsetting err: " + err.Error())
		return a.genError("设置失败")
	}

	if len(path) <= 0 {
		return a.genError("未选择路径")
	}

	return a.genOk(path)
}

func (a *AppQQ) getSetting() (*model.Setting, error) {
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

func (a *AppQQ) setSetting(setting model.Setting) error {
	if _, err := os.Stat(a.configPath); os.IsNotExist(err) {
		os.MkdirAll(a.configPath, os.ModePerm)
	}

	var file *os.File
	filename := a.configPath + "/conf.yaml"
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filename)
			if err != nil {
				return err
			}
			defer file.Close()
		} else {
			return err
		}
	}

	if file == nil {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, consts.DefaultPermOpen)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	settingJson, _ := yaml.Marshal(setting)
	_, err = file.Write(settingJson)
	if err != nil {
		return err
	}

	return nil
}

func (a *AppQQ) OnGetSetting() model.BaseResponse {
	setting, err := a.getSetting()
	if err != nil {
		a.log("读取配置文件错误: " + err.Error())
		return a.genError("读取配置文件错误")
	}

	return a.genOk(setting)
}

func (a *AppQQ) OnSetSetting(settingStr string) model.BaseResponse {
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

func (a *AppQQ) onDownloadResult(res model.BaseResponse) {
	item := res.Data.(model.DownloadQueueItem)
	a.log(fmt.Sprintf("[%s]%s", item.Name, res.Message))
	a.pushDownloadResult(res)
}

func (a *AppQQ) search(keyword string, pageIndex, pageSize int) (*model.SearchQQRes, error) {
	url := fmt.Sprintf(a.searchUrl, pageIndex, pageSize, keyword)
	url = "https://u.y.qq.com/cgi-bin/musicu.fcg"
	data := map[string]interface{}{
		"comm": map[string]interface{}{"ct": "19", "cv": "1845"},
		"music.search.SearchCgiService": map[string]interface{}{
			"method": "DoSearchForQQMusicDesktop",
			"module": "music.search.SearchCgiService",
			"param":  map[string]interface{}{"query": keyword, "num_per_page": pageSize, "page_num": pageIndex},
		},
	}
	res, err := resty.New().R().SetBody(data).Post(url)
	if err != nil {
		return nil, err
	}

	var resp model.SearchQQRes
	err = json.Unmarshal(res.Body(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (a *AppQQ) download(path string, item model.DownloadItem, downloadLrc, downloadCover bool) {
	a.downloader.Push(a.ctx, model.DownloadQueueItem{
		DownloadItem:  item,
		Path:          path,
		DownloadLrc:   downloadLrc,
		DownloadCover: downloadCover,
	})
}

func (a *AppQQ) genOk(data interface{}) model.BaseResponse {
	return model.BaseResponse{
		Code:    0,
		Message: "",
		Data:    data,
	}
}

func (a *AppQQ) genError(message string) model.BaseResponse {
	return model.BaseResponse{
		Code:    -1,
		Message: message,
		Data:    nil,
	}
}

func (a *AppQQ) log(message string) {
	util.Log(a.ctx, message)
}

func (a *AppQQ) pushDownloadResult(response model.BaseResponse) {
	util.PushDownloadResult(a.ctx, response)
}

func (a *AppQQ) Stop(ctx context.Context) {
	a.cancel()
	a.downloader.Stop()
}
