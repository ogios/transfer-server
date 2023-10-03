package byte

import (
	"fmt"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
)

func ByteCallback(conn *normal.Conn) error {
	data, err := conn.Si.GetSec()
	if err != nil {
		log.Error(nil, "Byte data get error: %v", err)
		return err
	}
	fmt.Println(data)

	return nil
}
