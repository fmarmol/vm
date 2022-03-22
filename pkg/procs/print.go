package procs

import (
	"fmt"

	"github.com/fmarmol/vm/pkg/inst"
)

func Print(vm VMer, _inst inst.Inst) error {
	top, err := vm.StackPop()
	if err != nil {
		return err
	}
	fmt.Println("->", top)
	return nil
}
