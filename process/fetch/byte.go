package fetch

import (
	"fmt"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/process"
	"github.com/ogios/transfer-server/storage/fetch"
)

var FETCH_MAX_FILENAME int = 255

func FetchFile(conn *normal.Conn) (err error) {
	defer conn.Close()
	defer func() {
		if e := recover(); e != nil {
			log.Error(nil, "Error in fetch: %v", e)
			err = e.(error)
			log.Debug(nil, "Adding error to output")
			conn.So.AddBytes([]byte(process.STATUS_ERROR))
			conn.So.AddBytes([]byte(err.Error()))
			log.Debug(nil, "Writing error to output")
			conn.So.WriteTo(conn.Raw)
			log.Debug(nil, "Writing error to output done")
		}
	}()

	length, err := conn.Si.Next()
	if err != nil {
		panic(err)
	}
	if length > FETCH_MAX_FILENAME {
		panic(fmt.Errorf("filename too long"))
	}
	name, err := conn.Si.GetSec()
	if err != nil {
		panic(err)
	}
	f, size, err := fetch.FetchByte(string(name))
	if err != nil {
		panic(err)
	}
	err = conn.So.AddBytes([]byte(process.STATUS_SUCCESS))
	if err != nil {
		panic(err)
	}
	err = conn.So.AddBytes([]byte(name))
	if err != nil {
		panic(err)
	}
	err = conn.So.AddReader(f, int(size))
	if err != nil {
		panic(err)
	}
	err = conn.So.WriteTo(conn.Raw)
	if err != nil {
		panic(err)
	}
	return nil
}
