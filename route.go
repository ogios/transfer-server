package main

import (
	"github.com/ogios/simple-socket-server/server/normal"

	by "github.com/ogios/transfer-server/process/byte"
	"github.com/ogios/transfer-server/process/text"
)

func AddRouters(server *normal.Server) {
	server.AddTypeCallback("text", text.TextCallback)
	server.AddTypeCallback("byte", by.ByteCallback)
}
