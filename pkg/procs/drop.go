package procs

import "github.com/fmarmol/vm/pkg/inst"

func Drop(vm VMer, _inst inst.Inst) error {
	_, err := vm.StackPop()
	return err
}
