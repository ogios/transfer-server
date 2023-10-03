package storage

import (
	"os"
	"sync"

	"github.com/ogios/transfer-server/config"
)

var BASE_PATH string
var BASE_PATH_TEXT string
var BASE_PATH_BYTE string
var BASE_PATH_META string

var META_FILE_LOCK sync.Cond
var TEXT_FILE_LOCK sync.Cond
var BYTE_FILE_LOCK sync.Cond

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
	META_FILE_LOCK.L.Lock()
	META_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
	META_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
	defer META_FILE_LOCK.L.Unlock()
	makeDir(BASE_PATH)
	makeDir(BASE_PATH_TEXT)
	makeDir(BASE_PATH_BYTE)
	makeFile(BASE_PATH_META)
}
