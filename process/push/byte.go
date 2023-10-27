package push

import (
	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/storage/save"
)

func PushByte(conn *normal.Conn) error {
	// data, err := conn.Si.GetSec()
	// if err != nil {
	// 	log.Error(nil, "Byte data get error: %v", err)
	// 	return err
	// }
	// fmt.Println(data)

	defer conn.Close()
	log.Info(nil, "Storage text start")
	err := save.SaveByte(conn)
	log.Info(nil, "Storage text done")
	if err != nil {
		return err
	}
	Notify()
	return nil
}
