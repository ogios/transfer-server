package main

import (
	"fmt"
	"time"

	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/log"
	"github.com/ogios/transfer-server/process/fetch"
	"github.com/ogios/transfer-server/process/push"
)

func AddRouters(server *normal.Server) {
	server.AddMiddlewareOnEnd(func(conn *normal.Conn, err any) {
		if err != nil {
			defer conn.Close()
			log.Debug(nil, "sending error to connection<%s>: %s", conn.Raw.RemoteAddr().String(), err)
			e := []byte(fmt.Sprintf("%s", err))
			conn.So.AddBytes([]byte("error"))
			conn.So.AddBytes(e)
			conn.So.WriteTo(conn.Raw)
		}
	})

	server.AddMiddlewareOnStart(func(conn *normal.Conn) error {
		err := conn.Raw.SetDeadline(time.Now().Add(time.Second * 10))
		if err != nil {
			return err
		}
		// err = conn.Raw.SetReadDeadline(time.Now().Add(time.Second * 10))
		// if err != nil {
		// 	return err
		// }
		return nil
	})

	// push
	server.AddTypeCallback("text", push.PushText)
	server.AddTypeCallback("byte", push.PushByte)

	// fetch
	server.AddTypeCallback("fetch", fetch.FetchMeta)
	server.AddTypeCallback("fetch_byte", fetch.FetchFile)
}
