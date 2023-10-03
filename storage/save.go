package storage

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

var TEXT_FILE_MAX_SIZE int = 16 * 1024

func getTextFile() (*os.File, error) {
	files, err := os.ReadDir(BASE_PATH_TEXT)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return os.Create("1")
	}
	m := 2
	var ff os.DirEntry = nil
	for _, f := range files {
		ind, err := strconv.Atoi(f.Name())
		if err == nil {
			if ind > m {
				m = ind
				ff = f
			}
		}
	}
	if ff == nil {
		return os.Create(strconv.Itoa(m))
	} else {
		info, err := ff.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() > 
	}
	return os.Create(strconv.Itoa(m))
}

func save(r io.Reader, f *os.File) error {
	defer f.Close()
	temp := make([]byte, 1024)
	reader := bufio.NewReader(r)
	for {
		read, err := reader.Read(temp)
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
		f.Write(temp[:read])
	}
}

func SaveText(reader io.Reader) error {
	TEXT_FILE_LOCK.L.Lock()
	defer TEXT_FILE_LOCK.L.Unlock()
	f, err := getTextFile()
	if err != nil {
		return err
	}
	err = save(reader, f)
	if err != nil {
		return err
	}
}
