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
	Start    int64  `json:"start,omitempty"`
	End      int64  `json:"end,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type MetaDataByte struct {
	Filename string `json:"filename,omitempty"`
	Size     int64  `json:"size,omitempty"`
}

type MetaData struct {
	Type uint8 `json:"type,omitempty"`
	Time int64 `json:"time,omitempty"`
	Data any   `json:"data,omitempty"`
}

func AddMetaData(d MetaData) {
	META_FILE_LOCK.L.Lock()
	defer META_FILE_LOCK.L.Unlock()
	d.Time = time.Now().UnixMilli()
	MetaDataMap = append(MetaDataMap, d)
}
