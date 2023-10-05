package save

import (
	"os"
	"strconv"
	"sync"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/storage"
)

var TEXT_FILE_MAX_SIZE int64 = 16 * 1024
var TEXT_FILE_LOCK sync.Cond

func init() {
	TEXT_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
}

func getTextFile() (*os.File, error) {
	log.Debug(nil, "Reading text dir")
	files, err := os.ReadDir(storage.BASE_PATH_TEXT)
	if err != nil {
		log.Error(nil, "Read text dir error: %s", err)
		return nil, err
	}
	m := 1
	var ff os.DirEntry = nil
	log.Debug(nil, "Getting max text file start")
	for _, f := range files {
		ind, err := strconv.Atoi(f.Name())
		if err == nil {
			if ind >= m {
				log.Debug(nil, "max text file replace: &d", ind)
				m = ind
				ff = f
			}
		}
	}
	log.Debug(nil, "Getting max text file done")
	if ff == nil {
		log.Info(nil, "dir empty, create 1")
		return os.Create(storage.BASE_PATH_TEXT + "/" + strconv.Itoa(m))
	} else {
		log.Debug(nil, "checking if text file size over: %d", TEXT_FILE_MAX_SIZE)
		info, err := ff.Info()
		if err != nil {
			log.Error(nil, "Read text file info error: %s", err)
			return nil, err
		}
		log.Debug(nil, "text file size: %d", info.Size())
		if info.Size() >= TEXT_FILE_MAX_SIZE {
			log.Debug(nil, "creating new text file")
			return os.Create(storage.BASE_PATH_TEXT + "/" + strconv.Itoa(m+1))
		} else {
			log.Debug(nil, "using old text file")
			return os.OpenFile(storage.BASE_PATH_TEXT+"/"+strconv.Itoa(m), os.O_RDWR, 0644)
		}
	}
}

func saveText(reader *normal.Conn) (string, int64, int64, error) {
	TEXT_FILE_LOCK.L.Lock()
	defer TEXT_FILE_LOCK.L.Unlock()
	log.Debug(nil, "getting text file...")
	f, err := getTextFile()
	if err != nil {
		log.Error(nil, "get text file error: %s", err)
		return "", 0, 0, err
	}
	defer f.Close()
	log.Debug(nil, "saving text file...")
	start, end, err := save(reader, f)
	if err != nil {
		log.Error(nil, "save text file error: %s", err)
		return "", 0, 0, err
	}
	return f.Name(), start, end, err
}

func SaveText(reader *normal.Conn) error {
	log.Debug(nil, "saving text...")
	filename, start, end, err := saveText(reader)
	if err != nil {
		log.Error(nil, "save text error: %s", err)
		return err
	}
	log.Debug(nil, "save text done: start-%d end-%d", start, end)
	log.Debug(nil, "saving text metadata...")
	storage.AddMetaData(storage.MetaData{
		Type: storage.TYPE_TEXT,
		Data: storage.MetaDataText{
			Start:    start,
			End:      end,
			Filename: filename,
		},
	})
	log.Debug(nil, "save text metadata done")
	return nil
}
