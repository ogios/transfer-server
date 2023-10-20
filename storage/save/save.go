package save

import (
	"io"
	"os"

	"github.com/ogios/simple-socket-server/server/normal"
	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/process"
)

func save(conn *normal.Conn, f *os.File) (start int64, end int64, err error) {
	defer f.Close()
	bufsize := 1024
	total, err := conn.Si.Next()
	if err != nil {
		return 0, 0, err
	}
	log.Info(nil, "total length: %d", total)
	start, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, 0, err
	}
	log.Debug(nil, "start offset: %d", start)
	temp := make([]byte, bufsize)
	var read int
	for {
		read, err = conn.Si.Read(temp)
		if err != nil {
			return 0, 0, err
		}
		_, err = f.Write(temp[:read])
		if err != nil {
			return 0, 0, err
		}
		total -= read
		log.Debug(nil, "for total: %d", total)
		if total == 0 {
			end, err = f.Seek(0, io.SeekEnd)
			if err != nil {
				return 0, 0, err
			}
			log.Debug(nil, "end offset: %d", start)
			err = f.Sync()
			if err != nil {
				return 0, 0, err
			}
			return start, end, nil
		} else if total < bufsize {
			temp = make([]byte, total)
		}
		continue
	}
}

func WriteSuccess(conn *normal.Conn) error {
	err := conn.So.AddBytes([]byte(process.STATUS_SUCCESS))
	if err != nil {
		log.Error(nil, "add success response error: %s", err)
		return err
	}
	err = conn.So.WriteTo(conn.Raw)
	if err != nil {
		log.Error(nil, "write response error: %s", err)
		return err
	}
	return nil
}
