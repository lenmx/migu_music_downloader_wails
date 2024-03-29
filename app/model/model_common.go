package model

type SearchItemRes struct {
	MusicId string   `json:"musicId"`
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	Singers []string `json:"singers"`
	Albums  []string `json:"albums"`
	LrcUrl  string   `json:"lrcUrl"`
	Cover   string   `json:"cover"`
}

type LrcProcessFn func(url, filename string)

type PicProcessFn func(url, filename string)

type Mp3ProcessFn func(url, filename string) error
