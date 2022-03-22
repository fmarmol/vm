package procs

import (
	"github.com/fmarmol/vm/pkg/inst"
)

func Push(v VMer, _inst inst.Inst) error {
	return v.StackPush(_inst.Operand)
}
