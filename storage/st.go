package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ogios/transfer-server/config"
	"github.com/ogios/transfer-server/log"
)

var BASE_PATH string
var BASE_PATH_TEXT string
var BASE_PATH_BYTE string
var BASE_PATH_META string

func makeDir(dir string) error {
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
	return nil
}

func makeFile(path string) error {
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
	return nil
}

func syncMeta() {
	for {
		time.Sleep(time.Second * 10)
		log.Info(nil, "sync meta file")
		f, err := os.OpenFile(BASE_PATH_META, os.O_WRONLY, 0644)
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
