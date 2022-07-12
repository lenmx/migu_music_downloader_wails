package util

import (
	"context"
	"github.com/go-resty/resty/v2"
	"migu_music_downloader_wails/app/consts"
	"migu_music_downloader_wails/app/model"
	"strings"
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

	_, err := resty.New().R().SetOutput(data.Path).Get(data.Url)
	if err != nil {
		resp.Code = -1
		resp.Message = "下载失败：" + err.Error()
	}

	idx := strings.LastIndex(data.Path, ".")
	if data.DownloadLrc && len(data.LrcUrl) > 0 {
		path := data.Path[:idx] + ".lrc"
		resty.New().R().SetOutput(path).Get(data.LrcUrl)
	}
	if data.DownloadCover && len(data.Cover) > 0 {
		path := data.Path[:idx] + ".png"
		resty.New().R().SetOutput(path).Get(data.Cover)
	}

	if d.onResult != nil {
		d.onResult(resp)
	}
}
