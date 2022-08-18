package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"migu_music_downloader_wails/app/model"
	"migu_music_downloader_wails/app/util"
	"strings"
)

type AppQQCore struct {
	mkey    string
	mqq     string
	headers map[string]string
}

func NewAppQQCore() (*AppQQCore, error) {
	cookie, err := getCookie()
	if err != nil {
		return nil, err
	}

	ok, mkey, mqq := decryptCookie(cookie)
	if !ok {
		return nil, errors.New("解析cookie失败")
	}

	core := &AppQQCore{
		mkey:    mkey,
		mqq:     mqq,
		headers: getHeader(),
	}

	return core, nil
}

func (a *AppQQCore) getMusicUrl(filename string) (string, error) {
	u := "http://8.136.185.193/api/MusicLink/link"
	d := encryptText(filename, a.mqq)
	dd := fmt.Sprintf(`"%s"`, d)

	res, err := resty.New().R().SetBody(dd).SetHeader("Content-Type", "application/json;charset=utf-8").Post(u)
	if err != nil {
		return "", err
	}
	return string(res.Body()), err
}

func (a *AppQQCore) getLyric(mid string) (string, error) {
	url := fmt.Sprintf("https://c.y.qq.com/lyric/fcgi-bin/fcg_query_lyric_new.fcg?songmid=%s&g_tk=5381", mid)
	d, err := a.getQQServerCallback(url, true, nil)
	if err != nil {
		return "", err
	}

	var resp model.LyricRes
	err = json.Unmarshal([]byte(d[18:len(d)-1]), &resp)
	if err != nil {
		return "", err
	}

	if resp.Retcode != 0 {
		return "", errors.New("歌词解析失败")
	}

	content, err := base64.StdEncoding.DecodeString(resp.Lyric)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (a *AppQQCore) getQQServerCallback(url string, methodGet bool, data interface{}) (string, error) {
	header := getHeader()
	header["cookie"] = fmt.Sprintf("qqmusic_key=%s;qqmusic_uin=%s;", a.mkey, a.mqq)

	var (
		resp *resty.Response
		err  error
	)
	if methodGet {
		resp, err = resty.New().R().SetHeaders(header).Get(url)
	} else {
		cli := resty.New().R().SetHeaders(header)
		if data != nil {
			d, _ := json.Marshal(data)
			cli.SetBody(d)
		}

		resp, err = cli.Post(url)
	}

	if err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}

func (a *AppQQCore) reGetMusicUrl(filename, songMid string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("下载地址解析失败")
			return
		}
	}()
	url := "https://u.y.qq.com/cgi-bin/musicu.fcg"
	data := map[string]interface{}{
		"comm": map[string]interface{}{
			"ct": "19",
			"cv": "1777",
		},
		"queryvkey": map[string]interface{}{
			"method": "CgiGetVkey",
			"module": "vkey.GetVkeyServer",
			"param": map[string]interface{}{
				"uin":      a.mqq,
				"guid":     "QMD50",
				"referer":  "y.qq.com",
				"songtype": []int{1},
				"filename": []string{filename},
				"songmid":  []string{songMid},
			},
		},
	}
	d, err := a.getQQServerCallback(url, false, data)
	if err != nil {
		return "", err
	}

	var resp model.CgiGetVkey
	err = json.Unmarshal([]byte(d), &resp)
	if err != nil {
		return "", err
	}

	musicUrl := fmt.Sprintf("http://ws.stream.qqmusic.qq.com/%s&fromtag=140", resp.QueryVKey.Data.MidUrlInfo[0].PUrl)
	return musicUrl, nil
}

func getMusicFileName(code, mid, format string) string {
	return fmt.Sprintf("%s%s.%s", code, mid, format)
}

func decryptCookie(text string) (bool, string, string) {
	text = strings.ReplaceAll(text, "-", "")
	text = strings.ReplaceAll(text, "|", "")
	if len(text) < 10 || strings.Index(text, "%") == -1 {
		return false, "", ""
	}

	split := strings.Split(text, "%")
	key := split[0]
	iv := key[:8]

	str := split[1]
	strRaw, _ := base64.StdEncoding.DecodeString(str)
	qq := util.DesDecrypt(string(strRaw), iv)
	if len(qq) < 8 {
		qq += "QMD"
	}

	keyRaw, _ := base64.StdEncoding.DecodeString(key)
	mkey := util.DesDecrypt(string(keyRaw), qq[:8])
	return true, mkey, qq
}

func encryptText(text, qq string) string {
	key := ("QMD" + qq)[:8]
	return util.DesEncrypt(text, key)
}

func decryptText(text, qq string) string {
	data := strings.ReplaceAll(text, "-", "")
	key := ("QMD" + qq)[:8]
	return util.DesDecrypt(data, key)
}

func getHeader() map[string]string {
	return map[string]string{
		"user-agent":   "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
		"content-type": "application/json; charset=UTF-8",
		"referer":      "https://y.qq.com/portal/profile.html",
	}
}

func getCookie() (string, error) {
	uid := "822a3b85-a5c9-438e-a277-a8da412e8265"
	systemVersion := "1.7.2"
	versionCode := "76"
	deviceBrand := "360"
	deviceModel := "QK1707-A01"
	appVersion := "7.1.2"
	encIP := encryptText(fmt.Sprintf("%s%s%s%s%s%s", uid, deviceModel, deviceBrand, systemVersion, appVersion, versionCode), "F*ckYou!")

	d := map[string]string{
		"appVersion":    appVersion,
		"deviceBrand":   deviceBrand,
		"deviceModel":   deviceModel,
		"ip":            encIP,
		"systemVersion": systemVersion,
		"uid":           uid,
		"versionCode":   versionCode,
	}

	u := "http://8.136.185.193/api/Cookies"
	res, err := resty.New().R().SetBody(d).Get(u)
	if err != nil {
		return "", err
	}
	return string(res.Body()), nil
}
