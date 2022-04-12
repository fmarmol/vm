package vm

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/fmarmol/vm/pkg/mem"
	"github.com/fmarmol/vm/pkg/prog"
)

func Load(r io.Reader) (*VM, error) {
	var metaInnerVM MetaInnerVM

	err := binary.Read(r, binary.BigEndian, &metaInnerVM)
	if err != nil {
		return nil, fmt.Errorf("could not load metadata: %w", err)
	}

	var innerVM InnerVM

	innerVM.Memory = make(
		mem.Memory,
		metaInnerVM.MemorySize,
		metaInnerVM.MemorySize,
	)
	// read memory
	err = binary.Read(r, binary.BigEndian, &innerVM.Memory)
	if err != nil {
		return nil, fmt.Errorf("could not load memory: %w", err)
	}

	innerVM.Program = make(
		prog.Program,
		metaInnerVM.ProgramSize,
		metaInnerVM.ProgramSize,
	)
	// read program
	err = binary.Read(r, binary.BigEndian, &innerVM.Program)
	if err != nil {
		return nil, fmt.Errorf("could not load program: %w", err)
	}

	return &VM{
		InnerVM:     innerVM,
		MetaInnerVM: metaInnerVM,
	}, nil
}
