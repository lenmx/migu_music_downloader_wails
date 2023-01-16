package app

//
//import (
//	"context"
//	"embed"
//	"encoding/json"
//	"fmt"
//	"github.com/go-resty/resty/v2"
//	"github.com/wailsapp/wails/v2/pkg/runtime"
//	"gopkg.in/yaml.v2"
//	"io/ioutil"
//	"migu_music_downloader_wails/app/consts"
//	"migu_music_downloader_wails/app/i18n"
//	"migu_music_downloader_wails/app/model"
//	"migu_music_downloader_wails/app/util"
//	"os"
//	"strings"
//)
//
//type AppQQ struct {
//	ctx    context.Context
//	cancel context.CancelFunc
//	i18n   *i18n.I18n
//
//	searchUrl   string
//	downloadUrl string
//	configPath  string
//
//	core       *AppQQCore
//	downloader *util.Downloader
//}
//
////go:embed i18n/locale.*.json
//var LocaleFS embed.FS
//
//func NewAppQQ(searchUrl, downloadUrl, configPath string) (*AppQQ, error) {
//	app := &AppQQ{
//		i18n:        i18n.New("en"),
//		searchUrl:   searchUrl,
//		downloadUrl: downloadUrl,
//		configPath:  configPath,
//	}
//
//	_, err := app.i18n.LoadFile(LocaleFS, "i18n/locale.zh.json", "zh")
//	if err != nil {
//		fmt.Println("初始化失败", err)
//		return nil, err
//	}
//
//	_, err = app.i18n.LoadFile(LocaleFS, "i18n/locale.en.json", "en")
//	if err != nil {
//		fmt.Println("初始化失败", err)
//		return nil, err
//	}
//
//	// app core
//	core, err := NewAppQQCore()
//	if err != nil {
//		fmt.Println("初始化失败", err)
//		return nil, err
//	}
//	app.core = core
//
//	// downloader
//	d := util.NewDownloader(app.onDownloadResult, 10)
//	app.downloader = d
//	go app.downloader.Run()
//
//	return app, nil
//}
//
//func (a *AppQQ) Startup(ctx context.Context) {
//	a.ctx, a.cancel = context.WithCancel(ctx)
//}
//
//func (a *AppQQ) OnSearch(keyword string, pageIndex, pageSize int) model.BaseResponse {
//	res, err := a.search(keyword, pageIndex, pageSize)
//	if err != nil {
//		return a.GenError(err.Error())
//	}
//
//	if res.Code != 0 {
//		return a.GenError(a.TR("SearchFail"))
//	}
//
//	items := []model.SearchQQItem{}
//	for _, item := range res.MusicSearchSearchCgiService.Data.Body.Song.List {
//		singers := []string{}
//		for _, singer := range item.Singer {
//			singers = append(singers, singer.Title)
//		}
//
//		albums := []string{item.Album.Title}
//
//		fileInfos := map[string]model.FileInfo{}
//		//if item.File.SizeHires > 0 {
//		//	fileInfos["hires"] = model.FileInfo{
//		//		Code:   "RS01",
//		//		Format: "flac",
//		//		Tips:   "高解析无损 Hi-Res",
//		//		Size:   item.File.SizeHires,
//		//	}
//		//}
//		if item.File.SizeFlac > 0 {
//			fileInfos["flac"] = model.FileInfo{
//				Code:   "F000",
//				Format: "flac",
//				Tips:   "无损品质 FLAC",
//				Size:   item.File.SizeFlac,
//			}
//		}
//		//if item.File.Size320Mp3 > 0 {
//		//	fileInfos["320mp3"] = model.FileInfo{
//		//		Code:   "M800",
//		//		Format: "mp3",
//		//		Tips:   "超高品质 320kbps",
//		//		Size:   item.File.Size320Mp3,
//		//	}
//		//}
//		if item.File.Size128Mp3 > 0 {
//			fileInfos["mp3"] = model.FileInfo{
//				Code:   "M500",
//				Format: "mp3",
//				Tips:   "标准品质 128kbps",
//				Size:   item.File.SizeHires,
//			}
//		}
//
//		items = append(items, model.SearchQQItem{
//			ContentId:  item.Id,
//			DocId:      item.Docid,
//			Mid:        item.Mid,
//			Name:       item.Name,
//			Title:      item.Title,
//			TimePublic: item.TimePublic,
//			Singers:    singers,
//			Albums:     albums,
//			LrcUrl:     "",
//			Cover:      "",
//			File:       item.File,
//			FileInfos:  fileInfos,
//		})
//	}
//
//	resp := &model.PageRes{
//		Total: res.MusicSearchSearchCgiService.Data.Meta.Sum,
//		Items: items,
//	}
//
//	return a.GenOk(resp)
//}
//
//func (a *AppQQ) OnDownload(sourceType string, downloadItemsJson string) model.BaseResponse {
//	a.Log("receive item: " + downloadItemsJson)
//	var items []model.DownloadItem
//	err := json.Unmarshal([]byte(downloadItemsJson), &items)
//	if err != nil {
//		return a.GenError(err.Error())
//	}
//
//	path := ""
//	downloadLrc := false
//	downloadCover := false
//	setting, _ := a.GetSetting()
//	if setting != nil && len(setting.SavePath) > 0 {
//		path = setting.SavePath
//		downloadLrc = setting.DownloadLrc
//		downloadCover = setting.DownloadCover
//	} else {
//		path, err = runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
//			DefaultDirectory:           "",
//			DefaultFilename:            "",
//			Title:                      a.TR("ChooseFileSavePath"),
//			Filters:                    nil,
//			ShowHiddenFiles:            false,
//			CanCreateDirectories:       false,
//			ResolvesAliases:            false,
//			TreatPackagesAsDirectories: false,
//		})
//		if err != nil {
//			return a.GenError(a.TR("ChooseFileSavePathFail"))
//		}
//		if len(path) <= 0 {
//			return a.GenError(a.TR("CancelDownload"))
//		}
//	}
//	if path[len(path)-1] != '/' {
//		path += "/"
//	}
//
//	for _, item := range items {
//		fileInfo := item.FileInfos["mp3"]
//		if sourceType == "SQ" {
//			fileInfo = item.FileInfos["flac"]
//		}
//		filename := getMusicFileName(fileInfo.Code, item.File.MediaMid, fileInfo.Format)
//		url, err := a.core.getMusicUrl(filename)
//		if err != nil || (!strings.Contains(url, "qqmusic.qq.com") && strings.Contains(url, `"title":"Not Found"`)) {
//			url, err = a.core.reGetMusicUrl(filename, item.Mid)
//		}
//		if err != nil {
//			a.Log(fmt.Sprintf("[%s]%s", item.Name, a.TR("DownloadUrlParseFail")))
//			continue
//		}
//
//		item.Url = url
//		filePath := path + item.Name + "." + fileInfo.Format
//		if downloadLrc {
//			item.LrcContent, _ = a.core.getLyric(item.Mid)
//		}
//		a.download(filePath, item, downloadLrc, downloadCover)
//		a.Log(fmt.Sprintf("[%s]%s", item.Name, a.TR("AddToDownloadCenterSuccess")))
//	}
//
//	return a.GenOk(nil)
//}
//
//func (a *AppQQ) OnSelectSavePath() model.BaseResponse {
//	existPath := ""
//	setting, err := a.GetSetting()
//	if err != nil {
//		a.Log("getsetting err: " + err.Error())
//		return a.GenError(a.TR("SettingFail"))
//	}
//	if setting != nil {
//		existPath = setting.SavePath
//	}
//
//	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
//		DefaultDirectory:           existPath,
//		DefaultFilename:            "",
//		Title:                      a.TR("ChooseFileSavePathFail"),
//		Filters:                    nil,
//		ShowHiddenFiles:            false,
//		CanCreateDirectories:       false,
//		ResolvesAliases:            false,
//		TreatPackagesAsDirectories: false,
//	})
//	if err != nil {
//		a.Log("getsetting err: " + err.Error())
//		return a.GenError(a.TR("SettingFail"))
//	}
//
//	if len(path) <= 0 {
//		return a.GenError(a.TR("NoSavePathChoose"))
//	}
//
//	return a.GenOk(path)
//}
//
//func (a *AppQQ) GetSetting() (*model.Setting, error) {
//	if _, err := os.Stat(a.configPath); os.IsNotExist(err) {
//		os.MkdirAll(a.configPath, os.ModePerm)
//	}
//
//	filename := a.configPath + "/conf.yaml"
//	_, err := os.Stat(filename)
//	if err != nil {
//		if os.IsNotExist(err) {
//			file, err := os.Create(filename)
//			if err != nil {
//				return nil, err
//			}
//			defer file.Close()
//
//			setting := model.Setting{
//				SavePath:      "",
//				DownloadLrc:   false,
//				DownloadCover: false,
//			}
//			settingJson, _ := yaml.Marshal(setting)
//			_, err = file.Write(settingJson)
//			if err != nil {
//				return nil, err
//			}
//
//			return &setting, nil
//		} else {
//			return nil, err
//		}
//	}
//
//	content, err := ioutil.ReadFile(filename)
//	if err != nil {
//		return nil, err
//	}
//
//	var setting model.Setting
//	err = yaml.Unmarshal(content, &setting)
//	if err != nil {
//		return nil, err
//	}
//
//	return &setting, nil
//}
//
//func (a *AppQQ) SetSetting(setting model.Setting) error {
//	if _, err := os.Stat(a.configPath); os.IsNotExist(err) {
//		os.MkdirAll(a.configPath, os.ModePerm)
//	}
//
//	var file *os.File
//	filename := a.configPath + "/conf.yaml"
//	_, err := os.Stat(filename)
//	if err != nil {
//		if os.IsNotExist(err) {
//			file, err = os.Create(filename)
//			if err != nil {
//				return err
//			}
//			defer file.Close()
//		} else {
//			return err
//		}
//	}
//
//	if file == nil {
//		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, consts.DefaultPermOpen)
//		if err != nil {
//			return err
//		}
//		defer file.Close()
//	}
//
//	settingJson, _ := yaml.Marshal(setting)
//	_, err = file.Write(settingJson)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (a *AppQQ) OnGetSetting() model.BaseResponse {
//	setting, err := a.GetSetting()
//	if err != nil {
//		a.Log("读取配置文件错误: " + err.Error())
//		return a.GenError(a.TR("ConfigFileParseFail"))
//	}
//
//	return a.GenOk(setting)
//}
//
//func (a *AppQQ) OnSetSetting(settingStr string) model.BaseResponse {
//	var setting model.Setting
//	err := json.Unmarshal([]byte(settingStr), &setting)
//	if err != nil {
//		a.Log("保存配置失败: " + err.Error())
//		return a.GenError(a.TR("SettingFail"))
//	}
//
//	err = a.SetSetting(setting)
//	if err != nil {
//		a.Log("保存配置失败: " + err.Error())
//		return a.GenError(a.TR("SettingFail"))
//	}
//
//	if a.i18n.GetDefaultLanguage() != setting.Language {
//		a.i18n.SetDefaultLanguage(setting.Language)
//		runtime.WindowReloadApp(a.ctx)
//	}
//	return a.GenOk(setting)
//}
//
//func (a *AppQQ) GetI18nSource(key string) model.I18nSourceMap {
//	return model.I18nSourceMap{
//		CurrentLang: a.i18n.GetLangName(),
//		Sources:     a.i18n.GetI18nSource(),
//	}
//}
//
//func (a *AppQQ) onDownloadResult(res model.BaseResponse) {
//	item := res.Data.(model.DownloadQueueItem)
//	a.Log(fmt.Sprintf("[%s]%s", item.Name, res.Message))
//	a.pushDownloadResult(res)
//}
//
//func (a *AppQQ) search(keyword string, pageIndex, pageSize int) (*model.SearchQQRes, error) {
//	url := fmt.Sprintf(a.searchUrl, pageIndex, pageSize, keyword)
//	url = "https://u.y.qq.com/cgi-bin/musicu.fcg"
//	data := map[string]interface{}{
//		"comm": map[string]interface{}{"ct": "19", "cv": "1845"},
//		"music.search.SearchCgiService": map[string]interface{}{
//			"method": "DoSearchForQQMusicDesktop",
//			"module": "music.search.SearchCgiService",
//			"param":  map[string]interface{}{"query": keyword, "num_per_page": pageSize, "page_num": pageIndex},
//		},
//	}
//	res, err := resty.New().R().SetBody(data).Post(url)
//	if err != nil {
//		return nil, err
//	}
//
//	var resp model.SearchQQRes
//	err = json.Unmarshal(res.Body(), &resp)
//	if err != nil {
//		return nil, err
//	}
//
//	return &resp, nil
//}
//
//func (a *AppQQ) download(path string, item model.DownloadItem, downloadLrc, downloadCover bool) {
//	a.downloader.Push(a.ctx, model.DownloadQueueItem{
//		DownloadItem:  item,
//		Path:          path,
//		DownloadLrc:   downloadLrc,
//		DownloadCover: downloadCover,
//	})
//}
//
//func (a *AppQQ) TR(key string) string {
//	return a.i18n.Parse("server." + key)
//}
//
//func (a *AppQQ) GenOk(data interface{}) model.BaseResponse {
//	return model.BaseResponse{
//		Code:    0,
//		Message: "",
//		Data:    data,
//	}
//}
//
//func (a *AppQQ) GenError(message string) model.BaseResponse {
//	return model.BaseResponse{
//		Code:    -1,
//		Message: message,
//		Data:    nil,
//	}
//}
//
//func (a *AppQQ) Log(message string) {
//	util.Log(a.ctx, message)
//}
//
//func (a *AppQQ) pushDownloadResult(response model.BaseResponse) {
//	util.PushDownloadResult(a.ctx, response)
//}
//
//func (a *AppQQ) Stop(ctx context.Context) {
//	a.cancel()
//	a.downloader.Stop()
//}
