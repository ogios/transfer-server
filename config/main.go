package config

import (
	"flag"
	"fmt"
)

type Config struct {
	ProxyAddress string
	ProxyKey     string
	Address      string
	Path         string
	ProxyEnabled bool
	SubEnabled   bool
	Debug        bool
}

var GlobalConfig Config = Config{
	Debug:        false,
	Address:      ":15001",
	Path:         "",
	ProxyAddress: "",
	ProxyKey:     "",
	ProxyEnabled: false,
	SubEnabled:   true,
}

func init() {
	proxy := flag.String("proxy", "", "proxy host, send host address to it, leave empty to disable")
	proxy_key := flag.String("proxy-key", "", "proxy server identifier, it can be any thing, works only when proxy enabled")
	debug := flag.Bool("debug", false, "start with debug logging")
	sub := flag.Bool("sub", true, "open message subscription, default true")
	host := flag.String("h", "", "listen host, default empty string")
	port := flag.String("p", "15001", "listen port, default 15001")
	path := flag.String("r", "", "file save path")
	flag.Parse()
	fmt.Println("Flag parsed: ", *debug, *host, *port)
	GlobalConfig.Address = (*host + ":" + *port)
	GlobalConfig.Debug = *debug
	GlobalConfig.SubEnabled = *sub
	GlobalConfig.Path = *path
	if len(*proxy) > 0 {
		GlobalConfig.ProxyEnabled = true
		GlobalConfig.ProxyAddress = *proxy
		if len(*proxy_key) == 0 {
			panic(fmt.Errorf("please provide proxy key when proxy enabled"))
		} else {
			GlobalConfig.ProxyKey = *proxy_key
		}
	}
}
