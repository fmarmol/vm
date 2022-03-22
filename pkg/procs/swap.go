package procs

import "github.com/fmarmol/vm/pkg/inst"

// Swap the the top and the nth elements
func Swap(vm VMer, _inst inst.Inst) error {
	return vm.Swap(1, _inst.Operand.UInt32())
}
