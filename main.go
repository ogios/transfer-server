package main

import (
	"golang.org/x/exp/slog"

	"github.com/ogios/simple-socket-server/server/normal"

	_ "github.com/ogios/transfer-server/config"
	"github.com/ogios/transfer-server/log"
)

func main() {

	log.SetLevel(slog.LevelDebug)
	server, err := normal.NewSocketServer()
	if err != nil {
		log.Error(nil, "Socket server error: &v", err)
		panic(err)
	}

	AddRouters(server)

	if err := server.Serv(); err != nil {
		log.Error(nil, "Serv error: &v", err)
		panic(err)
	}
}