package consts

type SourceType string

const (
	SourceType_SQ SourceType = "SQ&formatType=SQ&resourceType=E"
	SourceType_HQ SourceType = "HQ&formatType=HQ&resourceType=2"
)

var (
	SourceType2FileExt = map[SourceType]string{
		SourceType_SQ: ".flac",
		SourceType_HQ: ".mp3",
	}
)

type EventType string

const (
	EventType_Log            EventType = "log"
	EventType_DownloadResult EventType = "download_result"
)

const (
	SearchSuccess = "000000"
)
