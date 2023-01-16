package model

import "strconv"

type KuwoSearchRes struct {
	PN      string `json:"PN"`
	RN      string `json:"RN"`
	TOTAL   string `json:"TOTAL"`
	Abslist []struct {
		AARTIST        string `json:"AARTIST"`
		ALBUM          string `json:"ALBUM"`
		ALBUMID        string `json:"ALBUMID"`
		ALIAS          string `json:"ALIAS"`
		ARTIST         string `json:"ARTIST"`
		ARTISTID       string `json:"ARTISTID"`
		CanSetRing     string `json:"CanSetRing"`
		CanSetRingback string `json:"CanSetRingback"`
		DCTARGETID     string `json:"DC_TARGETID"`
		DCTARGETTYPE   string `json:"DC_TARGETTYPE"`
		DURATION       string `json:"DURATION"`
		FARTIST        string `json:"FARTIST"`
		FORMAT         string `json:"FORMAT"`
		FSONGNAME      string `json:"FSONGNAME"`
		KMARK          string `json:"KMARK"`
		MINFO          string `json:"MINFO"`
		MUSICRID       string `json:"MUSICRID"`
		MVFLAG         string `json:"MVFLAG"`
		MVPIC          string `json:"MVPIC"`
		MVQUALITY      string `json:"MVQUALITY"`
		NAME           string `json:"NAME"`
		NEW            string `json:"NEW"`
		NMINFO         string `json:"N_MINFO"`
		ONLINE         string `json:"ONLINE"`
		PAY            string `json:"PAY"`
		PROVIDER       string `json:"PROVIDER"`
		SONGNAME       string `json:"SONGNAME"`
		//SUBLIST          []interface{} `json:"SUBLIST"`
		//SUBTITLE         string        `json:"SUBTITLE"`
		//TAG              string        `json:"TAG"`
		//AdSubtype        string        `json:"ad_subtype"`
		//AdType           string        `json:"ad_type"`
		//Allartistid      string        `json:"allartistid"`
		//Audiobookpayinfo struct {
		//	Download string `json:"download"`
		//	Play     string `json:"play"`
		//} `json:"audiobookpayinfo"`
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
}

func (k *KuwoSearchRes) GetTotal() int {
	if k.TOTAL == "" {
		return 0
	}

	total, err := strconv.Atoi(k.TOTAL)
	if err != nil {
		return 0
	}

	return total
}

type KuwoDownloadItem struct {
	MusicId   string `json:"musicId"`
	MusicName string `json:"MusicName"`
	//Filename    string `json:"filename"`
	//DownloadUrl string `json:"downloadUrl"`
	//LrcUrl      string `json:"lrcUrl"`
	//PicUrl      string `json:"picUrl"`
}

type KuwoGetMusicUrlRes struct {
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
