package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ogios/transfer-server/config"
	"github.com/ogios/transfer-server/log"
	"golang.org/x/exp/slog"
)

var (
	BASE_PATH      string
	BASE_PATH_TEXT string
	BASE_PATH_BYTE string
	BASE_PATH_META string
)

func makeDir(dir string) {
	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		log.Info(nil, "Creating dir: %s", dir)
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	p, _ := filepath.Abs(dir)
	log.Debug(nil, "exist dir: %s", p)
}

func makeFile(path string) {
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		f, err := os.Create(path)
		if err != nil {
			f.Close()
			panic(err)
		}
		_, err = f.WriteString("[]")
		if err != nil {
			f.Close()
			panic(err)
		}
		f.Close()
	}
}

func syncMeta() {
	for {
		time.Sleep(time.Second * 10)
		log.Info(nil, "sync meta file")
		log.Info(nil, "metadatamap: %v", MetaDataMap)
		if len(MetaDataMap) > 0 {
			log.Info(nil, "metadatamap[0]: %v", *MetaDataMap[0])
			log.Info(nil, "metadatamap.data: %v", MetaDataMap[0].Data)
		}
		f, err := os.OpenFile(BASE_PATH_META, os.O_WRONLY|os.O_TRUNC, 0644)
		if err == nil {
			log.Debug(nil, "json encoding meta file")
			encoder := json.NewEncoder(f)
			err = encoder.Encode(&MetaDataMap)
		}
		if err != nil {
			log.Error(nil, "sync meta file error: %s", err)
		}
	}
}

func startMeta() {
	log.Info(nil, "loading meta file")
	f, err := os.OpenFile(BASE_PATH_META, os.O_RDONLY, 0644)
	if err == nil {
		log.Debug(nil, "json parsing meta file")
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&MetaDataMap)
		if err != nil {
			panic(err)
		}

		for _, metadata := range MetaDataMap {
			var raw []byte
			raw, err = json.Marshal(metadata.Data)
			if err == nil {
				var data any
				switch metadata.Type {
				case TYPE_BYTE:
					data = &MetaDataByte{}
				case TYPE_TEXT:
					data = &MetaDataText{}
				default:
					log.Error([]any{slog.String("Function", "startMeta")}, "metadata type mismathc: %d", metadata.Type)
					continue
				}
				err = json.Unmarshal(raw, data)
				if err == nil {
					metadata.Data = data
				} else {
					break
				}
			}
		}
	}
	if err != nil {
		panic(err)
	}
	go syncMeta()
}

func init() {
	path := config.GLOBAL_CONFIG.Storage.Path
	if path == "" {
		path = "./data"
	}
	log.Info(nil, "initializing storage")
	BASE_PATH = path
	BASE_PATH_TEXT = path + "/text"
	BASE_PATH_BYTE = path + "/byte"
	BASE_PATH_META = path + "/meta.json"
	META_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
	makeDir(BASE_PATH)
	makeDir(BASE_PATH_TEXT)
	makeDir(BASE_PATH_BYTE)
	makeFile(BASE_PATH_META)
	startMeta()
}
