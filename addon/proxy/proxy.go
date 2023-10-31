package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ogios/sutils"
	"github.com/ogios/transfer-server/config"
	"github.com/ogios/transfer-server/log"
)

type HostAddrs struct {
	V6 []string
	V4 []string
}

func StartProxy() {
	if config.GlobalConfig.ProxyEnabled {
		go startServ()
	}
}

func startServ() {
	for {
		conn, err := net.Dial("tcp", config.GlobalConfig.ProxyAddress)
		if err != nil {
			panic(err)
		}
		process(conn)
		time.Sleep(time.Second * 30)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if err := recover(); err != nil {
			log.Error(nil, "Proxy process ERROR: %v", err)
		}
	}()

	// set timeout 20s
	err := conn.SetDeadline(time.Now().Add(time.Second * 20))
	if err != nil {
		panic(err)
	}

	// make data and send
	so := sutils.NewSBodyOUT()
	err = so.AddBytes([]byte("server"))
	if err != nil {
		panic(err)
	}
	addrs := GetInetAddr()
	saddrs, err := json.Marshal(addrs)
	if err != nil {
		panic(err)
	}
	err = so.AddBytes([]byte(saddrs))
	if err != nil {
		panic(err)
	}
	err = so.WriteTo(conn)
	if err != nil {
		panic(err)
	}

	// read response
	buf := make([]byte, 32)
	read, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	if read == 1 && buf[0] == 200 {
		log.Info(nil, "Proxy address updated")
	} else {
		panic(fmt.Errorf("something wrong about the response: %d - %+v", read, buf))
	}
}

func GetInetAddr() *HostAddrs {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	ha := &HostAddrs{
		V6: make([]string, 0),
		V4: make([]string, 0),
	}
	for _, addr := range addrs {
		str := addr.String()
		if strings.Contains(str, ":") {
			ha.V6 = append(ha.V6, str)
		} else {
			ha.V4 = append(ha.V4, str)
		}
	}
	fmt.Println(*ha)
	return ha
}
