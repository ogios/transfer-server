package delete

import (
	"fmt"

	"github.com/ogios/simple-socket-server/server/normal"
	"github.com/ogios/transfer-server/process"
	"github.com/ogios/transfer-server/storage"
)

func getID(conn *normal.Conn) (string, error) {
	length, err := conn.Si.Next()
	if err != nil {
		return "", err
	}
	if length != storage.ID_LENGTH {
		return "", fmt.Errorf("id length mismatch: %d-%d", storage.ID_LENGTH, length)
	}
	id, err := conn.Si.GetSec()
	return string(id), err
}

func SendSuccess(conn *normal.Conn) error {
	err := conn.So.AddBytes([]byte(process.STATUS_SUCCESS))
	if err != nil {
		return err
	}
	err = conn.So.WriteTo(conn.Raw)
	if err != nil {
		return err
	}
	return nil
}

func DeleteData(conn *normal.Conn) error {
	id, err := getID(conn)
	if err != nil {
		return err
	}
	storage.DeleteMetaData(id)
	if err = SendSuccess(conn); err != nil {
		return err
	}
	return nil
}
