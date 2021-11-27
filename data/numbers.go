package data

import (
	"bytes"
	"encoding/binary"
)

func itob(v uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, v)
  return buf.Bytes(), err
}

func btoi(b []byte) (uint64, error) {
	var i uint64

	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}

	return i, nil
}