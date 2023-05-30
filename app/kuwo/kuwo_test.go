package kuwo

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"migu_music_downloader_wails/app/model"
	"migu_music_downloader_wails/app/util"
	"os"
	"path"
	"strings"
	"testing"
)

type GetMusicUrlRes struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	ReqId string `json:"reqId"`
	Data  struct {
		Url string `json:"url"`
	} `json:"data"`
	ProfileId string `json:"profileId"`
	CurTime   int64  `json:"curTime"`
	Success   bool   `json:"success"`
}

func TestDownload(t *testing.T) {
	musicId := "440616"
	urlBase := fmt.Sprintf("https://www.kuwo.cn/api/v1/www/music/playUrl?mid=%s&type=flac&httpsStatus=1&reqId=80b33650-8a62-11ed-a069-8d99eba73f2a", musicId)

	cli := resty.New().R()
	cli.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := cli.Get(urlBase)
	if err != nil {
		t.Error(err)
		return
	}

	tmp := resp.Body()
	var getMusicUrlRes GetMusicUrlRes
	json.Unmarshal(tmp, &getMusicUrlRes)
	fmt.Println(getMusicUrlRes)
	if getMusicUrlRes.Code == 200 && getMusicUrlRes.Data.Url != "" {

		getMusicRes, err := cli.Get(getMusicUrlRes.Data.Url)
		if err != nil {
			t.Error(err)
			return
		}

		filename := path.Join("C:\\Users\\xlano\\projects\\side_project\\migu_music_downloader_wails\\app\\kuwo", musicId+".mp3")
		err = os.WriteFile(filename, getMusicRes.Body(), 0666)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestDownload1(t *testing.T) {
	musicId := "440616"
	urlBase := fmt.Sprintf("http://antiserver.kuwo.cn/anti.s?type=convert_url&rid=%s&format=mp3|acc&response=url", musicId)

	cli := resty.New().R()
	cli.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := cli.Get(urlBase)
	if err != nil {
		t.Error(err)
		return
	}

	tmp := string(resp.Body())
	//var getMusicUrlRes GetMusicUrlRes
	//json.Unmarshal(tmp, &getMusicUrlRes)
	fmt.Println(string(tmp))
	if tmp != "" && strings.HasPrefix(tmp, "http") {

		getMusicRes, err := cli.Get(tmp)
		if err != nil {
			t.Error(err)
			return
		}

		filename := path.Join("C:\\Users\\xlano\\projects\\side_project\\migu_music_downloader_wails\\app\\kuwo", musicId+"_1"+".mp3")
		err = os.WriteFile(filename, getMusicRes.Body(), 0666)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

type SearchRes struct {
	ARTISTPIC     string `json:"ARTISTPIC"`
	HIT           string `json:"HIT"`
	HITMODE       string `json:"HITMODE"`
	HITBUTOFFLINE string `json:"HIT_BUT_OFFLINE"`
	MSHOW         string `json:"MSHOW"`
	NEW           string `json:"NEW"`
	PN            string `json:"PN"`
	RN            string `json:"RN"`
	SHOW          string `json:"SHOW"`
	TOTAL         string `json:"TOTAL"`
	UK            string `json:"UK"`
	Abslist       []struct {
		AARTIST          string        `json:"AARTIST"`
		ALBUM            string        `json:"ALBUM"`
		ALBUMID          string        `json:"ALBUMID"`
		ALIAS            string        `json:"ALIAS"`
		ARTIST           string        `json:"ARTIST"`
		ARTISTID         string        `json:"ARTISTID"`
		CanSetRing       string        `json:"CanSetRing"`
		CanSetRingback   string        `json:"CanSetRingback"`
		DCTARGETID       string        `json:"DC_TARGETID"`
		DCTARGETTYPE     string        `json:"DC_TARGETTYPE"`
		DURATION         string        `json:"DURATION"`
		FARTIST          string        `json:"FARTIST"`
		FORMAT           string        `json:"FORMAT"`
		FSONGNAME        string        `json:"FSONGNAME"`
		KMARK            string        `json:"KMARK"`
		MINFO            string        `json:"MINFO"`
		MUSICRID         string        `json:"MUSICRID"`
		MVFLAG           string        `json:"MVFLAG"`
		MVPIC            string        `json:"MVPIC"`
		MVQUALITY        string        `json:"MVQUALITY"`
		NAME             string        `json:"NAME"`
		NEW              string        `json:"NEW"`
		NMINFO           string        `json:"N_MINFO"`
		ONLINE           string        `json:"ONLINE"`
		PAY              string        `json:"PAY"`
		PROVIDER         string        `json:"PROVIDER"`
		SONGNAME         string        `json:"SONGNAME"`
		SUBLIST          []interface{} `json:"SUBLIST"`
		SUBTITLE         string        `json:"SUBTITLE"`
		TAG              string        `json:"TAG"`
		AdSubtype        string        `json:"ad_subtype"`
		AdType           string        `json:"ad_type"`
		Allartistid      string        `json:"allartistid"`
		Audiobookpayinfo struct {
			Download string `json:"download"`
			Play     string `json:"play"`
		} `json:"audiobookpayinfo"`
		Barrage     string `json:"barrage"`
		CacheStatus string `json:"cache_status"`
		ContentType string `json:"content_type"`
		Fpay        string `json:"fpay"`
		HtsMVPIC    string `json:"hts_MVPIC"`
		Info        string `json:"info"`
		IotInfo     string `json:"iot_info"`
		Isdownload  string `json:"isdownload"`
		Isshowtype  string `json:"isshowtype"`
		Isstar      string `json:"isstar"`
		Mvpayinfo   struct {
			Download string `json:"download"`
			Play     string `json:"play"`
			Vid      string `json:"vid"`
		} `json:"mvpayinfo"`
		Nationid          string `json:"nationid"`
		Opay              string `json:"opay"`
		Originalsongtype  string `json:"originalsongtype"`
		OverseasCopyright string `json:"overseas_copyright"`
		OverseasPay       string `json:"overseas_pay"`
		PayInfo           struct {
			CannotDownload   string `json:"cannotDownload"`
			CannotOnlinePlay string `json:"cannotOnlinePlay"`
			Download         string `json:"download"`
			FeeType          struct {
				Album   string `json:"album"`
				Bookvip string `json:"bookvip"`
				Song    string `json:"song"`
				Vip     string `json:"vip"`
			} `json:"feeType"`
			Limitfree      string `json:"limitfree"`
			ListenFragment string `json:"listen_fragment"`
			LocalEncrypt   string `json:"local_encrypt"`
			Ndown          string `json:"ndown"`
			Nplay          string `json:"nplay"`
			OverseasNdown  string `json:"overseas_ndown"`
			OverseasNplay  string `json:"overseas_nplay"`
			Play           string `json:"play"`
			RefrainEnd     string `json:"refrain_end"`
			RefrainStart   string `json:"refrain_start"`
			TipsIntercept  string `json:"tips_intercept"`
		} `json:"payInfo"`
		ReactType         string `json:"react_type"`
		SpPrivilege       string `json:"spPrivilege"`
		SubsStrategy      string `json:"subsStrategy"`
		SubsText          string `json:"subsText"`
		Terminal          string `json:"terminal"`
		TmeMusicianAdtype string `json:"tme_musician_adtype"`
		Tpay              string `json:"tpay"`
		WebAlbumpicShort  string `json:"web_albumpic_short"`
		WebArtistpicShort string `json:"web_artistpic_short"`
		WebTimingonline   string `json:"web_timingonline"`
	} `json:"abslist"`
	Searchgroup string `json:"searchgroup"`
}

func TestSearch(t *testing.T) {
	//http://img3.kuwo.cn/star/albumcover/140/45/77/2616358704.jpg
	keyword := "jay"
	url := fmt.Sprintf("http://search.kuwo.cn/r.s?client=kt&all=%s&pn=0&rn=1&uid=221260053&ver=kwplayer_ar_99.99.99.99&vipver=1&ft=music&cluster=0&strategy=2012&encoding=utf8&rformat=json&vermerge=1&mobi=1", keyword)

	//url := fmt.Sprintf("http://www.kuwo.cn/api/www/search/searchMusicBykeyWord?key=%s&pn=1&rn=30&httpsStatus=1&reqId=80b33650-8a62-11ed-a069-8d99eba73f2a", keyword)

	cli := resty.New().R()
	//cli.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	//cli.SetHeader("Referer", fmt.Sprintf("http://www.kuwo.cn/search/list?key=%s", keyword))

	resp, err := cli.Get(url)
	if err != nil {
		t.Error(err)
		return
	}

	tmp := string(resp.Body())

	//tmp = strings.ReplaceAll(tmp, `\'`, "")

	var data SearchRes
	err = json.Unmarshal([]byte(tmp), &data)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(data)
}

//func TestDownloadLrc(t *testing.T) {
//	musicId := ""
//	url := fmt.Sprintf("http://m.kuwo.cn/newh5/singles/songinfoandlrc?musicId=%s", musicId)
//}

//func TestDownloadPic(t *testing.T) {
//	url := "http://artistpicserver.kuwo.cn/pic.web?corp=kuwo&type=rid_pic&pictype=url&content=list&size=640&rid=440616"
//}

func TestConvertLrc(t *testing.T) {

	url := "http://m.kuwo.cn/newh5/singles/songinfoandlrc?musicId=93151"
	cli := resty.New().R()
	resp, err := cli.Get(url)
	if err != nil {
		t.Error(err)
		return
	}

	tmp := string(resp.Body())

	var data model.KuwoLrcRes
	json.Unmarshal([]byte(tmp), &data)

	var sb strings.Builder
	for _, item := range data.Data.Lrclist {
		time := fmt.Sprintf("[%s]", convertTime(item.Time))
		sb.WriteString(time)
		sb.WriteString(item.LineLyric + "\n")
	}

	fmt.Println(sb.String())
}

func TestPic(t *testing.T) {
	url := "http://artistpicserver.kuwo.cn/pic.web?corp=kuwo&type=rid_pic&pictype=url&content=list&size=640&rid=156522"
	mp3Filename := "C:\\Users\\xlano\\projects\\side_project\\migu_music_downloader_wails\\frontend\\Justin Bieber - What do you mean ? (Remix).mp3"
	mp3Filename = util.FixWindowsFileName2Normal(mp3Filename)

	core := &AppKuwoCore{}
	core.ProcessPic(url, mp3Filename)
}
