package procs

import (
	"fmt"

	"github.com/fmarmol/vm/pkg/fatal"
	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/word"
)

// DEBUG instruction DO NOT CONSUME top stack
func Eq(v VMer, _inst inst.Inst) error {
	top := v.StackPeek()
	op := _inst.Operand
	if top.Kind != op.Kind {
		return fmt.Errorf("incompatible types: tried comparison between %v and %v\n", top.Kind, op.Kind)
	}

	switch top.Kind {
	case word.Float64:
		if top.Float64() != op.Float64() {
			fatal.Panic("invalid assertion top[%v] != eq[%v]", top.Float64(), op.Float64())
		}
	case word.Int64:
		if top.Int64() != op.Int64() {
			fatal.Panic("invalid assertion top[%v] != eq[%v]", top.Int64(), op.Int64())
		}
	case word.UInt32:
		if top.UInt32() != op.UInt32() {
			fatal.Panic("invalid assertion top[%v] != eq[%v]", top.UInt32(), op.UInt32())
		}
	case word.Ptr:
		if top.Ptr() != op.Ptr() {
			fatal.Panic("invalid assertion top[%v] != eq[%v]", top.Ptr(), op.Ptr())
		}
	default:
		fatal.Panic("eq not implemented for type: %v", top.Kind)

	}
	return nil
}
