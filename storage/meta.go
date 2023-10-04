package storage

import "sync"

var META_FILE_LOCK sync.Cond

var TYPE_TEXT uint8 = 1
var TYPE_BYTE uint8 = 2

var MetaDataMap []MetaData

type MetaDataText struct {
	Start int64
	End   int64
}

type MetaDataByte struct {
	Filename string
	Size     int64
}

type MetaData struct {
	Type uint8
	Data any
}

func AddMetaData(d MetaData) {
	META_FILE_LOCK.L.Lock()
	defer META_FILE_LOCK.L.Unlock()
	MetaDataMap = append(MetaDataMap, d)
}
