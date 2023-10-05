package save

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/storage"
)

var BYTE_FILE_LOCK sync.Cond

func init() {
	BYTE_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
}

func SaveByte(conn *normal.Conn) error {
	filename, err := getFilename(conn)
	if err != nil {
		log.Error(nil, "Get filename error: %s", err)
		return err
	}
	log.Debug(nil, "byte filename: %s", filename)
	f, err := getByteFile(filename)
	if err != nil {
		log.Error(nil, "Get byte file error: %s", err)
		return err
	}
	log.Debug(nil, "byte filename get: %s", f.Name())
	log.Debug(nil, "saving filename")
	start, end, err := save(conn, f)
	if err != nil {
		log.Error(nil, "save byte file error: %s", err)
		return err
	}
	log.Debug(nil, "save byte file done: length-%d", end-start)
	log.Debug(nil, "saving byte metadata...")
	storage.AddMetaData(storage.MetaData{
		Type: storage.TYPE_BYTE,
		Data: storage.MetaDataByte{
			Filename: filename,
			Size:     end - start,
		},
	})

	return nil
}

func getByteFile(name string) (*os.File, error) {
	BYTE_FILE_LOCK.L.Lock()
	defer BYTE_FILE_LOCK.L.Unlock()
	log.Debug(nil, "Matching file name")
	name, err := matchFilename(name)
	if err != nil {
		log.Error(nil, "Match file name error: %s", err)
		return nil, err
	}
	return os.Create(storage.BASE_PATH_BYTE + "/" + name)
}

func matchFilename(name string) (string, error) {
	log.Debug(nil, "Reading byte dir")
	files, err := os.ReadDir(storage.BASE_PATH_BYTE)
	if err != nil {
		log.Error(nil, "Read byte dir error: %s", err)
		return "", err
	}
	for _, file := range files {
		if file.Name() == name {
			log.Debug(nil, "Compiling regexp")
			reg, err := regexp.Compile(getFilenameReg(name))
			if err != nil {
				log.Error(nil, "regexp compile error: %s", err)
				return "", err
			}
			log.Debug(nil, "matching max byte file name")
			last, err := matchMaxFilename(reg, files)
			if err != nil {
				log.Error(nil, "match max byte file name error: %s", err)
				return "", err
			}
			return makeFilename(name, last+1), nil
		}
	}
	return name, err
}

func makeFilename(name string, num int) string {
	suffix := ""
	if dot := strings.LastIndex(name, "."); dot != -1 {
		suffix = name[dot:]
		name = name[:dot]
	}
	return name + "(" + strconv.Itoa(num) + ")" + suffix
}

func getFilenameReg(name string) string {
	suffix := ""
	if dot := strings.LastIndex(name, "."); dot != -1 {
		suffix = "\\" + name[dot:]
		name = name[:dot]
	}
	return name + `\(([0-9]*?)\)` + suffix + "$"
}

func matchMaxFilename(reg *regexp.Regexp, fs []fs.DirEntry) (int, error) {
	maxd := 0
	for _, s := range fs {
		l := reg.FindStringSubmatch(s.Name())
		if len(l) > 0 {
			last := l[len(l)-1]
			if len(last) > 0 {
				count, err := strconv.Atoi(last)
				if err != nil {
					return 0, err
				}
				if count > maxd {
					maxd = count
				}
			}
		}
	}
	return maxd, nil
}

func getFilename(conn *normal.Conn) (string, error) {
	length, err := conn.Si.Next()
	if err != nil {
		return "", err
	}
	if length > 255 {
		return "", fmt.Errorf("filename too long: %d", length)
	}
	name, err := conn.Si.GetSec()
	if err != nil {
		return "", err
	}
	return string(name), nil
}
