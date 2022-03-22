package procs

import (
	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/rorre"
	"github.com/fmarmol/vm/pkg/word"
)

func Call(vm VMer, _inst inst.Inst) error {
	if _inst.Operand.UInt32() < 0 || _inst.Operand.UInt32() >= vm.ProgramSize() {
		return rorre.Err_OutOfIndexInstruction
	}
	return vm.StackPush(word.NewWord(vm.IP()+1, word.UInt32))
}
