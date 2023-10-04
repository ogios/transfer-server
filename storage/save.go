package storage

import (
	"io"
	"os"
	"strconv"

	"github.com/ogios/simple-socket-server/server/normal"
)

var TEXT_FILE_MAX_SIZE int64 = 16 * 1024

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
		if info.Size() >= TEXT_FILE_MAX_SIZE {
			return os.Create(strconv.Itoa(m + 1))
		} else {
			return os.OpenFile(strconv.Itoa(m), os.O_RDWR, 0644)
		}
	}
}

func save(conn *normal.Conn, f *os.File) (start int64, end int64, err error) {
	defer f.Close()
	bufsize := 1024
	total, err := conn.Si.Next()
	if err == nil {
		start, err = f.Seek(0, io.SeekEnd)
		if err == nil {
			temp := make([]byte, bufsize)
			for {
				read, err := conn.Si.Read(temp)
				if err != nil {
					break
				}
				f.Write(temp[:read])
				total -= read
				if total == 0 {
					end, err := f.Seek(0, io.SeekEnd)
					if err != nil {
						break
					}
					return start, end, nil
				} else if total < bufsize {
					temp = make([]byte, total)
				}
			}
		}
	}
	return 0, 0, err
}

func saveText(reader *normal.Conn) (int64, int64, error) {
	TEXT_FILE_LOCK.L.Lock()
	defer TEXT_FILE_LOCK.L.Unlock()
	f, err := getTextFile()
	if err != nil {
		return 0, 0, err
	}
	start, end, err := save(reader, f)
	if err != nil {
		return 0, 0, err
	}
	return start, end, err
}

func SaveText(reader *normal.Conn) error {
	start, end, err := saveText(reader)
	if err != nil {
		return err
	}
	AddMetaData(MetaData{
		Type: TYPE_TEXT,
		Data: MetaDataText{
			Start: start,
			End:   end,
		},
	})
	return nil
}
