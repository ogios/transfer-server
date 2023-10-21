package fetch

import (
	"fmt"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/storage"
)

func fetchFromData(start, size int) ([]*storage.MetaData, int) {
	storage.META_FILE_LOCK.L.Lock()
	defer storage.META_FILE_LOCK.L.Unlock()
	b := make([]*storage.MetaData, size)
	offset := 0
	startoff := start
	if len(storage.MetaDataDelList) == 0 {
		copy(b, storage.MetaDataMap[startoff:startoff+size])
		return b, len(storage.MetaDataMap)
	}
	for index, deleted := range storage.MetaDataDelList {
		if deleted <= start {
			startoff++
		} else {
			if offset >= len(b) || startoff >= len(storage.MetaDataMap) {
				break
			}
			length := deleted - startoff
			left := len(b) - offset
			if length > left {
				copy(b[offset:], storage.MetaDataMap[startoff:startoff+left])
				break
			} else {
				copy(b[offset:offset+length], storage.MetaDataMap[startoff:startoff+length])
				startoff = deleted + 1
				offset += length
			}
			if index == len(storage.MetaDataDelList)-1 && offset != len(b) && startoff < len(storage.MetaDataMap) {
				copy(b[offset:], storage.MetaDataMap[startoff:])
			}
		}
	}
	return b[:offset], len(storage.MetaDataMap) - len(storage.MetaDataDelList)
}

func Fetch(pageindex int, size int) ([]*storage.MetaData, int, error) {
	index := pageindex * size
	log.Debug(nil, "fetch index: %d", index)
	data, total := fetchFromData(index, size)
	log.Debug(nil, "fetch copied data: %v", data)

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
	return data, total, nil
}
