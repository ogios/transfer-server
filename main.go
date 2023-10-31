package main

import (
	"github.com/ogios/simple-socket-server/server/normal"

	"github.com/ogios/transfer-server/addon/proxy"
	"github.com/ogios/transfer-server/addon/udps"
	"github.com/ogios/transfer-server/config"
	"github.com/ogios/transfer-server/log"
)

func main() {
	server, err := normal.NewSocketServer(config.GlobalConfig.Address)
	if err != nil {
		log.Error(nil, "Socket server error: &v", err)
		panic(err)
	}

	AddRouters(server)

	log.Info(nil, "Start serving")
	udps.StartUdps()
	proxy.StartProxy()
	if err := server.Serv(); err != nil {
		log.Error(nil, "Serv error: &v", err)
		panic(err)
	}
}
