package vm

import (
	"encoding/binary"
	"io"
)

func (v *VM) WriteToFile(w io.Writer) error {
	v.MetaInnerVM.MemorySize = uint16(len(v.InnerVM.Memory))
	v.MetaInnerVM.ProgramSize = uint16(len(v.InnerVM.Program))

	err := binary.Write(w, binary.BigEndian, v.MetaInnerVM)
	if err != nil {
		return err
	}

	// write mem
	err = binary.Write(w, binary.BigEndian, v.Memory)
	if err != nil {
		return err
	}

	// write program second
	err = binary.Write(w, binary.BigEndian, v.Program)
	if err != nil {
		return err
	}

	return nil
}
