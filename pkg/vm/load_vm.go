package vm

import (
	"encoding/binary"
	"fmt"
	"io"
)

func LoadVM(r io.Reader) (*VM, error) {
	var metaInnerVM MetaInnerVM

	err := binary.Read(r, binary.BigEndian, &metaInnerVM)
	if err != nil {
		return nil, fmt.Errorf("could load metadata: %w", err)
	}

	var innerVM InnerVM
	// read memory
	err = binary.Read(r, binary.BigEndian, &innerVM.Memory)
	if err != nil {
		return nil, fmt.Errorf("could load memory: %w", err)
	}

	// read program
	err = binary.Read(r, binary.BigEndian, &innerVM.Program)
	if err != nil {
		return nil, fmt.Errorf("could load program: %w", err)
	}

	return &VM{
		InnerVM:     innerVM,
		MetaInnerVM: metaInnerVM,
	}, nil
}
