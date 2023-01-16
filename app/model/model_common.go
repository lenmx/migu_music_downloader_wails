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
