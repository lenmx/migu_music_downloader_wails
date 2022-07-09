package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Downloader struct {
	queue chan DownloadQueueItem
}

var downloader *Downloader

func NewDownloader(count int) *Downloader {
	if downloader != nil {
		return downloader
	}

	return &Downloader{queue: make(chan DownloadQueueItem, count)}
}

func (d *Downloader) Push(item DownloadQueueItem) {
	d.queue <- item
}

func (d *Downloader) Run() {
	for {
		select {
		case item, ok := <-d.queue:
			if !ok {
				return
			}

			d.Download(item)
		default:
		}
	}
}

func (d *Downloader) Download(data DownloadQueueItem) {
	resp := BaseResponse{Code: 0, Message: "下载成功", Data: data}

	_, err := resty.New().R().SetOutput(data.Path).Get(data.Url)
	if err != nil {
		resp.Code = -1
		resp.Message = "下载失败：" + err.Error()
	}

	respJson, _ := json.Marshal(resp)
	runtime.EventsEmit(data.ctx, "download_result", string(respJson))
}
