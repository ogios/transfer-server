package storage

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/ogios/transfer-server/config"
)

var BASE_PATH string
var BASE_PATH_TEXT string
var BASE_PATH_BYTE string
var BASE_PATH_META string

func makeDir(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
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
		f, err := os.OpenFile(BASE_PATH_META, os.O_WRONLY, 0644)
		if err == nil {
			encoder := json.NewEncoder(f)
			err = encoder.Encode(&MetaDataMap)
		}
		if err != nil {

		}
	}
}

func startMeta() {
	f, err := os.OpenFile(BASE_PATH_META, os.O_RDONLY, 0644)
	if err == nil {
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
	BASE_PATH = path
	BASE_PATH_TEXT = path + "/text"
	BASE_PATH_BYTE = path + "/byte"
	BASE_PATH_META = path + "/meta.json"
	META_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
	META_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
	META_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
	makeDir(BASE_PATH)
	makeDir(BASE_PATH_TEXT)
	makeDir(BASE_PATH_BYTE)
	makeFile(BASE_PATH_META)
	startMeta()
}
