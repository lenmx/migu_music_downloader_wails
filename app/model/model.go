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
		} `json:"result"`
	} `json:"songResultData"`
}

type PageRes struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

type DownloadItem struct {
	ContentId string `json:"contentId"`
	Name      string `json:"name"`
}

type DownloadQueueItem struct {
	DownloadItem
	Path string `json:"path"`
	Url  string `json:"url"`
}

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
