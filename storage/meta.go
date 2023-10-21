package storage

import (
	"math/rand"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	META_FILE_LOCK sync.Cond
	TEXT_FILE_LOCK sync.Cond
	BYTE_FILE_LOCK sync.Cond
)

var (
	TYPE_TEXT uint8 = 1
	TYPE_BYTE uint8 = 2
)

var ID_LENGTH = 5

var (
	// MetaDataMap   []MetaData
	MetaDataMap     []*MetaData    = make([]*MetaData, 0)
	MetaDataIDMap   map[string]int = make(map[string]int)
	MetaDataDelList []int          = make([]int, 0)
)

type IDPH struct{}

type MetaDataText struct {
	Filename string `json:"filename"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
}

type MetaDataByte struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type MetaData struct {
	Data any    `json:"data"`
	ID   string `json:"id"`
	Time int64  `json:"time"`
	Type uint8  `json:"type"`
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
		for i := 0; i < ID_LENGTH; i++ {
			id += string(RAND_FIELD[rand.Intn(len(RAND_FIELD))])
		}
		if _, ok := MetaDataIDMap[id]; !ok {
			return id
		}
	}
}

func AddMetaData(d *MetaData) {
	META_FILE_LOCK.L.Lock()
	defer META_FILE_LOCK.L.Unlock()
	d.Time = time.Now().UnixMilli()
	d.ID = getNewID()
	MetaDataMap = append(MetaDataMap, d)
	MetaDataIDMap[d.ID] = len(MetaDataMap) - 1
}

func DeleteMetaData(id string) {
	META_FILE_LOCK.L.Lock()
	defer META_FILE_LOCK.L.Unlock()
	if index, ok := MetaDataIDMap[id]; ok {
		MetaDataDelList = append(MetaDataDelList, index)
		sort.Slice(MetaDataDelList, func(i, j int) bool {
			return MetaDataDelList[i] < MetaDataDelList[j]
		})
		delete(MetaDataIDMap, id)
	}
}

func ClearDeleteMetaData() {
	META_FILE_LOCK.L.Lock()
	defer func() {
		runtime.GC()
	}()
	defer META_FILE_LOCK.L.Unlock()
	if len(MetaDataDelList) == 0 {
		return
	}
	var (
		startoff int
		delindex int
		temp     = make([]*MetaData, 0)
	)
	for i := 0; i < len(MetaDataDelList); i++ {
		delindex = MetaDataDelList[i]
		// original startoff index
		if i == 0 {
			startoff = delindex
		}
		// if is the last one
		if delindex == len(MetaDataMap)-1 {
			continue
		} else {
			// if is the last deletion
			if i == len(MetaDataDelList)-1 {
				temp = append(temp, MetaDataMap[delindex+1:]...)
			} else {
				temp = append(temp, MetaDataMap[delindex+1:MetaDataDelList[i+1]]...)
			}
		}
	}
	MetaDataMap = append(MetaDataMap[:startoff], temp...)
	MetaDataDelList = make([]int, 0)
}

func init() {
	TEXT_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
	BYTE_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
}
