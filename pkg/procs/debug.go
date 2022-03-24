package procs

import (
	"fmt"

	"github.com/fmarmol/vm/pkg/inst"
)

func Debug(vm VMer, _inst inst.Inst) error {
	top := vm.StackPeek()
	fmt.Println("->", top)
	return nil
}
