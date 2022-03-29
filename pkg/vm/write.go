package vm

import (
	"encoding/binary"
	"os"
)

func (v InnerVM) WriteToFile(pathFile string) error {
	fd, err := os.Create(pathFile)
	if err != nil {
		return err
	}
	defer fd.Close()
	return binary.Write(fd, binary.BigEndian, v)
}
