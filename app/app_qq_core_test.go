package app

import (
	"os"
	"testing"
)

func TestGetMusicUrl(t *testing.T) {
	core, err := NewAppQQCore()
	if err != nil {
		t.Error(err)
		return
	}

	filename := "M50000400jk23JDWwJ.mp3"
	tmp, err := core.getMusicUrl(filename)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(tmp)
}

func TestGetLrc(t *testing.T) {
	core, err := NewAppQQCore()
	if err != nil {
		t.Error(err)
		return
	}

	mid := "003cI52o4daJJL"
	path := "e:/tmp.lrc"
	content, err := core.getLyric(mid)
	if err != nil {
		t.Error(err)
		return
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err == nil {
		defer file.Close()
		file.WriteString(content)
	}
}

func TestSearch(t *testing.T) {
	app := &AppQQ{}

	res := app.OnSearch("jay", 1, 10)
	t.Log(res)
}
