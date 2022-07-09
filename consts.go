package main

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

const (
	searchSuccess = "000000"
)
