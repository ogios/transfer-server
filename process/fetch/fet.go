package fetch

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/storage/fetch"
)

var FETCH_PARAM_MAXLEN int = 255

func fetchParam(conn *normal.Conn) ([]byte, error) {
	length, err := conn.Si.Next()
	if err != nil {
		return nil, err
	}
	if length > FETCH_PARAM_MAXLEN {
		return nil, fmt.Errorf("fetch param length too long: accept-%d received-%d", FETCH_PARAM_MAXLEN, length)
	}
	return conn.Si.GetSec()
}

func fetchInt(conn *normal.Conn) (int, error) {
	sec, err := fetchParam(conn)
	if err != nil {
		return 0, err
	}
	total := 0
	for index, b := range sec {
		feat := int(math.Pow(255, float64(index)))
		total += int(b) * feat
	}

	return total, nil
}

func FetchMeta(conn *normal.Conn) (err error) {
	defer conn.Close()
	defer func() {
		if e := recover(); e != nil {
			log.Error(nil, "Error in fetch: %v", e)
			err = e.(error)
			log.Debug(nil, "Adding error to output")
			conn.So.AddBytes([]byte("error"))
			conn.So.AddBytes([]byte(err.Error()))
			log.Debug(nil, "Writing error to output")
			conn.So.WriteTo(conn.Raw)
			log.Debug(nil, "Writing error to output done")
		}
	}()
	index, err := fetchInt(conn)
	if err != nil {
		panic(err)
	}
	size, err := fetchInt(conn)
	if err != nil {
		panic(err)
	}

	log.Debug(nil, "fetch params: index-%d size-%d", index, size)
	metas, err := fetch.Fetch(index, size)
	if err != nil {
		panic(err)
	}
	bs, err := json.Marshal(metas)
	if err != nil {
		panic(err)
	}
	err = conn.So.AddBytes([]byte{200})
	if err != nil {
		panic(err)
	}
	err = conn.So.AddBytes(bs)
	if err != nil {
		panic(err)
	}
	err = conn.So.WriteTo(conn.Raw)
	if err != nil {
		panic(err)
	}

	return nil
}
