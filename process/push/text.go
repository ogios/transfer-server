package push

import (
	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/storage/save"
)

func PushText(conn *normal.Conn) error {
	// data, err := conn.Si.GetSec()
	// if err != nil {
	// 	log.Error(nil, "Text data get error: %v", err)
	// 	return err
	// }
	// fmt.Println(string(data))
	defer conn.Close()
	log.Info(nil, "Storage text start")
	err := save.SaveText(conn)
	log.Info(nil, "Storage text done")
	if err != nil {
		return err
	}
	return nil
}
