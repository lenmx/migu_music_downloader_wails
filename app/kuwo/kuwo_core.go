package kuwo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"migu_music_downloader_wails/app/model"
	"migu_music_downloader_wails/app/util"
	url2 "net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

type AppKuwoCore struct {
	headers map[string]string
}

func NewAppKuwoCore() *AppKuwoCore {
	core := &AppKuwoCore{
		headers: map[string]string{},
	}

	return core
}

func (a *AppKuwoCore) Search(keyword string, pageIndex, pageSize int) (*model.KuwoSearchRes, error) {
	keyword = url2.QueryEscape(keyword)
	url := fmt.Sprintf("http://search.kuwo.cn/r.s?client=kt&all=%s&pn=%d&rn=%d&uid=221260053&ver=kwplayer_ar_99.99.99.99&vipver=1&ft=music&cluster=0&strategy=2012&encoding=utf8&rformat=json&vermerge=1&mobi=1", keyword, pageIndex, pageSize)
	res, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}

	var resp model.KuwoSearchRes
	err = json.Unmarshal(res.Body(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (a *AppKuwoCore) GetDownloadUrl(musicId string) string {
	url := fmt.Sprintf("http://antiserver.kuwo.cn/anti.s?type=convert_url&rid=%s&format=mp3|acc&response=url", musicId)

	cli := resty.New().R()
	cli.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := cli.Get(url)
	if err != nil {
		return ""
	}

	tmp := string(resp.Body())
	if tmp != "" && strings.HasPrefix(tmp, "http") {
		return tmp
	}

	return ""
}

func (a *AppKuwoCore) GetDownloadUrl1(musicId string) string {
	url := fmt.Sprintf("https://www.kuwo.cn/api/v1/www/music/playUrl?mid=%s&type=flac&httpsStatus=1&reqId=80b33650-8a62-11ed-a069-8d99eba73f2a", musicId)

	cli := resty.New().R()
	cli.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := cli.Get(url)
	if err != nil {
		return ""
	}

	tmp := resp.Body()
	var getMusicUrlRes model.KuwoGetMusicUrlRes
	json.Unmarshal(tmp, &getMusicUrlRes)
	fmt.Println(getMusicUrlRes)
	if getMusicUrlRes.Code == 200 && getMusicUrlRes.Data.Url != "" {
		return getMusicUrlRes.Data.Url
	}

	return ""
}

func (a *AppKuwoCore) GetLrcUrl(musicId string) string {
	url := fmt.Sprintf("http://m.kuwo.cn/newh5/singles/songinfoandlrc?musicId=%s", musicId)
	return url
}

func (a *AppKuwoCore) GetPicUrl(musicId string) string {
	url := fmt.Sprintf("http://artistpicserver.kuwo.cn/pic.web?corp=kuwo&type=rid_pic&pictype=url&content=list&size=640&rid=%s", musicId)
	return url
}

func (a *AppKuwoCore) ProcessMp3(url, mp3Filename string) error {
	mp3Res, err := resty.New().R().Get(url)
	if err != nil {
		return err
	}

	// MP3文件不存在
	if len(mp3Res.Body()) < 1000 {
		return errors.New("file not exist")
	}

	file, err := os.OpenFile(mp3Filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err == nil {
		defer file.Close()
		file.Write(mp3Res.Body())
	}

	return nil
}

func (a *AppKuwoCore) ProcessLrc(url, mp3Filename string) {
	res, err := resty.New().R().Get(url)
	if err != nil {
		return
	}

	var data model.KuwoLrcRes
	json.Unmarshal(res.Body(), &data)

	var sb strings.Builder
	for _, item := range data.Data.Lrclist {
		time := fmt.Sprintf("[%s]", convertTime(item.Time))
		sb.WriteString(time)
		sb.WriteString(item.LineLyric + "\n")
	}

	filename := path.Join(path.Dir(mp3Filename), util.GetFilename(mp3Filename)+".lrc")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err == nil {
		defer file.Close()
		file.WriteString(sb.String())
	}
}

func (a *AppKuwoCore) ProcessPic(url, mp3Filename string) {
	r := resty.New().R()
	picUrlRes, err := r.Get(url)
	picUrl := strings.Trim(picUrlRes.String(), " ")

	if err == nil && strings.HasPrefix(picUrl, "http") {
		picRes, err := r.Get(picUrl)
		if err == nil {
			filename := path.Join(path.Dir(mp3Filename), util.GetFilename(mp3Filename)+path.Ext(picUrl))
			file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err == nil {
				defer file.Close()
				file.Write(picRes.Body())
			}
		}
	}
}

func convertTime(str string) string {
	timeArr := strings.Split(str, ".")
	ms := time2MS(timeArr[0])

	return fmt.Sprintf("%s.%s", ms, timeArr[1])
}

func time2MS(str string) string {
	t, _ := strconv.Atoi(str)
	minute := t / 60
	second := t % 60

	return fmt.Sprintf("%02d:%02d", minute, second)
}
