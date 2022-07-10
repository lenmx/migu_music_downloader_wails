package util

import (
	"context"
	"encoding/json"
	"migu_music_downloader_wails/app/consts"
	"migu_music_downloader_wails/app/model"
)

func Log(ctx context.Context, message string) {
	Emit(ctx, consts.EventType_Log, message)
}

func PushDownloadResult(ctx context.Context, response model.BaseResponse) {
	respJson, _ := json.Marshal(response)
	Emit(ctx, consts.EventType_DownloadResult, string(respJson))
}
