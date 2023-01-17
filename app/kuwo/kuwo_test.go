package kuwo

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"migu_music_downloader_wails/app/model"
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
	tmp := `{"data":{"lrclist":[{"lineLyric":"搁浅 - 周杰伦","time":"0.0"},{"lineLyric":"词：宋健彰","time":"1.99"},{"lineLyric":"曲：周杰伦","time":"4.08"},{"lineLyric":"编曲：钟兴民","time":"7.84"},{"lineLyric":"久未放晴的天空","time":"16.61"},{"lineLyric":"依旧留着你的笑容","time":"20.15"},{"lineLyric":"哭过却无法掩埋歉疚","time":"24.4"},{"lineLyric":"风筝在阴天搁浅","time":"32.02"},{"lineLyric":"想念还在等待救援","time":"35.83"},{"lineLyric":"我拉着线复习你给的温柔","time":"40.35"},{"lineLyric":"曝晒在一旁的寂寞","time":"46.63"},{"lineLyric":"笑我给不起承诺","time":"50.66"},{"lineLyric":"怎么会怎么会你竟原谅了我","time":"53.98"},{"lineLyric":"我只能永远读着对白","time":"61.2"},{"lineLyric":"读着我给你的伤害","time":"65.5"},{"lineLyric":"我原谅不了我","time":"69.39"},{"lineLyric":"就请你当作我已不在","time":"72.17"},{"lineLyric":"我睁开双眼看着空白","time":"76.74"},{"lineLyric":"忘记你对我的期待","time":"81.17"},{"lineLyric":"读完了依赖","time":"85.21"},{"lineLyric":"我很快就离开","time":"87.76"},{"lineLyric":"久未放晴的天空","time":"114.21"},{"lineLyric":"依旧留着你的笑容","time":"117.94"},{"lineLyric":"哭过却无法掩埋歉疚","time":"122.59"},{"lineLyric":"风筝在阴天搁浅","time":"129.74"},{"lineLyric":"想念还在等待救援","time":"133.53"},{"lineLyric":"我拉着线复习你给的温柔","time":"138.03"},{"lineLyric":"曝晒在一旁的寂寞","time":"144.57"},{"lineLyric":"笑我给不起承诺","time":"148.6"},{"lineLyric":"怎么会怎么会你竟原谅了我","time":"151.6"},{"lineLyric":"我只能永远读着对白","time":"158.98"},{"lineLyric":"读着我给你的伤害","time":"163.19"},{"lineLyric":"我原谅不了我","time":"167.22"},{"lineLyric":"就请你当作我已不在","time":"170.04001"},{"lineLyric":"我睁开双眼看着空白","time":"174.53"},{"lineLyric":"忘记你对我的期待","time":"178.92"},{"lineLyric":"读完了依赖","time":"182.95"},{"lineLyric":"我很快就","time":"185.74"},{"lineLyric":"我只能永远读着对白","time":"188.97"},{"lineLyric":"读着我给你的伤害","time":"194.59"},{"lineLyric":"我原谅不了我","time":"198.49"},{"lineLyric":"就请你当作我已不在","time":"201.34"},{"lineLyric":"我睁开双眼看着空白","time":"206.1"},{"lineLyric":"忘记你对我的期待","time":"210.18"},{"lineLyric":"读完了依赖","time":"214.17"},{"lineLyric":"我很快就离开","time":"217.14"}],"simpl":{"musiclist":[{"album":"","albumId":null,"artist":"","artistId":null,"contentType":null,"coopFormats":[],"copyRight":null,"duration":null,"formats":null,"hasEcho":null,"hasMv":null,"id":"228908","isExt":null,"isNew":null,"isPoint":null,"isbatch":null,"isdownload":null,"isstar":"0","mkvNsig1":null,"mkvNsig2":null,"mkvRid":null,"mp3Nsig1":null,"mp3Nsig2":null,"mp3Rid":null,"mp3Size":"","mp4sig1":"","mp4sig2":"","musicrId":"228908","mutiVer":null,"mvpayinfo":null,"mvpic":null,"nsig1":null,"nsig2":null,"online":null,"params":null,"pay":null,"pic":"","playCnt":null,"rankChange":null,"reason":"猜你喜欢","score":"1","score100":null,"songName":"","songTimeMinutes":"","tpay":null,"trend":null,"upTime":null,"uploader":null},{"album":"","albumId":null,"artist":"","artistId":null,"contentType":null,"coopFormats":[],"copyRight":null,"duration":null,"formats":null,"hasEcho":null,"hasMv":null,"id":"228912","isExt":null,"isNew":null,"isPoint":null,"isbatch":null,"isdownload":null,"isstar":"0","mkvNsig1":null,"mkvNsig2":null,"mkvRid":null,"mp3Nsig1":null,"mp3Nsig2":null,"mp3Rid":null,"mp3Size":"","mp4sig1":"","mp4sig2":"","musicrId":"228912","mutiVer":null,"mvpayinfo":null,"mvpic":null,"nsig1":null,"nsig2":null,"online":null,"params":null,"pay":null,"pic":"","playCnt":null,"rankChange":null,"reason":"猜你喜欢","score":"1","score100":null,"songName":"","songTimeMinutes":"","tpay":null,"trend":null,"upTime":null,"uploader":null},{"album":"","albumId":null,"artist":"","artistId":null,"contentType":null,"coopFormats":[],"copyRight":null,"duration":null,"formats":null,"hasEcho":null,"hasMv":null,"id":"226543302","isExt":null,"isNew":null,"isPoint":null,"isbatch":null,"isdownload":null,"isstar":"0","mkvNsig1":null,"mkvNsig2":null,"mkvRid":null,"mp3Nsig1":null,"mp3Nsig2":null,"mp3Rid":null,"mp3Size":"","mp4sig1":"","mp4sig2":"","musicrId":"226543302","mutiVer":null,"mvpayinfo":null,"mvpic":null,"nsig1":null,"nsig2":null,"online":null,"params":null,"pay":null,"pic":"","playCnt":null,"rankChange":null,"reason":"猜你喜欢","score":"1","score100":null,"songName":"","songTimeMinutes":"","tpay":null,"trend":null,"upTime":null,"uploader":null},{"album":"","albumId":null,"artist":"","artistId":null,"contentType":null,"coopFormats":[],"copyRight":null,"duration":null,"formats":null,"hasEcho":null,"hasMv":null,"id":"118980","isExt":null,"isNew":null,"isPoint":null,"isbatch":null,"isdownload":null,"isstar":"0","mkvNsig1":null,"mkvNsig2":null,"mkvRid":null,"mp3Nsig1":null,"mp3Nsig2":null,"mp3Rid":null,"mp3Size":"","mp4sig1":"","mp4sig2":"","musicrId":"118980","mutiVer":null,"mvpayinfo":null,"mvpic":null,"nsig1":null,"nsig2":null,"online":null,"params":null,"pay":null,"pic":"","playCnt":null,"rankChange":null,"reason":"猜你喜欢","score":"1","score100":null,"songName":"","songTimeMinutes":"","tpay":null,"trend":null,"upTime":null,"uploader":null},{"album":"","albumId":null,"artist":"","artistId":null,"contentType":null,"coopFormats":[],"copyRight":null,"duration":null,"formats":null,"hasEcho":null,"hasMv":null,"id":"94237","isExt":null,"isNew":null,"isPoint":null,"isbatch":null,"isdownload":null,"isstar":"0","mkvNsig1":null,"mkvNsig2":null,"mkvRid":null,"mp3Nsig1":null,"mp3Nsig2":null,"mp3Rid":null,"mp3Size":"","mp4sig1":"","mp4sig2":"","musicrId":"94237","mutiVer":null,"mvpayinfo":null,"mvpic":null,"nsig1":null,"nsig2":null,"online":null,"params":null,"pay":null,"pic":"","playCnt":null,"rankChange":null,"reason":"猜你喜欢","score":"1","score100":null,"songName":"","songTimeMinutes":"","tpay":null,"trend":null,"upTime":null,"uploader":null}],"playlist":[{"digest":"8","disname":"他来了，说好不哭","extend":"0","info":"","isnew":"0","name":"他来了，说好不哭","newcount":"0","nodeid":"0","pic":"http://img1.kwcdn.kuwo.cn/star/userpl2015/47/81/1568615726642_458076447b.jpg","playcnt":"4630809","source":"8","sourceid":"2863447252","tag":"华语#流行#经典"},{"digest":"8","disname":"终于等到周杰伦，说好不哭你今天哭了吗？","extend":"0","info":"","isnew":"0","name":"终于等到周杰伦，说好不哭你今天哭了吗？","newcount":"0","nodeid":"0","pic":"http://img1.kwcdn.kuwo.cn/star/userpl2015/81/23/1568684821020_182253281b.jpg","playcnt":"4388986","source":"8","sourceid":"2867496601","tag":"华语#流行#90年代"}]},"songinfo":{"album":"七里香","albumId":"4533","artist":"周杰伦","artistId":"336","contentType":"0","coopFormats":["320kmp3","192kmp3","128kmp3"],"copyRight":"0","duration":"238","formats":"WMA96|WMA128|MP3H|MP3192|MP3128|ALFLAC|AL|AAC48|AAC24|EXMV500|EXMP4UL|EXMP4L|EXMP4HV|EXMP4BD|EXMP4","hasEcho":null,"hasMv":"1","id":"94239","isExt":null,"isNew":null,"isPoint":"0","isbatch":null,"isdownload":"0","isstar":"0","mkvNsig1":"3265342408","mkvNsig2":"4284284954","mkvRid":"MV_4012819","mp3Nsig1":"77326800","mp3Nsig2":"1631030247","mp3Rid":"MP3_94239","mp3Size":"","mp4sig1":"","mp4sig2":"","musicrId":"MUSIC_94239","mutiVer":"0","mvpayinfo":null,"mvpic":null,"nsig1":"2103892599","nsig2":"657931830","online":"1","params":null,"pay":"16711935","pic":"http://img1.kwcdn.kuwo.cn/star/albumcover/240/30/97/4276557883.jpg","playCnt":"","rankChange":null,"reason":null,"score":null,"score100":"82","songName":"搁浅","songTimeMinutes":"03:58","tpay":null,"trend":null,"upTime":"","uploader":""}},"msg":"成功","msgs":null,"profileid":"site","reqid":"854a6553Xba25X4107Xa10aXe58814a9d4c7","status":200}`

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
