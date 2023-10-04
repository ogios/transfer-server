package save

import "sync"

var BYTE_FILE_LOCK sync.Cond

func init() {
	BYTE_FILE_LOCK = *sync.NewCond(&sync.Mutex{})
}
