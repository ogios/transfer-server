package fetch

import "github.com/ogios/transfer-server/storage"

func Fetch(pageindex int, size int) ([]storage.MetaData, error) {
	index := pageindex * size
	if len(storage.MetaDataMap)-1 < index {
		return make([]storage.MetaData, 0), nil
	} else {
		data := make([]storage.MetaData, size)
		start := len(storage.MetaDataMap) - 1 - index
		if start+size > len(storage.MetaDataMap) {
			copy(data, storage.MetaDataMap[start:])
		} else {
			copy(data, storage.MetaDataMap[start:start+size])
		}
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}

		for index, m := range data {
			if m.Type == storage.TYPE_TEXT {
				d := m.Data.(storage.MetaDataText)
				text, err := FetchText(d.Start, d.End, d.Filename)
				if err != nil {
					return nil, err
				}
				data[index].Data = text
			}
		}
		return data, nil
	}
}
