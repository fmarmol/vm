package procs

import (
	"fmt"

	"github.com/fmarmol/vm/pkg/inst"
)

func MemR8(vm VMer, _inst inst.Inst) error {
	top := vm.StackPeek().UInt32()

	m := vm.Mem()
	res := m.Read8(top)
	fmt.Printf("%c", res)
	return nil
}
