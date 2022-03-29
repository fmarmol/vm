package vm

import (
	"encoding/binary"
	"fmt"
	"os"
)

func LoadInnerVM(pathFile string) (*InnerVM, error) {
	fd, err := os.Open(pathFile)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var v InnerVM

	err = binary.Read(fd, binary.BigEndian, &v)
	if err != nil {
		return nil, fmt.Errorf("could load file: %v: %w", pathFile, err)
	}
	return &v, nil
}
