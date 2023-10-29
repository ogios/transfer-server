package udps

import (
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/ogios/transfer-server/log"
)

var (
	CONN_OFFLINE_TIME = time.Second * 30
	MSG_LOST_TIME     = time.Second * 10
	TAKE_DATA_TIME    = time.Second * 5
)

type Udps struct {
	Conn   *net.UDPConn
	Submap map[string]*UdpsCtx
}

var GlobalUdps Udps

func StartUdps() {
	GlobalUdps = Udps{
		Conn:   createServer(),
		Submap: map[string]*UdpsCtx{},
	}
	go GlobalUdps.heartBeat()
	go GlobalUdps.clearSub()
	go GlobalUdps.startServ()
}

// create and return server
func createServer() *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", ":15002")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	return conn
}

// listen and make goroutine to process bytes received
func (u *Udps) startServ() {
	log.Info(nil, "Start Udps serv")
	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := u.Conn.ReadFromUDP(buf)
		if err != nil {
			panic(fmt.Sprintln("Error reading UDP message:", err))
		}
		if n == 0 {
			continue
		}
		temp := make([]byte, n)
		copy(temp, buf[:n])
		go u.process(clientAddr, temp)
	}
}

// process bytes
//
// if `sub` add `addr` to `submap` else add it to msg list
func (u *Udps) process(addr *net.UDPAddr, buf []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(nil, "Udps process ERROR: %v", err)
		}
	}()
	saddr := addr.String()
	msg := string(buf)
	uc, ok := u.Submap[saddr]
	if msg == "sub" {
		log.Info(nil, "[Subscribe] %s sub", saddr)
		_, err := u.Conn.WriteToUDP([]byte("ok"), addr)
		if err != nil {
			panic(err)
		}
		if ok {
			uc.RefreshOffline()
		} else {
			u.Submap[saddr] = NewUdpsCtx()
		}
	} else {
		if ok {
			log.Debug(nil, "add data to %s", saddr)
			uc.AddData(buf)
		}
	}
}

// heartbeat runs every 30s
func (u *Udps) heartBeat() {
	for {
		log.Info(nil, "sending heart beat")
		for key := range u.Submap {
			log.Debug(nil, "sending heart beat to: %s", key)
			err := u.sendHeartBeat(key)
			if err != nil {
				log.Error(nil, err.Error())
				continue
			}
			go func(addr string) {
				log.Debug(nil, "wating heart beat res: %s", addr)
				err = u.afterHeartBeat(addr)
				if err != nil {
					log.Debug(nil, "heart beat no res, checking offline: %s", addr)
					u.checkOffline(addr)
				}
			}(key)
		}
		time.Sleep(time.Duration(CONN_OFFLINE_TIME))
	}
}

// send heartbeat data = [1]
func (u *Udps) sendHeartBeat(saddr string) error {
	addr, err := net.ResolveUDPAddr("udp", saddr)
	if err != nil {
		msg := fmt.Sprintf("heart beat addr resolve ERROR for: %s", saddr)
		log.Error(nil, msg)
		return fmt.Errorf(msg)

	}
	_, err = u.Conn.WriteTo([]byte{1}, addr)
	if err != nil {
		msg := fmt.Sprintf("heart beat addr write ERROR for: %s", saddr)
		log.Error(nil, msg)
		return fmt.Errorf(msg)
	}
	return nil
}

// wait until heart beat msg returned (5s timeout)
func (u *Udps) afterHeartBeat(addr string) error {
	uc, ok := u.Submap[addr]
	if !ok {
		return nil
	}
	for {
		b, err := uc.TakeData(TAKE_DATA_TIME)
		if err != nil {
			return err
		}
		if len(b) == 1 {
			if b[0] == 1 {
				return nil
			}
		}
	}
}

// runs after heartbeat of one client failed
//
// send heartbeat 3 times
//
// if no response then client will be deemed as offline
func (u *Udps) checkOffline(saddr string) {
	uc, ok := u.Submap[saddr]
	if !ok {
		return
	}
	for i := 0; i < 3; i++ {
		err := u.sendHeartBeat(saddr)
		if err != nil {
			log.Error(nil, err.Error())
			continue
		}
		err = u.afterHeartBeat(saddr)
		if err == nil {
			return
		}
	}
	uc.MakeOffline()
}

// clear offline sub from Submap every 10s
func (u *Udps) clearSub() {
	for {
		log.Debug(nil, "clearing sub")
		for key, val := range u.Submap {
			val.Lock.Lock()
			if val.IsOffline() {
				log.Debug(nil, "clearing sub: &s", key)
				delete(u.Submap, key)
			} else {
				temp := make([]*MsgCtx, len(val.Msgs))
				index := 0
				for _, msg := range val.Msgs {
					if !msg.IsLost() {
						temp[index] = msg
						index++
					}
				}
				if index != len(val.Msgs)-1 {
					val.Msgs = append([]*MsgCtx{}, temp[:index]...)
					// val.Msgs = temp[:index]
				}
			}
			val.Lock.Unlock()
		}
		runtime.GC()
		time.Sleep(time.Duration(MSG_LOST_TIME))
	}
}

func (u *Udps) BoardCast(data []byte) error {
	for key, val := range u.Submap {
		if val.IsOffline() {
			continue
		}
		addr, err := net.ResolveUDPAddr("udp", key)
		if err != nil {
			return err
		}
		_, err = u.Conn.WriteToUDP(data, addr)
		if err != nil {
			return err
		}
	}
	return nil
}
