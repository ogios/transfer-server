package udps

import (
	"fmt"
	"sync"
	"time"

	"github.com/ogios/transfer-server/log"
)

type MsgCtx struct {
	Lost time.Time
	Data []byte
}

func NewMsgCtx(data []byte) *MsgCtx {
	m := &MsgCtx{
		Lost: time.Now().Add(time.Duration(time.Second * 10)),
		Data: data,
	}
	return m
}

// check if current msg is out of date
func (mc *MsgCtx) IsLost() bool {
	return time.Now().Compare(mc.Lost) != -1
}

// mark current msg as out of date
func (uc *MsgCtx) MakeLost() {
	uc.Lost = time.Now()
}

type UdpsCtx struct {
	Offline time.Time
	Lock    *sync.Mutex
	Msgs    []*MsgCtx
}

func NewUdpsCtx() *UdpsCtx {
	u := &UdpsCtx{
		Lock: &sync.Mutex{},
	}
	u.RefreshOffline()
	return u
}

// mark current UDP client as out of date
func (uc *UdpsCtx) MakeOffline() {
	uc.Offline = time.Now()
}

// increase life
func (uc *UdpsCtx) RefreshOffline() {
	uc.Offline = time.Now().Add(time.Duration(time.Second * 30))
}

// check if current UDP client is out of date
func (uc *UdpsCtx) IsOffline() bool {
	return time.Now().Compare(uc.Offline) != -1
}

// add MsgCtx to client data for other usage, and refresh life
func (uc *UdpsCtx) AddData(data []byte) {
	uc.Lock.Lock()
	defer uc.Lock.Unlock()
	uc.Msgs = append(uc.Msgs, NewMsgCtx(data))
	uc.RefreshOffline()
}

// take the first data from Data, if no then wait until given timeout
//
// stop if timeout or UDP client is out of date
func (uc *UdpsCtx) TakeData(timeout time.Duration) ([]byte, error) {
	stop := false
	ch := make(chan []byte)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error(nil, "Take Data ERROR: %v", err)
			}
		}()
		for !stop && !uc.IsOffline() {
			uc.Lock.Lock()
			for _, val := range uc.Msgs {
				lost := val.IsLost()
				if !lost {
					log.Debug(nil, "Taking: %v", val.Data)
					val.MakeLost()
					temp := make([]byte, len(val.Data))
					copy(temp, val.Data)
					uc.Lock.Unlock()
					ch <- temp
				}
			}
			uc.Lock.TryLock()
			uc.Lock.Unlock()
			time.Sleep(time.Microsecond * 10)
		}
	}()
	select {
	case <-time.After(timeout):
		return nil, fmt.Errorf("TakeData timeout")
	case b := <-ch:
		stop = true
		return b, nil
	}
}
