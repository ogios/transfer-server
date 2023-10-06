package storage

import (
	"sync"
	"time"
)

var META_FILE_LOCK sync.Cond

var TYPE_TEXT uint8 = 1
var TYPE_BYTE uint8 = 2

var MetaDataMap []MetaData

type MetaDataText struct {
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Filename string `json:"filename"`
}

type MetaDataByte struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type MetaData struct {
	Type uint8 `json:"type"`
	Time int64 `json:"time"`
	Data any   `json:"data"`
}

func AddMetaData(d MetaData) {
	META_FILE_LOCK.L.Lock()
	defer META_FILE_LOCK.L.Unlock()
	d.Time = time.Now().UnixMilli()
	MetaDataMap = append(MetaDataMap, d)
}
