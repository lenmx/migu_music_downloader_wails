package model

type SearchRes struct {
	Code           string `json:"code"`
	Info           string `json:"info"`
	SongResultData struct {
		TotalCount string        `json:"totalCount"`
		Correct    []interface{} `json:"correct"`
		Result     []struct {
			Id           string   `json:"id"`
			ResourceType string   `json:"resourceType"`
			ContentId    string   `json:"contentId"`
			CopyrightId  string   `json:"copyrightId"`
			Name         string   `json:"name"`
			HighlightStr []string `json:"highlightStr"`
			Singers      []struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"singers"`
			Albums []struct {
				Id   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"albums"`
			Tags     []string `json:"tags"`
			LyricUrl string   `json:"lyricUrl"`
			TrcUrl   string   `json:"trcUrl"`
			ImgItems []struct {
				ImgSizeType string `json:"imgSizeType"`
				Img         string `json:"img"`
			} `json:"imgItems"`
			RateFormats []struct {
				ResourceType         string `json:"resourceType"`
				FormatType           string `json:"formatType"`
				Url                  string `json:"url,omitempty"`
				Format               string `json:"format"`
				Size                 string `json:"size"`
				FileType             string `json:"fileType,omitempty"`
				Price                string `json:"price"`
				IosUrl               string `json:"iosUrl,omitempty"`
				AndroidUrl           string `json:"androidUrl,omitempty"`
				AndroidFileType      string `json:"androidFileType,omitempty"`
				IosFileType          string `json:"iosFileType,omitempty"`
				IosSize              string `json:"iosSize,omitempty"`
				AndroidSize          string `json:"androidSize,omitempty"`
				IosFormat            string `json:"iosFormat,omitempty"`
				AndroidFormat        string `json:"androidFormat,omitempty"`
				IosAccuracyLevel     string `json:"iosAccuracyLevel,omitempty"`
				AndroidAccuracyLevel string `json:"androidAccuracyLevel,omitempty"`
			} `json:"rateFormats"`
			NewRateFormats []struct {
				ResourceType         string `json:"resourceType"`
				FormatType           string `json:"formatType"`
				Url                  string `json:"url,omitempty"`
				Format               string `json:"format,omitempty"`
				Size                 string `json:"size,omitempty"`
				FileType             string `json:"fileType,omitempty"`
				Price                string `json:"price"`
				IosUrl               string `json:"iosUrl,omitempty"`
				AndroidUrl           string `json:"androidUrl,omitempty"`
				AndroidFileType      string `json:"androidFileType,omitempty"`
				IosFileType          string `json:"iosFileType,omitempty"`
				IosSize              string `json:"iosSize,omitempty"`
				AndroidSize          string `json:"androidSize,omitempty"`
				IosFormat            string `json:"iosFormat,omitempty"`
				AndroidFormat        string `json:"androidFormat,omitempty"`
				IosAccuracyLevel     string `json:"iosAccuracyLevel,omitempty"`
				AndroidAccuracyLevel string `json:"androidAccuracyLevel,omitempty"`
				AndroidNewFormat     string `json:"androidNewFormat,omitempty"`
				IosBit               int    `json:"iosBit,omitempty"`
				AndroidBit           int    `json:"androidBit,omitempty"`
			} `json:"newRateFormats"`
		} `json:"result"`
	} `json:"songResultData"`
}

type PageRes struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

type Setting struct {
	SavePath      string `json:"savePath",yaml:"savePath"`
	DownloadLrc   bool   `json:"downloadLrc",yaml:"downloadLrc"`
	DownloadCover bool   `json:"downloadCover",yaml:"downloadCover"`
}

type DownloadItem struct {
	ContentId string `json:"contentId"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	LrcUrl    string `json:"lrcUrl"`
	Cover     string `json:"cover"`
}

type DownloadQueueItem struct {
	DownloadItem
	Path string `json:"path"`

	DownloadLrc   bool
	DownloadCover bool
}

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
