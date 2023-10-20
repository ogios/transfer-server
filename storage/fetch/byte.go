package fetch

import (
	"fmt"
	"os"

	"github.com/ogios/transfer-server/storage"
)

func searchFileFromMeta(id string) *storage.MetaData {
	storage.META_FILE_LOCK.L.Lock()
	defer storage.META_FILE_LOCK.L.Unlock()
	index, ok := storage.MetaDataIDMap[id]
	if ok {
		return storage.MetaDataMap[index]
	}
	// for index, data := range storage.MetaDataMap {
	// 	if data.Type == storage.TYPE_BYTE {
	// 		d := (data.Data).(*storage.MetaDataByte)
	// 		if d.Filename == name {
	// 			return storage.MetaDataMap[index]
	// 		}
	// 	}
	// }
	return nil
}

func FetchByte(id string) (*os.File, int64, error) {
	meta := searchFileFromMeta(id)
	if meta == nil {
		return nil, 0, fmt.Errorf("file not found in metadata")
	}
	f, err := os.OpenFile(
		storage.BASE_PATH_BYTE+"/"+meta.Data.(storage.MetaDataByte).Filename,
		os.O_RDONLY,
		0644,
	)
	if err != nil {
		return nil, 0, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}
	return f, stat.Size(), nil
}
