package util

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"migu_music_downloader_wails/app/consts"
)

func Emit(ctx context.Context, eventType consts.EventType, message string)  {
	runtime.EventsEmit(ctx, string(eventType), message)
}