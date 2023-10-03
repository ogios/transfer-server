package text

import (
	"fmt"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
)

func TextCallback(conn *normal.Conn) error {
	data, err := conn.Si.GetSec()
	if err != nil {
		log.Error(nil, "Text data get error: %v", err)
		return err
	}
	fmt.Println(string(data))

	return nil
}
