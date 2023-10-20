package fetch

import (
	"fmt"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/storage"
)

func fetchFromData(start, size int) []*storage.MetaData {
	storage.META_FILE_LOCK.L.Lock()
	defer storage.META_FILE_LOCK.L.Unlock()
	add, deleted := 0, 0
	for _, deleted = range storage.MetaDataDelList {
		if deleted >= start {
			add++
		} else {
			break
		}
	}
	startoff := start + add
	total := len(storage.MetaDataMap)
	if startoff > total {
		return make([]*storage.MetaData, 0)
	}
	var b []*storage.MetaData
	if total >= startoff+size {
		b = make([]*storage.MetaData, size)
	} else {
		b = make([]*storage.MetaData, total-startoff)
	}

	for i := 0; i < len(b); i++ {
		b[i] = storage.MetaDataMap[startoff+i]
	}
	return b
}

// func fetchFromData(start int, end int) []*storage.MetaData {
// 	storage.META_FILE_LOCK.L.Lock()
// 	defer storage.META_FILE_LOCK.L.Unlock()
// 	var d []*storage.MetaData
// 	if end == -1 {
// 		d = make([]*storage.MetaData, len(storage.MetaDataIDMap)-start)
// 	} else {
// 		d = make([]*storage.MetaData, end-start)
// 	}
// 	for i := 0; i < len(d); i++ {
// 	}
// 	return storage.MetaDataMap[start:end]
// }

func Fetch(pageindex int, size int) ([]*storage.MetaData, error) {
	index := pageindex * size
	log.Debug(nil, "fetch index: %d", index)
	if len(storage.MetaDataMap)-1 < index {
		log.Debug(nil, "fetch index surpass length - %d", len(storage.MetaDataMap))
		return make([]*storage.MetaData, 0), nil
	} else {
		// start := len(storage.MetaDataMap) - 1 - index - size
		// log.Debug(nil, "fetch start: %d", start)
		var data []*storage.MetaData = fetchFromData(index, size)
		// if start < 0 {
		// 	d = fetchFromData(0, -1)
		// } else {
		// 	d = fetchFromData(start, start+size)
		// }
		// data := make([]*storage.MetaData, len(d))
		// log.Debug(nil, "fetch d: %v", d)
		// copy(data, d)
		log.Debug(nil, "fetch copied data: %v", data)
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
		log.Debug(nil, "fetch converted data: %v", data)

		for index, m := range data {
			if m.Type == storage.TYPE_TEXT {
				data[index] = func() (meta *storage.MetaData) {
					defer func() {
						if err := recover(); err != nil {
							log.Error(nil, "%v", err)
							meta = &storage.MetaData{
								ID:   data[index].ID,
								Time: data[index].Time,
								Type: data[index].Type,
								Data: "-1",
							}
						}
					}()
					log.Debug(nil, "fetching text content")
					temp := m.Data
					log.Debug(nil, "fetch temp: %v - %v", m, m.Data)
					d := temp.(*storage.MetaDataText)
					text, err := FetchText(d.Start, d.End, d.Filename)
					if err != nil {
						msg := fmt.Sprintf("Fetch text content error: %v", err)
						panic(fmt.Errorf(msg))
					}
					log.Debug(nil, "fetch text content: %s", text)
					return &storage.MetaData{
						ID:   data[index].ID,
						Time: data[index].Time,
						Type: data[index].Type,
						Data: text,
					}
				}()
			}
		}
		return data, nil
	}
}
