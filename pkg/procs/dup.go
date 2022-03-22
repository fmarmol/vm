package procs

import "github.com/fmarmol/vm/pkg/inst"

func Dup(vm VMer, _inst inst.Inst) error {
	w, err := vm.StackPeekIndex(_inst.Operand.UInt32())
	if err != nil {
		return err
	}
	return vm.StackPush(w)
}
