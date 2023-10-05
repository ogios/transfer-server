package fetch

import (
	"fmt"
	"os"
)

func FetchText(start int64, end int64, filename string) (string, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	defer f.Close()
	if err != nil {
		return "", err
	}
	b := make([]byte, end-start)
	read, err := f.ReadAt(b, start)
	if err != nil {
		return "", err
	}
	if int64(read) != (end - start) {
		return "", fmt.Errorf("read length not right: read-%d length-%d", read, end-start)
	}
	return string(b), nil
}
