package util

import (
	"context"
	"migu_music_downloader_wails/app/consts"
	"migu_music_downloader_wails/app/model"
	"time"
)

type Downloader struct {
	queue    chan model.DownloadQueueItem
	onResult func(response model.BaseResponse)
}

var downloader *Downloader

func NewDownloader(onResultCallback func(response model.BaseResponse), count int) *Downloader {
	if downloader != nil {
		return downloader
	}

	return &Downloader{
		queue:    make(chan model.DownloadQueueItem, count),
		onResult: onResultCallback,
	}
}

func (d *Downloader) Push(ctx context.Context, item model.DownloadQueueItem) {
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			d.queue <- item
		}
	}()
}

func (d *Downloader) Run() {
	exit := false

	for {
		select {
		case item, ok := <-d.queue:
			if !ok {
				d.queue = make(chan model.DownloadQueueItem)
				exit = true
			} else {
				go d.download(item)
				time.Sleep(time.Second * 1)
			}
		default:
			time.Sleep(consts.LoopInterval)
		}

		if exit {
			break
		}
	}
}

func (d *Downloader) Stop() {
	close(d.queue)
}

func (d *Downloader) download(data model.DownloadQueueItem) {
	resp := model.BaseResponse{Code: 0, Message: "下载成功", Data: data}

	if len(data.Url) > 0 && data.Mp3Process != nil {
		err := data.Mp3Process(data.Url, data.Filename)
		if err != nil {
			resp.Code = -1
			resp.Message = "下载失败: " + err.Error()
			if d.onResult != nil {
				d.onResult(resp)
			}
			return
		}
	}

	if len(data.LrcUrl) > 0 && data.LrcProcess != nil {
		data.LrcProcess(data.LrcUrl, data.Filename)
	}

	if len(data.PicUrl) > 0 && data.PicProcess != nil {
		data.PicProcess(data.PicUrl, data.Filename)
	}

	if d.onResult != nil {
		d.onResult(resp)
	}
}
