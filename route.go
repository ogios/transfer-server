package main

import (
	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/process/push"
)

func AddRouters(server *normal.Server) {
	// push
	server.AddTypeCallback("text", push.PushText)
	server.AddTypeCallback("byte", push.PushByte)

	// fetch
	server.AddTypeCallback("fetch")
}
