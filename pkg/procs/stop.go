package procs

import "github.com/fmarmol/vm/pkg/inst"

func Stop(vm VMer, _ inst.Inst) error {
	vm.Stop()
	return nil
}
