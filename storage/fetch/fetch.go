package fetch

import (
	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/storage"
)

func fetchFromData(start int, end int) []storage.MetaData {
	storage.META_FILE_LOCK.L.Lock()
	defer storage.META_FILE_LOCK.L.Unlock()
	if end == -1 {
		return storage.MetaDataMap[start:]
	}
	return storage.MetaDataMap[start:end]
}

func Fetch(pageindex int, size int) ([]storage.MetaData, error) {
	index := pageindex * size
	log.Debug(nil, "fetch index: %d", index)
	if len(storage.MetaDataMap)-1 < index {
		log.Debug(nil, "fetch index surpass length - %d", len(storage.MetaDataMap))
		return make([]storage.MetaData, 0), nil
	} else {
		start := len(storage.MetaDataMap) - 1 - index - size
		log.Debug(nil, "fetch start: %d", start)
		var d []storage.MetaData
		if start < 0 {
			d = fetchFromData(0, -1)
		} else {
			d = fetchFromData(start, start+size)
		}
		data := make([]storage.MetaData, len(d))
		log.Debug(nil, "fetch d: %v", d)
		copy(data, d)
		log.Debug(nil, "fetch copied data: %v", data)
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
		log.Debug(nil, "fetch converted data: %v", data)

		for index, m := range data {
			if m.Type == storage.TYPE_TEXT {
				data[index].Data = func() (text any) {
					defer func() {
						if err := recover(); err != nil {
							text = -1
						}
					}()
					log.Debug(nil, "fetching text content")
					temp := m.Data
					log.Debug(nil, "fetch temp: %v", m)
					d := temp.(*storage.MetaDataText)
					text, err := FetchText(d.Start, d.End, d.Filename)
					if err != nil {
						log.Error(nil, "Fetch text content error")
						return 0
					}
					log.Debug(nil, "fetch text content: %s", text)
					return text
				}()
			}
		}
		return data, nil
	}
}
