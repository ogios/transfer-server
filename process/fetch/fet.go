package fetch

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/storage"
	"github.com/ogios/transfer-server/storage/fetch"
)

var FETCH_PARAM_MAXLEN int = 255

type FetchRes struct {
	Data  []*storage.MetaData `json:"data"`
	Total int                 `json:"total"`
}

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
	for index, b := range sec[:len(sec)-1] {
		feat := int(math.Pow(255, float64(index)))
		total += int(b) * feat
	}

	return total, nil
}

func FetchMeta(conn *normal.Conn) (err error) {
	log.Debug(nil, "fetching index")
	index, err := fetchInt(conn)
	if err != nil {
		return err
	}
	log.Debug(nil, "fetching size")
	size, err := fetchInt(conn)
	if err != nil {
		return err
	}

	log.Debug(nil, "fetch params: index-%d size-%d", index, size)
	metas, total, err := fetch.Fetch(index, size)
	if err != nil {
		return err
	}
	data := FetchRes{
		Total: total,
		Data:  metas,
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = conn.So.AddBytes([]byte{200})
	if err != nil {
		return err
	}
	err = conn.So.AddBytes(bs)
	if err != nil {
		return err
	}
	err = conn.So.WriteTo(conn.Raw)
	if err != nil {
		return err
	}

	return nil
}
