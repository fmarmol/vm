package procs

import (
	"fmt"

	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/rorre"
	"github.com/fmarmol/vm/pkg/word"
)

type Number interface {
	~uint32 | ~int64 | ~float64
}

func binaryOp[T Number](arg1 T, arg2 T, ik inst.InstKind) (ret T, err error) {
	switch ik {
	case inst.Inst_Add:
		ret = arg1 + arg2
	case inst.Inst_Sub:
		ret = arg1 - arg2
	case inst.Inst_Mul:
		ret = arg1 * arg2
	case inst.Inst_Div:
		if arg2 == 0 {
			err = rorre.Err_DivisionByZero
		} else {
			ret = arg1 / arg2
		}
	case inst.Inst_Eq:
		if arg1 == arg2 {
			ret = 1
		} else {
			ret = 0
		}
	default:
		panic("unknown binaryOp")
	}
	return
}

func Bin(vm VMer, _inst inst.Inst) error {
	b, err := vm.StackPop()
	if err != nil {
		return err
	}

	a, err := vm.StackPop()
	if err != nil {
		return err
	}
	if a.Kind != b.Kind {
		return fmt.Errorf("incompatible types: tried to binary operation between %v and %v\n", a.Kind, b.Kind)
	}
	var result word.Word
	switch a.Kind {
	case word.Int64:
		res, err := binaryOp(a.Int64(), b.Int64(), _inst.Kind)
		if err != nil {
			return err
		}
		result = word.NewWord(res, word.Int64)
	case word.Float64:
		res, err := binaryOp(a.Float64(), b.Float64(), _inst.Kind)
		if err != nil {
			return err
		}
		result = word.NewWord(res, word.Float64)
	default:
		return fmt.Errorf("binary operation not implemented for type: %v", a.Kind)
	}
	return vm.StackPush(result)
}
