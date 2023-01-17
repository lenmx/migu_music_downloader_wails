package app

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v2"
	"migu_music_downloader_wails/app/consts"
	"migu_music_downloader_wails/app/i18n"
	"migu_music_downloader_wails/app/model"
	"migu_music_downloader_wails/app/util"
	"os"
	"path"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	i18n   *i18n.I18n

	configFilename string
}

var GApp *App

//go:embed i18n/locale.*.json
var LocaleFS embed.FS

func NewApp(configPath string) (*App, error) {
	if GApp != nil {
		return GApp, nil
	}

	GApp = &App{
		i18n:           i18n.New(),
		configFilename: path.Join(configPath, "config.yml"),
	}

	util.CreateFileIfN(GApp.configFilename)

	_, err := GApp.i18n.LoadFile(LocaleFS, "i18n/locale.zh.json", "zh")
	if err != nil {
		fmt.Println("初始化失败", err)
		return nil, err
	}

	_, err = GApp.i18n.LoadFile(LocaleFS, "i18n/locale.en.json", "en")
	if err != nil {
		fmt.Println("初始化失败", err)
		return nil, err
	}

	return GApp, nil
}

func (a *App) Startup(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}

func (a *App) Stop(ctx context.Context) {
	a.cancel()
}

func (a *App) OnSelectSavePath() model.BaseResponse {
	existPath := ""
	setting, err := a.GetSetting()
	if err != nil {
		a.Log("getSetting err: " + err.Error())
		return a.GenError(a.TR("SettingFail"))
	}
	if setting != nil {
		existPath = setting.SavePath
	}

	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:           existPath,
		DefaultFilename:            "",
		Title:                      a.TR("ChooseFileSavePathFail"),
		Filters:                    nil,
		ShowHiddenFiles:            false,
		CanCreateDirectories:       false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	})
	if err != nil {
		a.Log("getSetting err: " + err.Error())
		return a.GenError(a.TR("SettingFail"))
	}

	if len(path) <= 0 {
		return a.GenError(a.TR("NoSavePathChoose"))
	}

	return a.GenOk(path)
}

func (a *App) GetSetting() (*model.Setting, error) {
	//_path, _ := filepath.Abs(a.configFilename)
	//a.Log("config file path: " + _path)
	content, err := os.ReadFile(a.configFilename)
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

func (a *App) SetSetting(setting model.Setting) error {
	file, err := os.OpenFile(a.configFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, consts.DefaultPermOpen)
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
	setting, err := a.GetSetting()
	if err != nil {
		a.Log("读取配置文件错误: " + err.Error())
		return a.GenError(a.TR("ConfigFileParseFail"))
	}

	return a.GenOk(setting)
}

func (a *App) OnSetSetting(settingStr string) model.BaseResponse {
	var setting model.Setting
	err := json.Unmarshal([]byte(settingStr), &setting)
	if err != nil {
		a.Log("保存配置失败: " + err.Error())
		return a.GenError(a.TR("SettingFail"))
	}

	err = a.SetSetting(setting)
	if err != nil {
		a.Log("保存配置失败: " + err.Error())
		return a.GenError(a.TR("SettingFail"))
	}

	if a.i18n.GetDefaultLanguage() != setting.Language {
		a.i18n.SetDefaultLanguage(setting.Language)
		runtime.WindowReloadApp(a.ctx)
	}
	return a.GenOk(setting)
}

func (a *App) GetI18nSource(key string) model.I18nSourceMap {
	return model.I18nSourceMap{
		CurrentLang: a.i18n.GetLangName(),
		Sources:     a.i18n.GetI18nSource(),
	}
}

func (a *App) TR(key string) string {
	return a.i18n.Parse("server." + key)
}

func (a *App) GenOk(data interface{}) model.BaseResponse {
	return model.BaseResponse{
		Code:    0,
		Message: "",
		Data:    data,
	}
}

func (a *App) GenError(message string) model.BaseResponse {
	return model.BaseResponse{
		Code:    -1,
		Message: message,
		Data:    nil,
	}
}

func (a *App) Log(message string) {
	util.Log(a.ctx, message)
}
