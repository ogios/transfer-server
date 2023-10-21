package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/ogios/transfer-server/config"
	"github.com/ogios/transfer-server/log"
	"golang.org/x/exp/slog"
)

var (
	BASE_PATH          string
	BASE_PATH_TEXT     string
	BASE_PATH_BYTE     string
	BASE_PATH_META     string
	BASE_PATH_META_DEL string
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

func loadMeta() {
	// load file
	log.Info(nil, "loading meta file")
	f, err := os.OpenFile(BASE_PATH_META, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	// parse file (map)
	log.Debug(nil, "json parsing meta file")
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&MetaDataMap)
	if err != nil {
		panic(err)
	}

	// map to struct
	for _, metadata := range MetaDataMap {
		var raw []byte
		raw, err = json.Marshal(metadata.Data)
		if err != nil {
			panic(err)
		}
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

func loadMetaIndex() {
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

func loadMetaDel() {
	log.Info(nil, "loading meta_del file")
	f, err := os.OpenFile(BASE_PATH_META_DEL, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	log.Debug(nil, "json parsing meta_del file")
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&MetaDataDelList)
	if err != nil {
		panic(err)
	}
}

func clearDel() {
	err := ClearDeleteFromFile()
	if err != nil {
		panic(err)
	}
	ClearDeleteMetaData()
}

func startMeta() {
	loadMeta()
	loadMetaIndex()
	loadMetaDel()
	clearDel()
	go syncMeta()
}

func syncMeta() {
	for {
		time.Sleep(time.Second * 10)
		log.Info(nil, "sync meta file")
		log.Info(nil, "metadatamap: %v | metadatamap_del: %v", MetaDataMap, MetaDataDelList)

		// MetaDataMap
		f, err := os.OpenFile(BASE_PATH_META, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Error(nil, "open sync meta file error: %s", err)
			continue
		}
		log.Debug(nil, "json encoding meta file")
		encoder := json.NewEncoder(f)
		err = encoder.Encode(&MetaDataMap)
		f.Close()
		if err != nil {
			log.Error(nil, "write sync meta file error: %s", err)
			continue
		}

		// MetaDataDelList
		f, err = os.OpenFile(BASE_PATH_META_DEL, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Error(nil, "open sync meta_file error: %s", err)
			continue
		}
		log.Debug(nil, "json encoding meta_del file")
		encoder = json.NewEncoder(f)
		err = encoder.Encode(&MetaDataDelList)
		f.Close()
		if err != nil {
			log.Error(nil, "write sync meta_file error: %s", err)
			continue
		}

		// GC
		f, encoder = nil, nil
		runtime.GC()
	}
}

func init() {
	path := config.GlobalConfig.Path
	if path == "" {
		path = "./data"
	}
	log.Info(nil, "initializing storage")
	BASE_PATH = path
	BASE_PATH_TEXT = path + "/text"
	BASE_PATH_BYTE = path + "/byte"
	BASE_PATH_META = path + "/meta.json"
	BASE_PATH_META_DEL = path + "/meta_del.json"
	META_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
	makeDir(BASE_PATH)
	makeDir(BASE_PATH_TEXT)
	makeDir(BASE_PATH_BYTE)
	makeFile(BASE_PATH_META)
	makeFile(BASE_PATH_META_DEL)
	startMeta()
}
