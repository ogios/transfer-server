package push

import "github.com/ogios/transfer-server/addon/udps"

func Notify() {
	udps.GlobalUdps.BoardCast([]byte{2})
}
