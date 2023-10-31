package storage

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ogios/transfer-server/addon/udps"
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
	MetaDataMap     []*MetaData            = make([]*MetaData, 0)
	MetaDataDelList []int                  = make([]int, 0)
	MetaDataIDMap   map[string]int         = make(map[string]int)
	MetaDataTextMap map[string][]*MetaData = make(map[string][]*MetaData)
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

func ReloadMetaIndex() {
	MetaDataIDMap = make(map[string]int)
	MetaDataTextMap = make(map[string][]*MetaData)
	for index, metadata := range MetaDataMap {
		MetaDataIDMap[metadata.ID] = index
		if metadata.Type == TYPE_TEXT {
			if _, ok := MetaDataTextMap[metadata.Data.(*MetaDataText).Filename]; !ok {
				MetaDataTextMap[metadata.Data.(*MetaDataText).Filename] = []*MetaData{metadata}
			} else {
				MetaDataTextMap[metadata.Data.(*MetaDataText).Filename] = append(MetaDataTextMap[metadata.Data.(*MetaDataText).Filename], metadata)
			}
		}
	}
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
	if d.Type == TYPE_TEXT {
		if _, ok := MetaDataTextMap[d.Data.(*MetaDataText).Filename]; !ok {
			MetaDataTextMap[d.Data.(*MetaDataText).Filename] = []*MetaData{d}
		} else {
			MetaDataTextMap[d.Data.(*MetaDataText).Filename] = append(MetaDataTextMap[d.Data.(*MetaDataText).Filename], d)
		}
	}
	syncMeta()
	Notify()
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
	syncMeta()
	Notify()
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
		startoff = len(MetaDataMap)
		temp     = make([]*MetaData, 0)
		// for text filename
		fs = map[*MetaData]struct{}{}
	)
	for i := 0; i < len(MetaDataDelList); i++ {
		delindex := MetaDataDelList[i]
		// original startoff index
		if i == 0 {
			startoff = delindex
		}
		// if is the last deletion
		if i == len(MetaDataDelList)-1 {
			temp = append(temp, MetaDataMap[delindex+1:]...)
		} else {
			temp = append(temp, MetaDataMap[delindex+1:MetaDataDelList[i+1]]...)
		}

		// text filename add for index: MetaDataTextMap deletion
		m := MetaDataMap[delindex]
		if m.Type == TYPE_TEXT {
			fs[m] = struct{}{}
		}
	}
	MetaDataMap = append(MetaDataMap[:startoff], temp...)
	MetaDataDelList = make([]int, 0)
	ReloadMetaIndex()
	syncMeta()
}

func init() {
	TEXT_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
	BYTE_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
}

func Notify() {
	fmt.Println(MetaDataMap)
	udps.GlobalUdps.BoardCast([]byte{2})
}
