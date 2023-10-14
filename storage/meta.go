package storage

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

var META_FILE_LOCK sync.Cond

var TYPE_TEXT uint8 = 1
var TYPE_BYTE uint8 = 2

var MetaDataMap []MetaData
var MetaDataIDMap map[string]IDPH

type IDPH struct{}

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
	Type uint8  `json:"type"`
	Time int64  `json:"time"`
	ID   string `json:"id"`
	Data any    `json:"data"`
}

var RAND_FIELD = func() string {
	base := "0123456789"
	alphab := ""
	for i := 0; i < 26; i++ {
		alphab += string(rune(97 + i))
	}
	return base + alphab + strings.ToUpper(alphab)
}()

func getNewID() string {
	for {
		id := ""
		for i := 0; i < 5; i++ {
			id += string(RAND_FIELD[rand.Intn(len(RAND_FIELD))])
		}
		if _, ok := MetaDataIDMap[id]; !ok {
			return id
		}
	}
}

func AddMetaData(d MetaData) {
	META_FILE_LOCK.L.Lock()
	defer META_FILE_LOCK.L.Unlock()
	d.Time = time.Now().UnixMilli()
	d.ID = getNewID()
	MetaDataMap = append(MetaDataMap, d)
	MetaDataIDMap[d.ID] = IDPH{}
}
