package mem

import (
	"errors"
	"fmt"
	"unsafe"
)

type Memory []byte

func (m *Memory) Len() uint32 { return uint32(len(*m)) }

func (m *Memory) Read8(addr uint32) uint8 {
	return (*m)[addr]
}

func (m *Memory) Read16(addr uint32) uint16 {
	return *(*uint16)(unsafe.Pointer(&(*m)[addr]))
}
func (m *Memory) Write8(value uint8, addr uint32) {
	(*m)[addr] = value
}

func (m *Memory) Write16(value uint16, addr uint32) error {
	if addr > m.Len()-2 {
		return errors.New("MEMORY ILLEGAL ACCESS")
	}
	*(*uint16)(unsafe.Pointer(&(*m)[addr])) = value
	return nil
}

func (m *Memory) Dump() {
	for _, b := range *m {
		fmt.Printf("%02X ", b)
	}
	fmt.Println("")
}
