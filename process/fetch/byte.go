package fetch

import (
	"fmt"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/process"
	"github.com/ogios/transfer-server/storage/fetch"
)

var FETCH_MAX_FILENAME int = 255

func FetchFile(conn *normal.Conn) (err error) {
	length, err := conn.Si.Next()
	if err != nil {
		return err
	}
	if length > FETCH_MAX_FILENAME {
		return fmt.Errorf("filename too long")
	}
	name, err := conn.Si.GetSec()
	if err != nil {
		return err
	}
	f, size, err := fetch.FetchByte(string(name))
	if err != nil {
		return err
	}
	err = conn.So.AddBytes([]byte(process.STATUS_SUCCESS))
	if err != nil {
		return err
	}
	err = conn.So.AddBytes([]byte(name))
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
