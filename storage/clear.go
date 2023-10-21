package storage

import (
	"io"
	"os"
)

func ClearDeleteFromFile() error {
	TEXT_FILE_LOCK.L.Lock()
	defer TEXT_FILE_LOCK.L.Unlock()
	BYTE_FILE_LOCK.L.Lock()
	defer BYTE_FILE_LOCK.L.Unlock()
	b := make([]*MetaData, 0)
	t := make(map[*MetaData]struct{})

	// what text file contains deletion
	f := []string{}
	for _, index := range MetaDataDelList {
		p := MetaDataMap[index]
		switch p.Type {
		case TYPE_BYTE:
			b = append(b, p)
		case TYPE_TEXT:
			t[p] = struct{}{}
			f = append(f, p.Data.(*MetaDataText).Filename)
		}
	}

	err := ClearTextFile(t, f)
	if err != nil {
		return err
	}
	err = ClearByteFile(b)
	if err != nil {
		return err
	}
	return nil
}

func Save(reader, writer *os.File, startoff, endoff int64) (start int64, end int64, err error) {
	bufsize := 1024 * 4
	start, err = writer.Seek(0, io.SeekEnd)
	end = start
	if err != nil {
		return 0, 0, err
	}
	defer reader.Close()
	temp := make([]byte, bufsize)
	left := endoff - startoff
	for {
		if left == 0 {
			end, err = writer.Seek(0, io.SeekEnd)
			if err != nil {
				return 0, 0, err
			}
			return start, end, nil
		}
		var read int
		if left <= int64(bufsize) {
			temp = make([]byte, left)
		}
		read, err = reader.ReadAt(temp, startoff)
		if err != nil {
			return 0, 0, err
		}
		_, err = writer.Write(temp[:read])
		if err != nil {
			return 0, 0, err
		}
		startoff += int64(read)
		left -= int64(read)
	}
}

func ClearTextFile(texts map[*MetaData]struct{}, files []string) error {
	for _, filename := range files {
		tempfilename := filename + ".bak"
		writer, err := os.Create(tempfilename)
		if err != nil {
			return err
		}
		defer writer.Close()
		ms := MetaDataTextMap[filename]
		for _, m := range ms {
			if _, ok := texts[m]; ok {
				continue
			}
			textmeta := m.Data.(*MetaDataText)
			reader, err := os.OpenFile(textmeta.Filename, os.O_RDONLY, 0644)
			if err != nil {
				return err
			}
			defer reader.Close()
			start, end, err := Save(reader, writer, textmeta.Start, textmeta.End)
			if err != nil {
				return err
			}
			textmeta.Start = start
			textmeta.End = end
			reader.Close()
		}
		writer.Sync()
		writer.Close()
		err = os.Remove(filename)
		if err != nil {
			return err
		}
		err = os.Rename(tempfilename, filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func ClearByteFile(bytes []*MetaData) error {
	return nil
}
