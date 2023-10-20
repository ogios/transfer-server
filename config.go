package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Address string
	Debug   bool
}

var GlobalConfig Config = Config{
	Debug:   false,
	Address: ":15001",
}

func init() {
	debug := flag.Bool("debug", false, "start with debug logging")
	host := flag.String("h", "", "listen host, default empty string")
	port := flag.String("p", "15001", "listen port, default 9977")
	flag.Parse()
	fmt.Println("Flag parsed: ", *debug, *host, *port)
	GlobalConfig.Address = (*host + ":" + *port)
	GlobalConfig.Debug = *debug
}
