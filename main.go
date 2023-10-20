package main

import (
	"golang.org/x/exp/slog"

	"github.com/ogios/simple-socket-server/server/normal"

	_ "github.com/ogios/transfer-server/config"
	"github.com/ogios/transfer-server/log"
)

func main() {
	if GlobalConfig.Debug {
		log.SetLevel(slog.LevelDebug)
	}

	server, err := normal.NewSocketServer(GlobalConfig.Address)
	if err != nil {
		log.Error(nil, "Socket server error: &v", err)
		panic(err)
	}

	AddRouters(server)

	log.Info(nil, "Start serving")
	if err := server.Serv(); err != nil {
		log.Error(nil, "Serv error: &v", err)
		panic(err)
	}
}
