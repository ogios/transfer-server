package delete

import (
	"github.com/ogios/simple-socket-server/server/normal"
	"github.com/ogios/transfer-server/log"
	sd "github.com/ogios/transfer-server/storage/delete"
	"golang.org/x/exp/slog"
)

func DeleteByID(conn *normal.Conn) error {
	defer conn.Close()
	log.Info([]any{slog.String("addr", conn.Raw.RemoteAddr().String())}, "Start deleting")
	err := sd.DeleteData(conn)
	log.Info([]any{slog.String("addr", conn.Raw.RemoteAddr().String())}, "Delete done")
	if err != nil {
		return err
	}
	return nil
}

// sync metadata where is deleted
func SyncMeta(conn *normal.Conn) error {
	return nil
}
