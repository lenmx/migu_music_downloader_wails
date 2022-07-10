package util

import (
	"github.com/go-resty/resty/v2"
	"migu_music_downloader_wails/app/model"
	"time"
)

type Downloader struct {
	queue    chan model.DownloadQueueItem
	ResultCh chan model.BaseResponse
}

var downloader *Downloader

func NewDownloader(count int) *Downloader {
	if downloader != nil {
		return downloader
	}

	return &Downloader{
		queue:    make(chan model.DownloadQueueItem, count),
		ResultCh: make(chan model.BaseResponse, count),
	}
}

func (d *Downloader) Push(item model.DownloadQueueItem) {
	d.queue <- item
}

func (d *Downloader) Run() {
	for {
		select {
		case item, ok := <-d.queue:
			if !ok {
				return
			}

			go d.download(item)
			time.Sleep(time.Second * 1)
		default:
		}
	}
}

func (d *Downloader) download(data model.DownloadQueueItem) {
	resp := model.BaseResponse{Code: 0, Message: "下载成功", Data: data}

	_, err := resty.New().R().SetOutput(data.Path).Get(data.Url)
	if err != nil {
		resp.Code = -1
		resp.Message = "下载失败：" + err.Error()
	}

	d.ResultCh <- resp
}
