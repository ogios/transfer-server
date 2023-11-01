package fetch

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/process"
	"github.com/ogios/transfer-server/storage"
	"github.com/ogios/transfer-server/storage/fetch"
)

var FETCH_MAX_FILENAME int = 255

func FetchFile(conn *normal.Conn) (err error) {
	length, err := conn.Si.Next()
	if err != nil {
		return err
	}
	if length != storage.ID_LENGTH {
		return fmt.Errorf("fetch file id length mismatch")
	}
	id, err := conn.Si.GetSec()
	if err != nil {
		return err
	}
	f, size, err := fetch.FetchByte(string(id))
	if err != nil {
		return err
	}
	defer f.Close()
	defer runtime.GC()
	err = conn.So.AddBytes([]byte(process.STATUS_SUCCESS))
	if err != nil {
		return err
	}
	err = conn.So.AddBytes([]byte(filepath.Base(f.Name())))
	if err != nil {
		return err
	}
	err = conn.So.AddReader(f, int(size))
	if err != nil {
		return err
	}
	err = conn.So.WriteTo(conn.Raw)
	if err != nil {
		return err
	}
	return nil
}
