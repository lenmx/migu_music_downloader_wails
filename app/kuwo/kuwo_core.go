package kuwo

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"migu_music_downloader_wails/app/model"
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
